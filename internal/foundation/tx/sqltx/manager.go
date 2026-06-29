package sqltx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Manager struct {
	DB *sql.DB
}

type txKey struct{}

func (m *Manager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if m == nil || m.DB == nil {
		return errors.New("sqltx: db is nil")
	}
	if fn == nil {
		return errors.New("sqltx: fn is nil")
	}

	if existing := FromCtx(ctx); existing != nil {
		return fn(ctx)
	}

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sqltx: begin tx: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)
	if err := fn(txCtx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			return fmt.Errorf("sqltx: rollback after error: %w: %w", rollbackErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("sqltx: commit tx: %w", err)
	}

	return nil
}

func FromCtx(ctx context.Context) *sql.Tx {
	if ctx == nil {
		return nil
	}

	tx, _ := ctx.Value(txKey{}).(*sql.Tx)
	return tx
}
