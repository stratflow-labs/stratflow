package ports

import "context"

// TxManager describes the infrastructure mechanism for working with transactions.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
