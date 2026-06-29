-- +goose Up
-- Opaque access tokens table. Only HMAC hashes are stored; raw tokens are returned to clients once.
CREATE TABLE IF NOT EXISTS session (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_session_user_id ON session (user_id);

-- +goose Down
DROP TABLE IF EXISTS session;
