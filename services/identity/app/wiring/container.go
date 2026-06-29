package wiring

import (
	"database/sql"

	appconfig "github.com/stratflow-labs/stratflow/internal/app/config"

	"github.com/samber/do"
)

func BuildContainer(cfg *appconfig.Config, sqlDB *sql.DB) *do.Injector {
	injector := do.New()
	if cfg != nil {
		do.ProvideValue(injector, *cfg)
	}
	do.ProvideValue(injector, sqlDB)
	do.ProvideNamedValue(injector, "db", sqlDB)
	return injector
}
