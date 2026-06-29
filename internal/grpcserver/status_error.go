package grpcserver

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewStatusError(code codes.Code, reason, message string) error {
	st := status.New(code, message)
	withDetails, err := st.WithDetails(&errdetails.ErrorInfo{Reason: reason})
	if err != nil {
		return st.Err()
	}
	return withDetails.Err()
}
