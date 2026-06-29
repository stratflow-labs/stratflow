package authverifier

import (
	"context"
	"strings"

	"github.com/stratflow-labs/stratflow/internal/authkit"
	identityclient "github.com/stratflow-labs/stratflow/services/identity/client"
)

// ServiceTokenFallbackVerifier authorizes requests by first matching a service token,
// then delegating to an Identity-backed verifier for user tokens.
type ServiceTokenFallbackVerifier struct {
	ServiceToken     string
	IdentityVerifier authkit.AccessTokenVerifier
}

func (v ServiceTokenFallbackVerifier) Verify(ctx context.Context, token string) (authkit.Claims, error) {
	tok := strings.TrimSpace(token)
	if tok == "" {
		return authkit.Claims{}, identityclient.ErrUnauthorized
	}

	if v.ServiceToken != "" && tok == v.ServiceToken {
		return authkit.Claims{UserID: "strategy-registry", Role: "admin"}, nil
	}

	if v.IdentityVerifier != nil {
		return v.IdentityVerifier.Verify(ctx, tok)
	}

	return authkit.Claims{}, identityclient.ErrUnauthorized
}
