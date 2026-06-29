package auth_test

import (
	"context"
	"testing"

	usecases "github.com/stratflow-labs/stratflow/services/identity/internal/auth"
	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestLogoutSession_Success(t *testing.T) {
	userID := uuid.New()
	token := "refresh-token"

	access := &fakeAccessTokens{}
	svc := usecases.NewService(nil, nil, nil, access, &fakeTx{})
	err := svc.Logout(context.Background(), usecases.LogoutCommand{UserID: userID, AccessToken: token})
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{userID}, access.deletedUsers)
	require.Equal(t, []string{token}, access.deletedTokens)
}

func TestLogoutSession_InvalidID(t *testing.T) {
	svc := usecases.NewService(nil, nil, nil, &fakeAccessTokens{}, &fakeTx{})
	err := svc.Logout(context.Background(), usecases.LogoutCommand{UserID: uuid.Nil, AccessToken: "token"})
	require.ErrorIs(t, err, identitydomain.ErrAccessTokenNotFound)
}

func TestLogoutSession_EmptyToken(t *testing.T) {
	svc := usecases.NewService(nil, nil, nil, &fakeAccessTokens{}, &fakeTx{})
	err := svc.Logout(context.Background(), usecases.LogoutCommand{UserID: uuid.New(), AccessToken: "   "})
	require.ErrorIs(t, err, identitydomain.ErrAccessTokenNotFound)
}
