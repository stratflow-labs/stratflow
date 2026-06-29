package strategyhttp

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	apphttp "github.com/stratflow-labs/stratflow/internal/app/http"
	"github.com/stratflow-labs/stratflow/internal/foundation/httpx/respond"
	"github.com/stratflow-labs/stratflow/internal/grpcserver"
	"github.com/stratflow-labs/stratflow/internal/httpserver"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type GatewayHandler struct {
	handler http.Handler
}

func NewGatewayHandler(server strategyregistryv1.StrategyRegistryServiceServer) *GatewayHandler {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(incomingHeaderMatcher),
		runtime.WithErrorHandler(handleGatewayError),
		runtime.WithForwardResponseOption(applyHTTPStatus),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   false,
				EmitUnpopulated: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	return &GatewayHandler{handler: mux}
}

func (h *GatewayHandler) RegisterRoutes(r httpserver.Router) {
	r.Mount("/api", h.handler)
	r.Mount("/api/", h.handler)
}

func incomingHeaderMatcher(key string) (string, bool) {
	if strings.EqualFold(key, "Authorization") {
		return "authorization", true
	}

	return runtime.DefaultHeaderMatcher(key)
}

func applyHTTPStatus(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	values := md.HeaderMD.Get(grpcserver.MetadataHTTPStatus)
	if len(values) == 0 {
		return nil
	}

	delete(md.HeaderMD, grpcserver.MetadataHTTPStatus)
	w.Header().Del("Grpc-Metadata-X-Http-Code")

	statusCode, convErr := strconv.Atoi(values[0])
	if convErr == nil && statusCode >= 100 && statusCode <= 999 {
		w.WriteHeader(statusCode)
	}
	return nil
}

func handleGatewayError(ctx context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	st := status.Convert(err)
	statusCode := runtime.HTTPStatusFromCode(st.Code())
	respond.Problem(w, statusCode, gatewayErrorCode(st), gatewayErrorMessage(statusCode, st), nil)
}

func gatewayErrorCode(st *status.Status) string {
	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok && strings.TrimSpace(info.Reason) != "" {
			return info.Reason
		}
	}

	switch st.Code() {
	case codes.OK,
		codes.Canceled,
		codes.Unknown,
		codes.DeadlineExceeded,
		codes.AlreadyExists,
		codes.ResourceExhausted,
		codes.FailedPrecondition,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unimplemented,
		codes.Internal,
		codes.Unavailable,
		codes.DataLoss:
		return "strategyRegistry.internal"
	case codes.InvalidArgument:
		return "request.invalid"
	case codes.Unauthenticated:
		return "session.unauthorized"
	case codes.PermissionDenied:
		return "session.forbidden"
	case codes.NotFound:
		return "strategyRegistry.notFound"
	}

	return "strategyRegistry.internal"
}

func gatewayErrorMessage(statusCode int, st *status.Status) string {
	message := strings.TrimSpace(st.Message())
	if message == "" {
		message = http.StatusText(statusCode)
	}

	if statusCode == http.StatusInternalServerError {
		return apphttp.FinalizeErrorMessage(statusCode, "internal server error", st.Err())
	}

	return message
}

func SetHTTPStatus(ctx context.Context, statusCode int) {
	grpcserver.SetHTTPStatus(ctx, statusCode)
}
