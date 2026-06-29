// pkg/authkit/middleware.go
package authkit

import (
	"context"
	"net/http"
	"strings"

	"github.com/stratflow-labs/stratflow/internal/foundation/httpx/respond"
)

type Claims struct {
	UserID string
	Role   string
}

// Abstract interface - each service can implement it as needed.
type AccessTokenVerifier interface {
	Verify(ctx context.Context, token string) (Claims, error)
}

type ctxKey string

const ctxUserIDKey ctxKey = "user_id"
const ctxRoleKey ctxKey = "user_role"

type Skipper func(r *http.Request) bool

func Authenticate(verifier AccessTokenVerifier) func(http.Handler) http.Handler {
	return AuthenticateWithSkipper(verifier, nil)
}

func AuthenticateWithSkipper(verifier AccessTokenVerifier, skip Skipper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if skip != nil && skip(r) {
				next.ServeHTTP(w, r)
				return
			}

			header := strings.TrimSpace(r.Header.Get("Authorization"))
			if header == "" {
				respond.Problem(w, http.StatusUnauthorized, "session.unauthorized", "missing Authorization header", nil)
				return
			}

			const bearer = "Bearer "
			if !strings.HasPrefix(header, bearer) {
				respond.Problem(w, http.StatusUnauthorized, "session.unauthorized", "invalid authorization schema", nil)
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(header, bearer))
			if token == "" {
				respond.Problem(w, http.StatusUnauthorized, "session.unauthorized", "empty token", nil)
				return
			}

			claims, err := verifier.Verify(r.Context(), token)
			if err != nil {
				respond.Problem(w, http.StatusUnauthorized, "session.unauthorized", "token verification failed", nil)
				return
			}

			ctx := WithClaims(r.Context(), claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithClaims(ctx context.Context, claims Claims) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = context.WithValue(ctx, ctxUserIDKey, claims.UserID)
	ctx = context.WithValue(ctx, ctxRoleKey, claims.Role)
	return ctx
}

// Helpers for reading values from context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxUserIDKey).(string)
	return v, ok
}

func RoleFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxRoleKey).(string)
	return v, ok
}

// PathPrefixSkipper returns a skipper that bypasses auth for any request whose path
// has one of the provided prefixes.
func PathPrefixSkipper(prefixes ...string) Skipper {
	return func(r *http.Request) bool {
		if len(prefixes) == 0 || r == nil || r.URL == nil {
			return false
		}
		path := r.URL.Path
		for _, p := range prefixes {
			if strings.HasPrefix(path, p) {
				return true
			}
		}
		return false
	}
}
