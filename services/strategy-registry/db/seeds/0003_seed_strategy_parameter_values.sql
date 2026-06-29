-- +goose Up
-- Seed strategy attribute values for local/dev.

WITH seed_values (strategy_slug, attribute_slug, slug, value) AS (
    VALUES
        ('mean-reversion', 'risk_level', 'low', 'low'),
        ('mean-reversion', 'risk_level', 'medium', 'medium'),
        ('mean-reversion', 'risk_level', 'high', 'high'),
        ('mean-reversion', 'risk_level', 'very_high', 'very_high'),
        ('mean-reversion', 'timeframe', 'tf_1m', '1m'),
        ('mean-reversion', 'timeframe', 'tf_5m', '5m'),
        ('mean-reversion', 'timeframe', 'tf_15m', '15m'),
        ('mean-reversion', 'timeframe', 'tf_30m', '30m'),
        ('mean-reversion', 'lookback_bars', 'bars_20', '20'),
        ('mean-reversion', 'lookback_bars', 'bars_50', '50'),
        ('mean-reversion', 'lookback_bars', 'bars_100', '100'),
        ('mean-reversion', 'zscore_threshold', 'z_15', '1.5'),
        ('mean-reversion', 'zscore_threshold', 'z_20', '2.0'),
        ('mean-reversion', 'zscore_threshold', 'z_25', '2.5'),
        ('mean-reversion', 'entry_mode', 'limit', 'limit'),
        ('mean-reversion', 'entry_mode', 'market', 'market'),

        ('trend-momentum', 'risk_level', 'low', 'low'),
        ('trend-momentum', 'risk_level', 'medium', 'medium'),
        ('trend-momentum', 'risk_level', 'high', 'high'),
        ('trend-momentum', 'risk_level', 'very_high', 'very_high'),
        ('trend-momentum', 'timeframe', 'tf_1m', '1m'),
        ('trend-momentum', 'timeframe', 'tf_5m', '5m'),
        ('trend-momentum', 'timeframe', 'tf_15m', '15m'),
        ('trend-momentum', 'timeframe', 'tf_30m', '30m'),
        ('trend-momentum', 'momentum_window', 'win_10', '10'),
        ('trend-momentum', 'momentum_window', 'win_20', '20'),
        ('trend-momentum', 'momentum_window', 'win_50', '50'),
        ('trend-momentum', 'trend_filter', 'ema_50', 'ema_50'),
        ('trend-momentum', 'trend_filter', 'ema_100', 'ema_100'),
        ('trend-momentum', 'trend_filter', 'ema_200', 'ema_200'),
        ('trend-momentum', 'confirmation_bars', 'bars_1', '1'),
        ('trend-momentum', 'confirmation_bars', 'bars_2', '2'),
        ('trend-momentum', 'confirmation_bars', 'bars_3', '3'),

        ('volatility-breakout', 'risk_level', 'low', 'low'),
        ('volatility-breakout', 'risk_level', 'medium', 'medium'),
        ('volatility-breakout', 'risk_level', 'high', 'high'),
        ('volatility-breakout', 'risk_level', 'very_high', 'very_high'),
        ('volatility-breakout', 'timeframe', 'tf_1m', '1m'),
        ('volatility-breakout', 'timeframe', 'tf_5m', '5m'),
        ('volatility-breakout', 'timeframe', 'tf_15m', '15m'),
        ('volatility-breakout', 'timeframe', 'tf_30m', '30m'),
        ('volatility-breakout', 'breakout_window', 'bars_20', '20'),
        ('volatility-breakout', 'breakout_window', 'bars_55', '55'),
        ('volatility-breakout', 'breakout_window', 'bars_100', '100'),
        ('volatility-breakout', 'stop_atr_multiplier', 'x_15', '1.5'),
        ('volatility-breakout', 'stop_atr_multiplier', 'x_20', '2.0'),
        ('volatility-breakout', 'stop_atr_multiplier', 'x_30', '3.0'),
        ('volatility-breakout', 'volume_filter', 'off', 'off'),
        ('volatility-breakout', 'volume_filter', 'on', 'on')
)
INSERT INTO strategy_attribute_value (id, strategy_attribute_id, slug, value, created_at, updated_at)
SELECT
    gen_random_uuid(),
    sp.id,
    v.slug,
    v.value,
    NOW(),
    NOW()
FROM seed_values v
JOIN strategy s ON s.slug = v.strategy_slug
JOIN strategy_attribute sp
    ON sp.strategy_id = s.id
   AND sp.slug = v.attribute_slug
ON CONFLICT ON CONSTRAINT uq_strategy_param_value_param_slug DO NOTHING;

-- +goose Down
DELETE FROM strategy_attribute_value spv
USING strategy_attribute sp, strategy s
WHERE spv.strategy_attribute_id = sp.id
  AND sp.strategy_id = s.id
  AND (s.slug, sp.slug, spv.slug) IN (
    ('mean-reversion', 'risk_level', 'low'),
    ('mean-reversion', 'risk_level', 'medium'),
    ('mean-reversion', 'risk_level', 'high'),
    ('mean-reversion', 'risk_level', 'very_high'),
    ('mean-reversion', 'timeframe', 'tf_1m'),
    ('mean-reversion', 'timeframe', 'tf_5m'),
    ('mean-reversion', 'timeframe', 'tf_15m'),
    ('mean-reversion', 'timeframe', 'tf_30m'),
    ('mean-reversion', 'lookback_bars', 'bars_20'),
    ('mean-reversion', 'lookback_bars', 'bars_50'),
    ('mean-reversion', 'lookback_bars', 'bars_100'),
    ('mean-reversion', 'zscore_threshold', 'z_15'),
    ('mean-reversion', 'zscore_threshold', 'z_20'),
    ('mean-reversion', 'zscore_threshold', 'z_25'),
    ('mean-reversion', 'entry_mode', 'limit'),
    ('mean-reversion', 'entry_mode', 'market'),

    ('trend-momentum', 'risk_level', 'low'),
    ('trend-momentum', 'risk_level', 'medium'),
    ('trend-momentum', 'risk_level', 'high'),
    ('trend-momentum', 'risk_level', 'very_high'),
    ('trend-momentum', 'timeframe', 'tf_1m'),
    ('trend-momentum', 'timeframe', 'tf_5m'),
    ('trend-momentum', 'timeframe', 'tf_15m'),
    ('trend-momentum', 'timeframe', 'tf_30m'),
    ('trend-momentum', 'momentum_window', 'win_10'),
    ('trend-momentum', 'momentum_window', 'win_20'),
    ('trend-momentum', 'momentum_window', 'win_50'),
    ('trend-momentum', 'trend_filter', 'ema_50'),
    ('trend-momentum', 'trend_filter', 'ema_100'),
    ('trend-momentum', 'trend_filter', 'ema_200'),
    ('trend-momentum', 'confirmation_bars', 'bars_1'),
    ('trend-momentum', 'confirmation_bars', 'bars_2'),
    ('trend-momentum', 'confirmation_bars', 'bars_3'),

    ('volatility-breakout', 'risk_level', 'low'),
    ('volatility-breakout', 'risk_level', 'medium'),
    ('volatility-breakout', 'risk_level', 'high'),
    ('volatility-breakout', 'risk_level', 'very_high'),
    ('volatility-breakout', 'timeframe', 'tf_1m'),
    ('volatility-breakout', 'timeframe', 'tf_5m'),
    ('volatility-breakout', 'timeframe', 'tf_15m'),
    ('volatility-breakout', 'timeframe', 'tf_30m'),
    ('volatility-breakout', 'breakout_window', 'bars_20'),
    ('volatility-breakout', 'breakout_window', 'bars_55'),
    ('volatility-breakout', 'breakout_window', 'bars_100'),
    ('volatility-breakout', 'stop_atr_multiplier', 'x_15'),
    ('volatility-breakout', 'stop_atr_multiplier', 'x_20'),
    ('volatility-breakout', 'stop_atr_multiplier', 'x_30'),
    ('volatility-breakout', 'volume_filter', 'off'),
    ('volatility-breakout', 'volume_filter', 'on')
  );
