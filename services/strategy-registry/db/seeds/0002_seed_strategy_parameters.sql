-- +goose Up
-- Seed strategy attributes for local/dev.

WITH seed_params (strategy_slug, slug, name, description) AS (
    VALUES
        ('mean-reversion', 'risk_level', 'Risk Level', 'Risk profile level used by the strategy model'),
        ('mean-reversion', 'timeframe', 'Timeframe', 'Working timeframe used for feature extraction and execution'),
        ('mean-reversion', 'lookback_bars', 'Lookback Bars', 'Lookback window size in candles'),
        ('mean-reversion', 'zscore_threshold', 'Z-Score Threshold', 'Z-score threshold for signal activation'),
        ('mean-reversion', 'entry_mode', 'Entry Mode', 'Entry type selection for execution style'),

        ('trend-momentum', 'risk_level', 'Risk Level', 'Risk profile level used by the strategy model'),
        ('trend-momentum', 'timeframe', 'Timeframe', 'Working timeframe used for momentum evaluation'),
        ('trend-momentum', 'momentum_window', 'Momentum Window', 'Window size for momentum calculation'),
        ('trend-momentum', 'trend_filter', 'Trend Filter', 'Trend filter mode applied before entry'),
        ('trend-momentum', 'confirmation_bars', 'Confirmation Bars', 'Number of bars required to confirm trend signal'),

        ('volatility-breakout', 'risk_level', 'Risk Level', 'Risk profile level used by the strategy model'),
        ('volatility-breakout', 'timeframe', 'Timeframe', 'Working timeframe used for breakout detection'),
        ('volatility-breakout', 'breakout_window', 'Breakout Window', 'Window size in candles for breakout confirmation'),
        ('volatility-breakout', 'stop_atr_multiplier', 'Stop ATR Multiplier', 'ATR multiplier for stop placement'),
        ('volatility-breakout', 'volume_filter', 'Volume Filter', 'Volume filter mode for breakout validation')
)
INSERT INTO strategy_attribute (id, strategy_id, slug, name, description, created_at, updated_at)
SELECT
    gen_random_uuid(),
    s.id,
    p.slug,
    p.name,
    p.description,
    NOW(),
    NOW()
FROM seed_params p
JOIN strategy s ON s.slug = p.strategy_slug
ON CONFLICT ON CONSTRAINT uq_strategy_param_strategy_slug DO NOTHING;

-- +goose Down
DELETE FROM strategy_attribute sp
USING strategy s
WHERE sp.strategy_id = s.id
  AND (s.slug, sp.slug) IN (
    ('mean-reversion', 'risk_level'),
    ('mean-reversion', 'timeframe'),
    ('mean-reversion', 'lookback_bars'),
    ('mean-reversion', 'zscore_threshold'),
    ('mean-reversion', 'entry_mode'),
    ('trend-momentum', 'risk_level'),
    ('trend-momentum', 'timeframe'),
    ('trend-momentum', 'momentum_window'),
    ('trend-momentum', 'trend_filter'),
    ('trend-momentum', 'confirmation_bars'),
    ('volatility-breakout', 'risk_level'),
    ('volatility-breakout', 'timeframe'),
    ('volatility-breakout', 'breakout_window'),
    ('volatility-breakout', 'stop_atr_multiplier'),
    ('volatility-breakout', 'volume_filter')
  );
