-- +goose Up
-- Seed attribute value relations (directed edges) for local/dev.
-- Mapping model: risk_level -> timeframe.

WITH seed_relations (strategy_slug, from_attribute_slug, from_value_slug, to_attribute_slug, to_value_slug) AS (
    VALUES
        ('mean-reversion', 'risk_level', 'low', 'timeframe', 'tf_1m'),
        ('mean-reversion', 'risk_level', 'medium', 'timeframe', 'tf_5m'),
        ('mean-reversion', 'risk_level', 'high', 'timeframe', 'tf_15m'),
        ('mean-reversion', 'risk_level', 'very_high', 'timeframe', 'tf_30m'),

        ('trend-momentum', 'risk_level', 'low', 'timeframe', 'tf_1m'),
        ('trend-momentum', 'risk_level', 'medium', 'timeframe', 'tf_5m'),
        ('trend-momentum', 'risk_level', 'high', 'timeframe', 'tf_15m'),
        ('trend-momentum', 'risk_level', 'very_high', 'timeframe', 'tf_30m'),

        ('volatility-breakout', 'risk_level', 'low', 'timeframe', 'tf_1m'),
        ('volatility-breakout', 'risk_level', 'medium', 'timeframe', 'tf_5m'),
        ('volatility-breakout', 'risk_level', 'high', 'timeframe', 'tf_15m'),
        ('volatility-breakout', 'risk_level', 'very_high', 'timeframe', 'tf_30m')
)
INSERT INTO strategy_attribute_value_relation (
    id,
    from_attribute_id,
    from_value_id,
    to_attribute_id,
    to_value_id,
    created_at,
    updated_at
)
SELECT
    gen_random_uuid(),
    sp_from.id,
    pv_from.id,
    sp_to.id,
    pv_to.id,
    NOW(),
    NOW()
FROM seed_relations r
JOIN strategy s ON s.slug = r.strategy_slug
JOIN strategy_attribute sp_from
    ON sp_from.strategy_id = s.id
   AND sp_from.slug = r.from_attribute_slug
JOIN strategy_attribute sp_to
    ON sp_to.strategy_id = s.id
   AND sp_to.slug = r.to_attribute_slug
JOIN strategy_attribute_value pv_from
    ON pv_from.strategy_attribute_id = sp_from.id
   AND pv_from.slug = r.from_value_slug
JOIN strategy_attribute_value pv_to
    ON pv_to.strategy_attribute_id = sp_to.id
   AND pv_to.slug = r.to_value_slug
ON CONFLICT ON CONSTRAINT uq_strategy_param_value_relation_edge DO NOTHING;

-- +goose Down
DELETE FROM strategy_attribute_value_relation rel
USING strategy s,
      strategy_attribute sp_from,
      strategy_attribute sp_to,
      strategy_attribute_value pv_from,
      strategy_attribute_value pv_to
WHERE rel.from_attribute_id = sp_from.id
  AND rel.to_attribute_id = sp_to.id
  AND rel.from_value_id = pv_from.id
  AND rel.to_value_id = pv_to.id
  AND sp_from.strategy_id = s.id
  AND sp_to.strategy_id = s.id
  AND pv_from.strategy_attribute_id = sp_from.id
  AND pv_to.strategy_attribute_id = sp_to.id
  AND (s.slug, sp_from.slug, pv_from.slug, sp_to.slug, pv_to.slug) IN (
    ('mean-reversion', 'risk_level', 'low', 'timeframe', 'tf_1m'),
    ('mean-reversion', 'risk_level', 'medium', 'timeframe', 'tf_5m'),
    ('mean-reversion', 'risk_level', 'high', 'timeframe', 'tf_15m'),
    ('mean-reversion', 'risk_level', 'very_high', 'timeframe', 'tf_30m'),
    ('trend-momentum', 'risk_level', 'low', 'timeframe', 'tf_1m'),
    ('trend-momentum', 'risk_level', 'medium', 'timeframe', 'tf_5m'),
    ('trend-momentum', 'risk_level', 'high', 'timeframe', 'tf_15m'),
    ('trend-momentum', 'risk_level', 'very_high', 'timeframe', 'tf_30m'),
    ('volatility-breakout', 'risk_level', 'low', 'timeframe', 'tf_1m'),
    ('volatility-breakout', 'risk_level', 'medium', 'timeframe', 'tf_5m'),
    ('volatility-breakout', 'risk_level', 'high', 'timeframe', 'tf_15m'),
    ('volatility-breakout', 'risk_level', 'very_high', 'timeframe', 'tf_30m')
  );
