package authverifier

import (
	"context"

	"github.com/stratflow-labs/stratflow/internal/authkit"
	auth "github.com/stratflow-labs/stratflow/services/identity/internal/auth"
)

type VerifyTokenAdapter struct {
	Auth *auth.Service
}

var _ authkit.AccessTokenVerifier = (*VerifyTokenAdapter)(nil)

func (v *VerifyTokenAdapter) Verify(ctx context.Context, token string) (authkit.Claims, error) {
	out, err := v.Auth.VerifyToken(ctx, &auth.VerifyTokenInput{AccessToken: token})
	if err != nil {
		return authkit.Claims{}, err
	}
	return authkit.Claims{
		UserID: out.Claims.UserID.String(),
		Role:   out.Claims.Role,
	}, nil
}
