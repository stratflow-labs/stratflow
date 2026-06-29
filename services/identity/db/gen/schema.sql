-- Code generated from db/migrations. DO NOT EDIT.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

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

CREATE INDEX IF NOT EXISTS idx_users_name_trgm
    ON users USING gin (name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_users_last_name_trgm
    ON users USING gin (last_name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_users_email_trgm
    ON users USING gin (email gin_trgm_ops);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_login_lower_unique
    ON users (LOWER(login));

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_lower_unique
    ON users (LOWER(email))
    WHERE email IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_users_created_at_desc
    ON users (created_at DESC);

CREATE TABLE IF NOT EXISTS session (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_session_user_id ON session (user_id);

