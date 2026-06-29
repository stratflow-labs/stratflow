package auth

import (
	"context"

	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"

	"github.com/google/uuid"
)

// CredentialFinder fetches stored credentials for authentication.
type CredentialFinder interface {
	FindByIdentity(ctx context.Context, identity string) (UserCredentials, error)
}

// UserCredentials holds the minimum data required to authenticate a user.
type UserCredentials struct {
	UserID       uuid.UUID
	PasswordHash string
	Role         string
	Email        string
}

type LogoutCommand struct {
	UserID      uuid.UUID
	AccessToken string
}

// Ports

type AccessTokenRevoker interface {
	DeleteAccessToken(ctx context.Context, userID uuid.UUID, rawToken string) error
}

// PasswordVerifier compares plaintext password with stored hash.
type PasswordVerifier interface {
	Compare(ctx context.Context, plain, hashed string) (bool, error)
}

// TokenService issues and validates opaque access tokens.
type TokenService interface {
	IssueAccessToken(ctx context.Context, claims identitydomain.TokenClaims) (identitydomain.IssuedToken, error)
	VerifyAccessToken(ctx context.Context, token string) (AccessTokenPayload, error)
}

// AccessTokenPayload returned after token verification.
type AccessTokenPayload struct {
	Claims identitydomain.TokenClaims
}

// LoginInput aggregates data required to authenticate user.
type LoginInput struct {
	Login    string
	Password string
}

// LoginOutput returns an issued access token.
type LoginOutput struct {
	AccessToken string
}
