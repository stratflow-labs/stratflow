package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	appbootstrap "github.com/stratflow-labs/stratflow/internal/app/bootstrap"
	appconfig "github.com/stratflow-labs/stratflow/internal/app/config"
	"github.com/stratflow-labs/stratflow/internal/foundation/logger"
)

type LoadConfigFunc func() (appconfig.Config, error)
type OpenDBFunc func(cfg *appconfig.Config) (*sql.DB, func() error, error)
type WithDBActionFunc func(context.Context, *appconfig.Config, *sql.DB) error
type BuildHTTPHandlerFunc func(context.Context, *appconfig.Config, *sql.DB) (http.Handler, error)
type MigrateActionFunc func(*sql.DB) error
type SeedActionFunc func(context.Context, *sql.DB) error
type SeedNoContextActionFunc func(*sql.DB) error

// WithConfigAndDB executes action with loaded config and opened database.
func WithConfigAndDB(
	ctx context.Context,
	loadConfig LoadConfigFunc,
	openDB OpenDBFunc,
	action WithDBActionFunc,
) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	logger.ReloadFromEnv()

	dbConn, closeDB, err := openDB(&cfg)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() {
		if closeDB == nil {
			return
		}
		if cerr := closeDB(); cerr != nil {
			logger.Err("close database connection", cerr)
		}
	}()

	if action == nil {
		return nil
	}
	return action(ctx, &cfg, dbConn)
}

// ServeHTTP starts HTTP server with defaults from config.
func ServeHTTP(ctx context.Context, cfg *appconfig.Config, handler http.Handler) error {
	server := &http.Server{
		Addr:              cfg.HTTP.ListenAddr,
		Handler:           handler,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	}
	return appbootstrap.GracefulServe(ctx, server, cfg)
}

// Serve runs service-specific HTTP server builder with shared config/database lifecycle.
func Serve(ctx context.Context, loadConfig LoadConfigFunc, openDB OpenDBFunc, build BuildHTTPHandlerFunc) error {
	return WithConfigAndDB(ctx, loadConfig, openDB, func(ctx context.Context, cfg *appconfig.Config, dbConn *sql.DB) error {
		if build == nil {
			return errors.New("http handler builder is nil")
		}
		handler, err := build(ctx, cfg, dbConn)
		if err != nil {
			return fmt.Errorf("build http handler: %w", err)
		}
		if err := ServeHTTP(ctx, cfg, handler); err != nil {
			return fmt.Errorf("serve http: %w", err)
		}
		return nil
	})
}

// Migrate runs service-specific migrations with shared config/database lifecycle.
func Migrate(ctx context.Context, loadConfig LoadConfigFunc, openDB OpenDBFunc, run MigrateActionFunc) error {
	return WithConfigAndDB(ctx, loadConfig, openDB, func(_ context.Context, _ *appconfig.Config, dbConn *sql.DB) error {
		if run == nil {
			return nil
		}
		if err := run(dbConn); err != nil {
			return fmt.Errorf("run migrations: %w", err)
		}
		return nil
	})
}

// Seed runs service-specific seeders with shared config/database lifecycle.
func Seed(ctx context.Context, loadConfig LoadConfigFunc, openDB OpenDBFunc, run SeedActionFunc) error {
	return WithConfigAndDB(ctx, loadConfig, openDB, func(ctx context.Context, _ *appconfig.Config, dbConn *sql.DB) error {
		if run == nil {
			return nil
		}
		if err := run(ctx, dbConn); err != nil {
			return fmt.Errorf("run seeders: %w", err)
		}
		return nil
	})
}

// SeedNoContext adapts services whose seed runner does not require context.
func SeedNoContext(ctx context.Context, loadConfig LoadConfigFunc, openDB OpenDBFunc, run SeedNoContextActionFunc) error {
	return Seed(ctx, loadConfig, openDB, func(_ context.Context, dbConn *sql.DB) error {
		if run == nil {
			return nil
		}
		return run(dbConn)
	})
}
