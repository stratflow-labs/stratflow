package integration_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stratflow-labs/stratflow/internal/authkit"
	"github.com/stratflow-labs/stratflow/internal/authz"
	"github.com/stratflow-labs/stratflow/internal/httpserver"
	testsupport "github.com/stratflow-labs/stratflow/services/identity/tests/integration/testkit"
)

func TestSessionAuthz(t *testing.T) {
	engine := testsupport.LoadEngineFromRepo(t)

	reqLogout, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/auth/logout", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if err := engine.Authorize(reqLogout, authz.Subject{Role: "user", UserID: "user-1"}); err != nil {
		t.Fatalf("authorize user logout should allow, got %v", err)
	}

	router := httpserver.New()
	router.Use(authkit.AuthenticateWithSkipper(testsupport.FakeVerifier{}, authkit.PathPrefixSkipper("/api/auth/login")))
	router.Use(authz.Middleware(engine))

	router.Mount("/api/auth/logout", testsupport.StatusHandler(http.StatusNoContent))
	router.Mount("/api/auth/refresh", testsupport.StatusHandler(http.StatusOK))

	token := func(v string) string { return "Bearer " + v }

	testsupport.RunCases(t, router, []testsupport.Case{
		{Name: "user logout -> 204", Method: http.MethodPost, Path: "/api/auth/logout", Token: token("user-token"), WantStatus: http.StatusNoContent},
		{Name: "no token logout -> 401", Method: http.MethodPost, Path: "/api/auth/logout", WantStatus: http.StatusUnauthorized},
		{Name: "expired token -> 401", Method: http.MethodPost, Path: "/api/auth/logout", Token: token("expired-token"), WantStatus: http.StatusUnauthorized},
		{Name: "refresh valid -> 200", Method: http.MethodPost, Path: "/api/auth/refresh", Token: token("user-token"), WantStatus: http.StatusOK},
	})
}
