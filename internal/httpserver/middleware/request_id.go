package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey struct{}

// RequestID adds a request ID to the context and response header.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewString()
		ctx := context.WithValue(r.Context(), ctxKey{}, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFromContext returns request id if it was set by RequestID middleware.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	id, ok := ctx.Value(ctxKey{}).(string)
	return id, ok
}
