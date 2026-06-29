package authverifier

import (
	"context"

	"github.com/stratflow-labs/stratflow/internal/authkit"
	identityclient "github.com/stratflow-labs/stratflow/services/identity/client"
)

// VerifyTokenAdapter bridges Identity client TokenVerifier to authkit.AccessTokenVerifier.
type VerifyTokenAdapter struct {
	Verifier identityclient.TokenVerifier
}

func (a VerifyTokenAdapter) Verify(ctx context.Context, token string) (authkit.Claims, error) {
	if a.Verifier == nil {
		return authkit.Claims{}, identityclient.ErrUnauthorized
	}
	out, err := a.Verifier.Verify(ctx, token)
	if err != nil {
		return authkit.Claims{}, err
	}
	return authkit.Claims{UserID: out.UserID.String(), Role: out.Role}, nil
}
