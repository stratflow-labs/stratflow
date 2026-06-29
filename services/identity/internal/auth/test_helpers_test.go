package auth_test

import (
	"context"

	"github.com/google/uuid"
)

type fakeAccessTokens struct {
	deletedUsers  []uuid.UUID
	deletedTokens []string
	err           error
}

func (a *fakeAccessTokens) DeleteAccessToken(ctx context.Context, userID uuid.UUID, token string) error {
	if a.err != nil {
		return a.err
	}
	a.deletedUsers = append(a.deletedUsers, userID)
	a.deletedTokens = append(a.deletedTokens, token)
	return nil
}

type fakeTx struct{ calls int }

func (t *fakeTx) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	t.calls++
	return fn(ctx)
}
