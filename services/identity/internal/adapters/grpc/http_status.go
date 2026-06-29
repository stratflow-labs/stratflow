package identitygrpc

import (
	"context"
	"net/http"

	"github.com/stratflow-labs/stratflow/internal/grpcserver"

	"google.golang.org/protobuf/types/known/emptypb"
)

func noContent(ctx context.Context) *emptypb.Empty {
	setHTTPStatus(ctx, http.StatusNoContent)
	return &emptypb.Empty{}
}

func setHTTPStatus(ctx context.Context, statusCode int) {
	grpcserver.SetHTTPStatus(ctx, statusCode)
}
