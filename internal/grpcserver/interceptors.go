package grpcserver

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/stratflow-labs/stratflow/internal/authkit"
	"github.com/stratflow-labs/stratflow/internal/foundation/logger"
)

type MetricsRecorder interface {
	ObserveRPC(ctx context.Context, fullMethod string, code codes.Code, duration time.Duration)
}

type AuthSkipper func(fullMethod string) bool

func SkipMethods(methods ...string) AuthSkipper {
	allowed := make(map[string]struct{}, len(methods))
	for _, method := range methods {
		if method == "" {
			continue
		}
		allowed[method] = struct{}{}
	}

	return func(fullMethod string) bool {
		_, ok := allowed[fullMethod]
		return ok
	}
}

func RecoveryUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (_ any, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("grpc panic recovered",
					"method", info.FullMethod,
					"panic", rec,
					"stack", string(debug.Stack()),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

func RequestMetadataUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		requestID, correlationID := requestMetadata(ctx)
		ctx = WithRequestMetadata(ctx, requestID, correlationID)

		_ = grpc.SetHeader(ctx, metadata.Pairs(
			HeaderRequestID, requestID,
			HeaderCorrelationID, correlationID,
		))

		return handler(ctx, req)
	}
}

func LoggingUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		code := status.Code(err)
		requestID, _ := RequestIDFromContext(ctx)
		correlationID, _ := CorrelationIDFromContext(ctx)
		userID, _ := authkit.UserIDFromContext(ctx)

		logger.Info("grpc request",
			"method", info.FullMethod,
			"code", code.String(),
			"duration", time.Since(start).Seconds(),
			"request_id", requestID,
			"correlation_id", correlationID,
			"user_id", userID,
		)

		return resp, err
	}
}

func MetricsUnaryInterceptor(recorder MetricsRecorder) grpc.UnaryServerInterceptor {
	if recorder == nil {
		return func(
			ctx context.Context,
			req any,
			info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler,
		) (any, error) {
			return handler(ctx, req)
		}
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		recorder.ObserveRPC(ctx, info.FullMethod, status.Code(err), time.Since(start))
		return resp, err
	}
}

func AuthUnaryInterceptor(verifier authkit.AccessTokenVerifier, skip AuthSkipper) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if verifier == nil || (skip != nil && skip(info.FullMethod)) {
			return handler(ctx, req)
		}

		token, err := BearerTokenFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "missing Authorization metadata")
		}

		claims, err := verifier.Verify(ctx, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "token verification failed")
		}

		return handler(authkit.WithClaims(ctx, claims), req)
	}
}

func requestMetadata(ctx context.Context) (string, string) {
	requestID := firstMetadataValue(ctx, HeaderRequestID)
	if requestID == "" {
		requestID = uuid.NewString()
	}

	correlationID := firstMetadataValue(ctx, HeaderCorrelationID)
	if correlationID == "" {
		correlationID = requestID
	}

	return requestID, correlationID
}

func firstMetadataValue(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}
