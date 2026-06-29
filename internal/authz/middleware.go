package authz

import (
	"errors"
	"net/http"

	"github.com/stratflow-labs/stratflow/internal/authkit"
	"github.com/stratflow-labs/stratflow/internal/foundation/httpx/respond"
)

// Middleware enforces authorization using provided engine.
func Middleware(engine *Engine) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if engine == nil {
				next.ServeHTTP(w, r)
				return
			}

			// Public endpoint: no auth context required.
			if !engine.RequireAuth(r) {
				next.ServeHTTP(w, r)
				return
			}

			role, okRole := authkit.RoleFromContext(r.Context())
			userID, okUser := authkit.UserIDFromContext(r.Context())
			if !okRole || !okUser {
				respond.Problem(w, http.StatusUnauthorized, "session.unauthorized", "missing auth context", nil)
				return
			}

			sub := Subject{Role: role, UserID: userID}
			if err := engine.Authorize(r, sub); err != nil {
				if errors.Is(err, ErrForbidden) {
					respond.Problem(w, http.StatusForbidden, "session.forbidden", "access denied", nil)
					return
				}
				respond.Problem(w, http.StatusInternalServerError, "session.internal", err.Error(), nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
