package grpcserver

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const MetadataHTTPStatus = "x-http-code"

func SetHTTPStatus(ctx context.Context, statusCode int) {
	if ctx == nil || statusCode == 0 {
		return
	}

	_ = grpc.SetHeader(ctx, metadata.Pairs(MetadataHTTPStatus, strconv.Itoa(statusCode)))
}
