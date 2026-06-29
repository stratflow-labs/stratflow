-- +goose Up
-- Seed strategies for local/dev.

INSERT INTO strategy (id, slug, name, description, created_at, updated_at)
VALUES
    (
        gen_random_uuid(),
        'mean-reversion',
        'Mean Reversion',
        'Mean reversion strategy template for models with configurable attributes.',
        NOW(),
        NOW()
    ),
    (
        gen_random_uuid(),
        'trend-momentum',
        'Trend Momentum',
        'Trend-following strategy with momentum filters and dynamic risk profile.',
        NOW(),
        NOW()
    ),
    (
        gen_random_uuid(),
        'volatility-breakout',
        'Volatility Breakout',
        'Breakout strategy focused on volatility expansion, stop control and volume filters.',
        NOW(),
        NOW()
    )
ON CONFLICT (slug) DO NOTHING;

-- +goose Down
DELETE FROM strategy
WHERE slug IN ('mean-reversion', 'trend-momentum', 'volatility-breakout');
