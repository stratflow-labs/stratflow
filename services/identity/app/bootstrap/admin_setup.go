package bootstrap

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	appcli "github.com/stratflow-labs/stratflow/internal/app/cli"
	appconfig "github.com/stratflow-labs/stratflow/internal/app/config"
	"github.com/stratflow-labs/stratflow/internal/foundation/password"
	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"

	"github.com/google/uuid"
)

// AdminSetupInput contains admin fields for create/update flow.
type AdminSetupInput struct {
	Login    string
	Email    string
	Name     string
	LastName string
	Role     string
	Gender   int
	Password string
}

// EnsureAdmin creates or updates admin user in identity DB.
func EnsureAdmin(ctx context.Context, input AdminSetupInput) (uuid.UUID, bool, error) {
	normalized := normalizeAdminInput(input)
	login, err := identitydomain.NormalizeUsername(normalized.Login)
	if err != nil {
		return uuid.Nil, false, err
	}
	normalized.Login = login

	hashedPassword, err := hashAdminPassword(ctx, normalized.Password)
	if err != nil {
		return uuid.Nil, false, err
	}

	return upsertAdmin(ctx, normalized, hashedPassword)
}

func normalizeAdminInput(input AdminSetupInput) AdminSetupInput {
	input.Login = strings.TrimSpace(input.Login)
	input.Email = strings.TrimSpace(input.Email)
	input.Name = strings.TrimSpace(input.Name)
	input.LastName = strings.TrimSpace(input.LastName)
	input.Role = strings.TrimSpace(input.Role)
	input.Password = strings.TrimSpace(input.Password)
	return input
}

func genderToDBValue(code int) string {
	if code == 2 {
		return "female"
	}
	return "male"
}

func hashAdminPassword(ctx context.Context, plain string) (string, error) {
	hasher := password.NewBcryptHasher(0)
	hashed, err := hasher.Hash(ctx, plain)
	if err != nil {
		return "", fmt.Errorf("hash admin password: %w", err)
	}
	return hashed, nil
}

func upsertAdmin(ctx context.Context, input AdminSetupInput, hashedPassword string) (uuid.UUID, bool, error) {
	var adminID uuid.UUID
	created := false

	err := appcli.WithConfigAndDB(ctx, LoadConfig, OpenDB, func(ctx context.Context, _ *appconfig.Config, dbConn *sql.DB) error {
		tx, err := dbConn.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin transaction: %w", err)
		}
		defer func() {
			_ = tx.Rollback()
		}()

		emailID, hasEmail, err := findUserIDByEmail(ctx, tx, input.Email)
		if err != nil {
			return err
		}

		if hasEmail {
			adminID = emailID
		} else {
			adminID = uuid.New()
			created = true
		}

		if created {
			if _, err := tx.ExecContext(ctx, `
INSERT INTO users (
	id,
	name,
	last_name,
	login,
	email,
	password,
	role,
	image_url,
	gender,
	is_email_verified,
	created_at,
	updated_at
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, '', $8, TRUE, NOW(), NOW()
)`,
				adminID,
				input.Name,
				input.LastName,
				input.Login,
				input.Email,
				hashedPassword,
				input.Role,
				genderToDBValue(input.Gender),
			); err != nil {
				return fmt.Errorf("insert admin: %w", err)
			}
		} else {
			if _, err := tx.ExecContext(ctx, `
UPDATE users
SET
	name = $2,
	last_name = $3,
	email = $4,
	login = $5,
	password = $6,
	role = $7,
	image_url = '',
	gender = $8,
	is_email_verified = TRUE,
	updated_at = NOW()
WHERE id = $1`,
				adminID,
				input.Name,
				input.LastName,
				input.Email,
				input.Login,
				hashedPassword,
				input.Role,
				genderToDBValue(input.Gender),
			); err != nil {
				return fmt.Errorf("update admin: %w", err)
			}
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit transaction: %w", err)
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, false, err
	}

	return adminID, created, nil
}

func findUserIDByEmail(ctx context.Context, tx *sql.Tx, email string) (uuid.UUID, bool, error) {
	var id uuid.UUID
	err := tx.QueryRowContext(
		ctx,
		`SELECT id FROM users WHERE LOWER(email) = LOWER($1) ORDER BY created_at ASC LIMIT 1 FOR UPDATE`,
		email,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, false, nil
	}
	if err != nil {
		return uuid.Nil, false, fmt.Errorf("find user by email: %w", err)
	}
	return id, true, nil
}
