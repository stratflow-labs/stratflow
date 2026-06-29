package grpcserver

import "context"

type (
	requestIDKey     struct{}
	correlationIDKey struct{}
)

func WithRequestMetadata(ctx context.Context, requestID, correlationID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = context.WithValue(ctx, requestIDKey{}, requestID)
	ctx = context.WithValue(ctx, correlationIDKey{}, correlationID)
	return ctx
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	value, ok := ctx.Value(requestIDKey{}).(string)
	return value, ok
}

func CorrelationIDFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	value, ok := ctx.Value(correlationIDKey{}).(string)
	return value, ok
}
