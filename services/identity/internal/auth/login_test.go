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

// ---- fakes ----
type fakeUserCreds struct {
	out usecases.UserCredentials
	err error
}

func (f *fakeUserCreds) FindByIdentity(ctx context.Context, identity string) (usecases.UserCredentials, error) {
	return f.out, f.err
}

type fakeHasher struct {
	match bool
	err   error
}

func (f fakeHasher) Compare(ctx context.Context, plain, hashed string) (bool, error) {
	return f.match, f.err
}

type fakeTokens struct {
	acc    identitydomain.IssuedToken
	errAcc error
}

func (f *fakeTokens) IssueAccessToken(ctx context.Context, claims identitydomain.TokenClaims) (identitydomain.IssuedToken, error) {
	if f == nil {
		return identitydomain.IssuedToken{}, errors.New("fakeTokens not configured")
	}
	return f.acc, f.errAcc
}

func (f *fakeTokens) VerifyAccessToken(ctx context.Context, token string) (usecases.AccessTokenPayload, error) {
	return usecases.AccessTokenPayload{}, errors.New("not used here")
}

type noTx struct{}

func (noTx) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) }

// ---- tests ----
func TestLogin_Success(t *testing.T) {
	userID := uuid.New()
	creds := usecases.UserCredentials{
		UserID:       userID,
		PasswordHash: "hash",
		Role:         "admin",
		Email:        "a@b.c",
	}

	acc := identitydomain.IssuedToken{Value: "a"}

	svc := usecases.NewService(
		&fakeUserCreds{out: creds},
		fakeHasher{match: true},
		&fakeTokens{acc: acc},
		nil,
		noTx{},
	)

	out, err := svc.Login(context.Background(), usecases.LoginInput{
		Login:    "a@b.c",
		Password: "secret",
	})
	require.NoError(t, err)
	require.Equal(t, "a", out.AccessToken)
}

func TestLogin_DoesNotMutateInput(t *testing.T) {
	t.Parallel()

	input := usecases.LoginInput{
		Login:    "  trader  ",
		Password: "  secret  ",
	}

	svc := usecases.NewService(
		&fakeUserCreds{out: usecases.UserCredentials{PasswordHash: "hash", Role: "user"}},
		fakeHasher{match: true},
		&fakeTokens{acc: identitydomain.IssuedToken{Value: "a"}},
		nil,
		noTx{},
	)

	_, err := svc.Login(context.Background(), input)
	require.NoError(t, err)
	require.Equal(t, "  trader  ", input.Login)
	require.Equal(t, "  secret  ", input.Password)
}

func TestLogin_InvalidCredentials_WhenEmpty(t *testing.T) {
	t.Parallel()
	svc := usecases.NewService(
		&fakeUserCreds{},
		fakeHasher{},
		&fakeTokens{},
		nil,
		noTx{},
	)
	_, err := svc.Login(context.Background(), usecases.LoginInput{})
	require.ErrorIs(t, err, identitydomain.ErrInvalidCredentials)
}

func TestLogin_InvalidCredentials_WhenUserNotFound(t *testing.T) {
	t.Parallel()
	svc := usecases.NewService(
		&fakeUserCreds{err: identitydomain.ErrUserNotFound},
		fakeHasher{},
		&fakeTokens{},
		nil,
		noTx{},
	)
	_, err := svc.Login(context.Background(), usecases.LoginInput{
		Login:    "trader",
		Password: "y",
	})
	require.ErrorIs(t, err, identitydomain.ErrUserNotFound)
}

func TestLogin_InvalidCredentials_WhenPasswordMismatch(t *testing.T) {
	t.Parallel()
	svc := usecases.NewService(
		&fakeUserCreds{out: usecases.UserCredentials{PasswordHash: "hash"}},
		fakeHasher{match: false},
		&fakeTokens{},
		nil,
		noTx{},
	)
	_, err := svc.Login(context.Background(), usecases.LoginInput{
		Login:    "trader",
		Password: "y",
	})
	require.ErrorIs(t, err, identitydomain.ErrPasswordMismatch)
}
