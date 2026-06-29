-- +goose Up
-- Strategy attribute value table.

CREATE TABLE IF NOT EXISTS strategy_attribute_value (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    strategy_attribute_id UUID NOT NULL REFERENCES strategy_attribute(id) ON DELETE CASCADE,
    slug TEXT NOT NULL CHECK (btrim(slug) <> ''),
    value TEXT NOT NULL CHECK (btrim(value) <> ''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_strategy_param_value_param_slug UNIQUE (strategy_attribute_id, slug),
    CHECK (updated_at >= created_at)
);

-- Backward compatibility for pre-existing tables created before timestamps were introduced.
ALTER TABLE strategy_attribute_value ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ;
ALTER TABLE strategy_attribute_value ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;
UPDATE strategy_attribute_value SET created_at = NOW() WHERE created_at IS NULL;
UPDATE strategy_attribute_value SET updated_at = COALESCE(updated_at, created_at, NOW()) WHERE updated_at IS NULL;
ALTER TABLE strategy_attribute_value ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE strategy_attribute_value ALTER COLUMN updated_at SET DEFAULT NOW();
ALTER TABLE strategy_attribute_value ALTER COLUMN created_at SET NOT NULL;
ALTER TABLE strategy_attribute_value ALTER COLUMN updated_at SET NOT NULL;

DROP TRIGGER IF EXISTS trg_strategy_attribute_value_set_updated_at ON strategy_attribute_value;
CREATE TRIGGER trg_strategy_attribute_value_set_updated_at
BEFORE UPDATE ON strategy_attribute_value
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE INDEX IF NOT EXISTS idx_strategy_attribute_value_created_at ON strategy_attribute_value (created_at);
CREATE INDEX IF NOT EXISTS idx_strategy_attribute_value_strategy_attribute_id
ON strategy_attribute_value (strategy_attribute_id);

-- +goose Down
DROP INDEX IF EXISTS idx_strategy_attribute_value_strategy_attribute_id;
DROP INDEX IF EXISTS idx_strategy_attribute_value_created_at;

DROP TRIGGER IF EXISTS trg_strategy_attribute_value_set_updated_at ON strategy_attribute_value;
DROP TABLE IF EXISTS strategy_attribute_value;
