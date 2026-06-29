package identitygrpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stratflow-labs/stratflow/internal/authkit"
	"github.com/stratflow-labs/stratflow/services/identity/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestMapError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		err     error
		code    codes.Code
		reason  string
		message string
	}{
		{
			name:    "wrapped invalid credentials",
			err:     errors.Join(errors.New("wrapped"), domain.ErrInvalidCredentials),
			code:    codes.Unauthenticated,
			reason:  "auth.invalidCredentials",
			message: "invalid email or password",
		},
		{
			name:    "user not found",
			err:     domain.ErrUserNotFound,
			code:    codes.NotFound,
			reason:  "user.notFound",
			message: "user not found",
		},
		{
			name:    "fallback internal",
			err:     errors.New("boom"),
			code:    codes.Internal,
			reason:  "identity.internal",
			message: "Internal Server Error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			st := status.Convert(mapError(tc.err))
			require.Equal(t, tc.code, st.Code())
			require.Equal(t, tc.message, st.Message())

			require.Len(t, st.Details(), 1)
			info, ok := st.Details()[0].(*errdetails.ErrorInfo)
			require.True(t, ok)
			require.Equal(t, tc.reason, info.Reason)
		})
	}
}

func TestContextHelpers(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	ctx := authkit.WithClaims(context.Background(), authkit.Claims{UserID: userID.String()})
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("authorization", "Bearer  token-123  "))

	gotUserID, err := currentUserID(ctx)
	require.NoError(t, err)
	require.Equal(t, userID, gotUserID)

	token, err := bearerToken(ctx)
	require.NoError(t, err)
	require.Equal(t, "token-123", token)
}

func TestContextHelpers_ReturnDomainErrors(t *testing.T) {
	t.Parallel()

	_, err := currentUserID(context.Background())
	require.ErrorIs(t, err, domain.ErrAccessTokenNotFound)

	_, err = currentUserID(authkit.WithClaims(context.Background(), authkit.Claims{UserID: "not-a-uuid"}))
	require.ErrorIs(t, err, domain.ErrUserNotFound)

	_, err = bearerToken(context.Background())
	require.ErrorIs(t, err, domain.ErrAccessTokenNotFound)
}

func TestMappingHelpers(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 18, 12, 0, 0, 0, time.UTC)
	user := domain.User{
		ID:              uuid.New(),
		Login:           "trader.one",
		Name:            "Alice",
		LastName:        "Smith",
		Email:           " alice@example.com ",
		Role:            "admin",
		ImageUrl:        " https://cdn.example/avatar.png ",
		Gender:          2,
		IsEmailVerified: true,
		CreatedAt:       now,
		UpdatedAt:       now.Add(time.Minute),
	}

	mapped := mapUser(user)
	require.Equal(t, user.ID.String(), mapped.Id)
	require.Equal(t, "alice@example.com", mapped.GetEmail())
	require.Equal(t, "https://cdn.example/avatar.png", mapped.GetImageUrl())
	require.NotNil(t, mapped.Gender)
	require.EqualValues(t, 2, *mapped.Gender)
	require.True(t, mapped.IsVerified)
	require.True(t, mapped.IsEmailVerified)
	require.Equal(t, now, mapped.CreatedAt.AsTime())
	require.Equal(t, now.Add(time.Minute), mapped.UpdatedAt.AsTime())

	require.Equal(t, int32(3), *mapGender(3))
	require.Nil(t, mapGender(0))
	require.Equal(t, "value", *optionalString("  value  "))
	require.Nil(t, optionalString(" \t "))
	require.NotNil(t, updateUserPatch(nil))

	gender := int32(7)
	require.Equal(t, 7, *int32PtrToIntPtr(&gender))
	require.Nil(t, int32PtrToIntPtr(nil))
	require.Equal(t, "ROLE_ADMIN", mapRole("admin").String())
	require.Equal(t, "ROLE_UNSPECIFIED", mapRole("unknown").String())
}

func TestUnimplemented(t *testing.T) {
	t.Parallel()

	st := status.Convert(unimplemented("DeleteEverything"))
	require.Equal(t, codes.Unimplemented, st.Code())
	require.Equal(t, "identity grpc: DeleteEverything is not implemented yet", st.Message())
}
