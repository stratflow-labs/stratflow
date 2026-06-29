package db

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/stratflow-labs/stratflow/internal/foundation/tx/sqltx"
)

const (
	pgUniqueViolationCode = "23505"
	maxCloneBatchSize     = 100000
)

func intToInt32(v int) (int32, error) {
	if v > math.MaxInt32 || v < math.MinInt32 {
		return 0, fmt.Errorf("value %d is out of int32 range", v)
	}
	return int32(v), nil
}

func timeNowUTC() time.Time {
	return time.Now().UTC()
}

func nonZeroTime(times ...time.Time) time.Time {
	for _, t := range times {
		if !t.IsZero() {
			return t
		}
	}
	return time.Now().UTC()
}

// execContext executes a query that doesn't return rows.
// Uses transaction from context if present.
func execContext(ctx context.Context, db *sql.DB, query string, args ...any) (sql.Result, error) {
	if tx := sqltx.FromCtx(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return db.ExecContext(ctx, query, args...)
}

// queryContext executes a query that returns rows.
// Uses transaction from context if present.
func queryContext(ctx context.Context, db *sql.DB, query string, args ...any) (*sql.Rows, error) {
	if tx := sqltx.FromCtx(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}
	return db.QueryContext(ctx, query, args...)
}

// queryRowContext executes a query that returns at most one row.
// Uses transaction from context if present.
func queryRowContext(ctx context.Context, db *sql.DB, query string, args ...any) *sql.Row {
	if tx := sqltx.FromCtx(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return db.QueryRowContext(ctx, query, args...)
}
