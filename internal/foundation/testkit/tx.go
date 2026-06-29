// internal/foundation/testkit/tx.go
package testkit

import "context"

type NoTx struct{}

func (NoTx) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type TxRecorder struct {
	Calls int
	Err   error
}

func (t *TxRecorder) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	t.Calls++
	if t.Err != nil {
		return t.Err
	}
	return fn(ctx)
}
