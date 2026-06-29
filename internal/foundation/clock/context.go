package clock

import (
	"context"
	"time"
)

type ctxKey struct{}

// IntoCtx stores a Clock in the context.
func IntoCtx(ctx context.Context, c Clock) context.Context {
	return context.WithValue(ctx, ctxKey{}, c)
}

// FromCtx retrieves a Clock from the context.
// Returns nil if not present.
//
//nolint:ireturn // Clock is the package port; callers must stay decoupled from implementations.
func FromCtx(ctx context.Context) Clock {
	if c, ok := ctx.Value(ctxKey{}).(Clock); ok {
		return c
	}
	return nil
}

// NowFromCtx returns the current time from the Clock in context,
// or time.Now().UTC() if no Clock is present.
func NowFromCtx(ctx context.Context) time.Time {
	if c := FromCtx(ctx); c != nil {
		return c.Now()
	}
	return time.Now().UTC()
}
