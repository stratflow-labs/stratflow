package db

import (
	"database/sql"

	"github.com/stratflow-labs/stratflow/internal/app/config"
	foundationdb "github.com/stratflow-labs/stratflow/internal/foundation/db"
)

// NewDB opens a postgres connection using shared foundation helper (retry + pool).
func NewDB(cfg *config.Config) (*sql.DB, func() error, error) {
	return foundationdb.ConnectPostgres(cfg.DB)
}
