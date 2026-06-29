-- Code generated from db/migrations. DO NOT EDIT.

CREATE TABLE IF NOT EXISTS strategy (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL UNIQUE CHECK (btrim(slug) <> ''),
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    description TEXT NOT NULL CHECK (btrim(description) <> ''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (updated_at >= created_at)
);

ALTER TABLE strategy ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ;

ALTER TABLE strategy ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;

ALTER TABLE strategy ALTER COLUMN created_at SET DEFAULT NOW();

ALTER TABLE strategy ALTER COLUMN updated_at SET DEFAULT NOW();

ALTER TABLE strategy ALTER COLUMN created_at SET NOT NULL;

ALTER TABLE strategy ALTER COLUMN updated_at SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_strategy_created_at ON strategy (created_at);

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

ALTER TABLE strategy_attribute ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ;

ALTER TABLE strategy_attribute ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;

ALTER TABLE strategy_attribute ALTER COLUMN created_at SET DEFAULT NOW();

ALTER TABLE strategy_attribute ALTER COLUMN updated_at SET DEFAULT NOW();

ALTER TABLE strategy_attribute ALTER COLUMN created_at SET NOT NULL;

ALTER TABLE strategy_attribute ALTER COLUMN updated_at SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_strategy_attribute_created_at 
ON strategy_attribute (created_at);

CREATE INDEX IF NOT EXISTS idx_strategy_attribute_strategy_id
ON strategy_attribute (strategy_id);

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

ALTER TABLE strategy_attribute_value ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ;

ALTER TABLE strategy_attribute_value ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;

ALTER TABLE strategy_attribute_value ALTER COLUMN created_at SET DEFAULT NOW();

ALTER TABLE strategy_attribute_value ALTER COLUMN updated_at SET DEFAULT NOW();

ALTER TABLE strategy_attribute_value ALTER COLUMN created_at SET NOT NULL;

ALTER TABLE strategy_attribute_value ALTER COLUMN updated_at SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_strategy_attribute_value_created_at ON strategy_attribute_value (created_at);

CREATE INDEX IF NOT EXISTS idx_strategy_attribute_value_strategy_attribute_id
ON strategy_attribute_value (strategy_attribute_id);

CREATE TABLE IF NOT EXISTS strategy_attribute_value_relation (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_attribute_id UUID NOT NULL REFERENCES strategy_attribute(id) ON DELETE CASCADE,
    from_value_id UUID NOT NULL REFERENCES strategy_attribute_value(id) ON DELETE CASCADE,
    to_attribute_id UUID NOT NULL REFERENCES strategy_attribute(id) ON DELETE CASCADE,
    to_value_id UUID NOT NULL REFERENCES strategy_attribute_value(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_strategy_param_value_relation_edge UNIQUE (from_value_id, to_value_id),
    CHECK (from_attribute_id <> to_attribute_id),
    CHECK (from_value_id <> to_value_id),
    CHECK (updated_at >= created_at)
);

CREATE INDEX IF NOT EXISTS idx_strategy_param_value_relation_from_value
ON strategy_attribute_value_relation (from_value_id);

CREATE INDEX IF NOT EXISTS idx_strategy_param_value_relation_to_value
ON strategy_attribute_value_relation (to_value_id);

CREATE INDEX IF NOT EXISTS idx_strategy_param_value_relation_from_to_param
ON strategy_attribute_value_relation (from_attribute_id, to_attribute_id);

