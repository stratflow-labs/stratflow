package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/stratflow-labs/stratflow/internal/app/config"
	"github.com/stratflow-labs/stratflow/internal/foundation/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// ConnectPostgres opens a postgres connection using database/sql with retry and applies pool settings.
func ConnectPostgres(cfg config.DBConfig) (*sql.DB, func() error, error) {
	db, err := openWithRetry(cfg)
	if err != nil {
		return nil, nil, err
	}

	configurePool(db, cfg)

	return db, func() error { return db.Close() }, nil
}

func openWithRetry(cfg config.DBConfig) (*sql.DB, error) {
	retries := cfg.ConnectRetries
	if retries <= 0 {
		retries = 1
	}
	delay := cfg.ConnectRetryInterval
	if delay <= 0 {
		delay = 2 * time.Second
	}

	var db *sql.DB
	var err error
	for attempt := 1; attempt <= retries; attempt++ {
		db, err = sql.Open("pgx", cfg.DSN)
		if err == nil {
			if pingErr := db.PingContext(context.Background()); pingErr == nil {
				return db, nil
			}
			err = db.PingContext(context.Background())
		}

		logger.Warn("database connection attempt failed",
			"attempt", attempt,
			"maxAttempts", retries,
			"error", err,
		)

		if attempt < retries {
			timer := time.NewTimer(delay)
			<-timer.C
		}
	}
	return nil, fmt.Errorf("open database connection: %w", err)
}

// configurePool applies connection pool settings.
// MaxOpenConns limits load; MaxIdleConns keeps hot connections; lifetimes handle LB/firewall timeouts.
func configurePool(db *sql.DB, cfg config.DBConfig) {
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	logger.Info("database connection pool configured",
		"maxOpenConns", cfg.MaxOpenConns,
		"maxIdleConns", cfg.MaxIdleConns,
		"connMaxLifetime", cfg.ConnMaxLifetime,
		"connMaxIdleTime", cfg.ConnMaxIdleTime,
	)
}
