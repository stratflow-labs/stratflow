package auth_test

import (
	"context"
	"errors"
	"testing"

	usecases "github.com/stratflow-labs/stratflow/services/identity/internal/auth"
	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type tokSvc struct {
	out usecases.AccessTokenPayload
	err error
}

func (t tokSvc) IssueAccessToken(context.Context, identitydomain.TokenClaims) (identitydomain.IssuedToken, error) {
	return identitydomain.IssuedToken{}, errors.New("IssueAccessToken must not be called in verify token tests")
}

func (t tokSvc) VerifyAccessToken(ctx context.Context, token string) (usecases.AccessTokenPayload, error) {
	return t.out, t.err
}

func TestVerify_Success(t *testing.T) {
	claims := identitydomain.TokenClaims{UserID: uuid.New(), Role: "user"}
	svc := usecases.NewService(nil, nil, tokSvc{
		out: usecases.AccessTokenPayload{Claims: claims},
	}, nil, nil)
	out, err := svc.VerifyToken(context.Background(), &usecases.VerifyTokenInput{AccessToken: "x"})
	require.NoError(t, err)
	require.Equal(t, claims.Role, out.Claims.Role)
}
