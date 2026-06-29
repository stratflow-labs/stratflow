-- +goose Up
-- +goose NO TRANSACTION
-- pg_trgm extension for efficient substring search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    login TEXT NOT NULL UNIQUE CHECK (
        login = LOWER(login)
        AND char_length(login) BETWEEN 3 AND 32
        AND login ~ '^[A-Za-z0-9]+([._-][A-Za-z0-9]+)*$'
    ),
    email TEXT UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL,
    image_url TEXT NOT NULL DEFAULT '',
    gender TEXT DEFAULT NULL CHECK (gender IN ('male', 'female')),
    is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DROP TRIGGER IF EXISTS trg_users_set_updated_at ON users;
CREATE TRIGGER trg_users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_name_trgm
    ON users USING gin (name gin_trgm_ops);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_last_name_trgm
    ON users USING gin (last_name gin_trgm_ops);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email_trgm
    ON users USING gin (email gin_trgm_ops);
CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_users_login_lower_unique
    ON users (LOWER(login));
CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email_lower_unique
    ON users (LOWER(email))
    WHERE email IS NOT NULL;

-- Indexes for sorting by created_at
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_created_at_desc
    ON users (created_at DESC);

-- +goose Down
DROP TRIGGER IF EXISTS trg_users_set_updated_at ON users;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS pg_trgm;
