package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	usecases "github.com/stratflow-labs/stratflow/services/identity/internal/auth"
	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"
)

type CredentialRepository struct {
	db *sql.DB
}

func NewCredentialRepository(db *sql.DB) *CredentialRepository {
	return &CredentialRepository{db: db}
}

var _ usecases.CredentialFinder = (*CredentialRepository)(nil)

func (r *CredentialRepository) FindByIdentity(ctx context.Context, identity string) (usecases.UserCredentials, error) {
	value := strings.TrimSpace(identity)
	if value == "" {
		return usecases.UserCredentials{}, identitydomain.ErrInvalidCredentials
	}

	const query = `
SELECT
    id AS user_id,
    password AS password_hash,
    role,
    COALESCE(email, '') AS email
FROM users
WHERE LOWER(login) = LOWER($1) OR LOWER(email) = LOWER($1)`

	row := r.db.QueryRowContext(ctx, query, value)
	var creds usecases.UserCredentials
	err := row.Scan(
		&creds.UserID,
		&creds.PasswordHash,
		&creds.Role,
		&creds.Email,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return usecases.UserCredentials{}, identitydomain.ErrUserNotFound
	}
	if err != nil {
		return usecases.UserCredentials{}, fmt.Errorf("find credentials: %w", err)
	}
	return creds, nil
}
