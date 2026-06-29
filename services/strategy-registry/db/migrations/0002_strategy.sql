-- +goose Up
-- Strategy table.

CREATE TABLE IF NOT EXISTS strategy (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL UNIQUE CHECK (btrim(slug) <> ''),
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    description TEXT NOT NULL CHECK (btrim(description) <> ''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (updated_at >= created_at)
);

-- Backward compatibility for pre-existing tables created before timestamps were introduced.
ALTER TABLE strategy ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ;
ALTER TABLE strategy ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;
UPDATE strategy SET created_at = NOW() WHERE created_at IS NULL;
UPDATE strategy SET updated_at = COALESCE(updated_at, created_at, NOW()) WHERE updated_at IS NULL;
ALTER TABLE strategy ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE strategy ALTER COLUMN updated_at SET DEFAULT NOW();
ALTER TABLE strategy ALTER COLUMN created_at SET NOT NULL;
ALTER TABLE strategy ALTER COLUMN updated_at SET NOT NULL;

DROP TRIGGER IF EXISTS trg_strategy_set_updated_at ON strategy;
CREATE TRIGGER trg_strategy_set_updated_at
BEFORE UPDATE ON strategy
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE INDEX IF NOT EXISTS idx_strategy_created_at ON strategy (created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_strategy_created_at;

DROP TRIGGER IF EXISTS trg_strategy_set_updated_at ON strategy;
DROP TABLE IF EXISTS strategy;
