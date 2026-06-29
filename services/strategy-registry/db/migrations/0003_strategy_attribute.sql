-- +goose Up
-- Attribute table.

CREATE TABLE IF NOT EXISTS strategy_attribute (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    strategy_id UUID NOT NULL REFERENCES strategy(id) ON DELETE CASCADE,
    slug TEXT NOT NULL CHECK (btrim(slug) <> ''),
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    description TEXT NOT NULL CHECK (btrim(description) <> ''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_strategy_param_strategy_slug UNIQUE (strategy_id, slug),
    CHECK (updated_at >= created_at)
);

-- Backward compatibility for pre-existing tables created before timestamps were introduced.
ALTER TABLE strategy_attribute ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ;
ALTER TABLE strategy_attribute ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;
UPDATE strategy_attribute SET created_at = NOW() WHERE created_at IS NULL;
UPDATE strategy_attribute SET updated_at = COALESCE(updated_at, created_at, NOW()) WHERE updated_at IS NULL;
ALTER TABLE strategy_attribute ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE strategy_attribute ALTER COLUMN updated_at SET DEFAULT NOW();
ALTER TABLE strategy_attribute ALTER COLUMN created_at SET NOT NULL;
ALTER TABLE strategy_attribute ALTER COLUMN updated_at SET NOT NULL;


DROP TRIGGER IF EXISTS trg_strategy_attribute_set_updated_at ON strategy_attribute;
CREATE TRIGGER trg_strategy_attribute_set_updated_at
BEFORE UPDATE ON strategy_attribute
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE INDEX IF NOT EXISTS idx_strategy_attribute_created_at 
ON strategy_attribute (created_at);
CREATE INDEX IF NOT EXISTS idx_strategy_attribute_strategy_id
ON strategy_attribute (strategy_id);

-- +goose Down
DROP INDEX IF EXISTS idx_strategy_attribute_strategy_id;
DROP INDEX IF EXISTS idx_strategy_attribute_created_at;

DROP TRIGGER IF EXISTS trg_strategy_attribute_set_updated_at ON strategy_attribute;
DROP TABLE IF EXISTS strategy_attribute;
