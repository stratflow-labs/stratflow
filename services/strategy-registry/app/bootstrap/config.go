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
	return nil
}
