package bootstrap

import (
	"errors"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/app/config"
)

func LoadConfig() (config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return config.Config{}, fmt.Errorf("load config: %w", err)
	}
	if err := validateConfig(&cfg); err != nil {
		return config.Config{}, fmt.Errorf("invalid config: %w", err)
	}
	return cfg, nil
}

func validateConfig(cfg *config.Config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	if cfg.Security.TokenHashSecret == "" {
		return errors.New("TOKEN_HASH_SECRET is required (generate with: openssl rand -hex 32)")
	}
	if len(cfg.Security.TokenHashSecret) < 32 {
		return fmt.Errorf("TOKEN_HASH_SECRET must be at least 32 characters (got %d)", len(cfg.Security.TokenHashSecret))
	}
	if cfg.AppEnv != "localhost" && cfg.Security.TokenHashSecret == "dev-secret-change-me-in-production-use-openssl-rand-hex-32" {
		return errors.New("TOKEN_HASH_SECRET must be changed in non-localhost environments")
	}
	return nil
}
