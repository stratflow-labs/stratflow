-- +goose Up
-- Directed graph of relations between attribute values (many-to-many).

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

-- Optional strictness mode for one-to-one mapping per target attribute.
-- Example: one risk level must map to only one timeframe.
-- CREATE UNIQUE INDEX uq_strategy_param_value_relation_from_target_param
-- ON strategy_attribute_value_relation (from_value_id, to_attribute_id);

CREATE INDEX IF NOT EXISTS idx_strategy_param_value_relation_from_value
ON strategy_attribute_value_relation (from_value_id);

CREATE INDEX IF NOT EXISTS idx_strategy_param_value_relation_to_value
ON strategy_attribute_value_relation (to_value_id);

CREATE INDEX IF NOT EXISTS idx_strategy_param_value_relation_from_to_param
ON strategy_attribute_value_relation (from_attribute_id, to_attribute_id);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION validate_strategy_attribute_value_relation()
RETURNS TRIGGER AS $$
BEGIN
    -- Both attributes must belong to the same strategy.
    PERFORM 1
    FROM strategy_attribute p_from
    JOIN strategy_attribute p_to ON p_from.strategy_id = p_to.strategy_id
    WHERE p_from.id = NEW.from_attribute_id
      AND p_to.id = NEW.to_attribute_id;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'from_attribute_id % and to_attribute_id % must belong to the same strategy',
            NEW.from_attribute_id, NEW.to_attribute_id
            USING ERRCODE = '23503';
    END IF;

    -- Values must belong to their declared attributes.
    PERFORM 1
    FROM strategy_attribute_value pv
    WHERE pv.id = NEW.from_value_id
      AND pv.strategy_attribute_id = NEW.from_attribute_id;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'from_value_id % does not belong to from_attribute_id %', NEW.from_value_id, NEW.from_attribute_id
            USING ERRCODE = '23503';
    END IF;

    PERFORM 1
    FROM strategy_attribute_value pv
    WHERE pv.id = NEW.to_value_id
      AND pv.strategy_attribute_id = NEW.to_attribute_id;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'to_value_id % does not belong to to_attribute_id %', NEW.to_value_id, NEW.to_attribute_id
            USING ERRCODE = '23503';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS trg_strategy_param_value_relation_validate ON strategy_attribute_value_relation;
CREATE TRIGGER trg_strategy_param_value_relation_validate
BEFORE INSERT OR UPDATE ON strategy_attribute_value_relation
FOR EACH ROW
EXECUTE FUNCTION validate_strategy_attribute_value_relation();

DROP TRIGGER IF EXISTS trg_strategy_param_value_relation_set_updated_at ON strategy_attribute_value_relation;
CREATE TRIGGER trg_strategy_param_value_relation_set_updated_at
BEFORE UPDATE ON strategy_attribute_value_relation
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_strategy_param_value_relation_set_updated_at ON strategy_attribute_value_relation;
DROP TRIGGER IF EXISTS trg_strategy_param_value_relation_validate ON strategy_attribute_value_relation;

DROP FUNCTION IF EXISTS validate_strategy_attribute_value_relation();

DROP TABLE IF EXISTS strategy_attribute_value_relation;
