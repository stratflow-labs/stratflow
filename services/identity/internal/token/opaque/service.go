package opaque

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/stratflow-labs/stratflow/internal/foundation/crypto/tokenhash"
	"github.com/stratflow-labs/stratflow/internal/foundation/tx/sqltx"
	auth "github.com/stratflow-labs/stratflow/services/identity/internal/auth"
	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"

	"github.com/google/uuid"
)

type Config struct {
	TokenHashSecret string
	TokenBytes      int
}

type Service struct {
	db  *sql.DB
	cfg Config
}

type accessTokenRow struct {
	UserID    uuid.UUID
	TokenHash string
	UserRole  string
}

type dbRunner interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

var (
	_ auth.TokenService       = (*Service)(nil)
	_ auth.AccessTokenRevoker = (*Service)(nil)
)

func NewService(db *sql.DB, cfg *Config) (*Service, error) {
	if db == nil {
		return nil, errors.New("db is required")
	}
	if cfg == nil {
		return nil, errors.New("config is required")
	}

	normalized := *cfg
	if strings.TrimSpace(normalized.TokenHashSecret) == "" {
		return nil, errors.New("token hash secret is required")
	}
	if normalized.TokenBytes <= 0 {
		normalized.TokenBytes = 32
	}

	return &Service{db: db, cfg: normalized}, nil
}

func (s *Service) IssueAccessToken(ctx context.Context, claims identitydomain.TokenClaims) (identitydomain.IssuedToken, error) {
	rawToken, err := generateToken(s.cfg.TokenBytes)
	if err != nil {
		return identitydomain.IssuedToken{}, fmt.Errorf("generate session token: %w", err)
	}
	hash, err := s.hashToken(rawToken)
	if err != nil {
		return identitydomain.IssuedToken{}, fmt.Errorf("hash session token: %w", err)
	}

	const query = `
	INSERT INTO session (
	    id,
	    user_id,
	    token_hash,
	    created_at
	) VALUES (
	    $1, $2, $3, NOW()
	)`
	if _, err := s.runner(ctx).ExecContext(ctx, query, uuid.New(), claims.UserID, hash); err != nil {
		return identitydomain.IssuedToken{}, fmt.Errorf("store session token: %w", err)
	}

	return identitydomain.IssuedToken{Value: rawToken}, nil
}

func (s *Service) VerifyAccessToken(ctx context.Context, rawToken string) (auth.AccessTokenPayload, error) {
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return auth.AccessTokenPayload{}, identitydomain.ErrAccessTokenInvalid
	}

	hash, err := s.hashToken(rawToken)
	if err != nil {
		return auth.AccessTokenPayload{}, identitydomain.ErrAccessTokenInvalid
	}

	row, err := s.findAccessToken(ctx, hash)
	if errors.Is(err, sql.ErrNoRows) {
		return auth.AccessTokenPayload{}, identitydomain.ErrAccessTokenInvalid
	}
	if err != nil {
		return auth.AccessTokenPayload{}, fmt.Errorf("find access token: %w", err)
	}
	if !tokenhash.Verify(rawToken, row.TokenHash, []byte(s.cfg.TokenHashSecret)) {
		return auth.AccessTokenPayload{}, identitydomain.ErrAccessTokenInvalid
	}

	return auth.AccessTokenPayload{
		Claims: identitydomain.TokenClaims{
			UserID: row.UserID,
			Role:   row.UserRole,
		},
	}, nil
}

func (s *Service) DeleteAccessTokensByUser(ctx context.Context, userID uuid.UUID) error {
	const query = `DELETE FROM session WHERE user_id = $1`
	_, err := s.runner(ctx).ExecContext(ctx, query, userID)
	return err
}

func (s *Service) DeleteAccessToken(ctx context.Context, userID uuid.UUID, rawToken string) error {
	if userID == uuid.Nil {
		return identitydomain.ErrAccessTokenNotFound
	}
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return identitydomain.ErrAccessTokenNotFound
	}
	hash, err := s.hashToken(rawToken)
	if err != nil {
		return fmt.Errorf("hash session token: %w", err)
	}
	return s.DeleteAccessTokenByHashAndUser(ctx, hash, userID)
}

func (s *Service) DeleteAccessTokenByHashAndUser(ctx context.Context, hash string, userID uuid.UUID) error {
	const query = `DELETE FROM session WHERE token_hash = $1 AND user_id = $2`
	result, err := s.runner(ctx).ExecContext(ctx, query, hash, userID)
	if err != nil {
		return err
	}
	return ensureDeleted(result)
}

func (s *Service) DeleteAccessTokenByHash(ctx context.Context, hash string) error {
	const query = `DELETE FROM session WHERE token_hash = $1`
	result, err := s.runner(ctx).ExecContext(ctx, query, hash)
	if err != nil {
		return err
	}
	return ensureDeleted(result)
}

func ensureDeleted(result sql.Result) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return identitydomain.ErrAccessTokenNotFound
	}
	return nil
}

func (s *Service) findAccessToken(ctx context.Context, hash string) (accessTokenRow, error) {
	const query = `
SELECT s.user_id, s.token_hash, u.role
FROM session s
JOIN users u ON u.id = s.user_id
WHERE s.token_hash = $1`

	var row accessTokenRow
	err := s.runner(ctx).QueryRowContext(ctx, query, hash).Scan(
		&row.UserID,
		&row.TokenHash,
		&row.UserRole,
	)
	return row, err
}

func (s *Service) hashToken(rawToken string) (string, error) {
	return tokenhash.Hash(rawToken, []byte(s.cfg.TokenHashSecret))
}

func (s *Service) runner(ctx context.Context) dbRunner {
	if tx := sqltx.FromCtx(ctx); tx != nil {
		return tx
	}
	return s.db
}

func generateToken(size int) (string, error) {
	randomBytes := make([]byte, size)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}
