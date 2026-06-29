package identitydbsqlc

import (
	"context"
	"database/sql"

	"github.com/stratflow-labs/stratflow/internal/foundation/tx/sqltx"
)

// NewQuerier returns sqlc queries bound to tx from context if present.
func NewQuerier(ctx context.Context, db *sql.DB) *Queries {
	base := New(db)
	if tx := sqltx.FromCtx(ctx); tx != nil {
		return base.WithTx(tx)
	}
	return base
}
