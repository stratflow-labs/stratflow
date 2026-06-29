package testkit

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RequireGRPCCode(t *testing.T, err error, want codes.Code) *status.Status {
	t.Helper()
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok, "expected gRPC status error, got %T", err)
	require.Equal(t, want, st.Code())
	return st
}

func RequireGRPCReason(t *testing.T, err error, want codes.Code, wantReason string) {
	t.Helper()
	st := RequireGRPCCode(t, err, want)
	require.NotEmpty(t, st.Details(), "expected gRPC error details")

	info, ok := st.Details()[0].(*errdetails.ErrorInfo)
	require.True(t, ok, "expected ErrorInfo details, got %T", st.Details()[0])
	require.Equal(t, wantReason, info.Reason)
}
