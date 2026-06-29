package bootstrap

import (
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/app/config"
)

func LoadConfig() (config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return config.Config{}, fmt.Errorf("load config: %w", err)
	}
	return cfg, nil
}
