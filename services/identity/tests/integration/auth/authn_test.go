package integration_test

import (
	"net/http"
	"testing"

	"github.com/stratflow-labs/stratflow/internal/authkit"
	"github.com/stratflow-labs/stratflow/internal/httpserver"
	testsupport "github.com/stratflow-labs/stratflow/services/identity/tests/integration/testkit"
)

func TestAuthn(t *testing.T) {
	router := httpserver.New()
	router.Use(authkit.AuthenticateWithSkipper(testsupport.FakeVerifier{}, authkit.PathPrefixSkipper("/api/auth/login")))

	// Public login handler
	router.Mount("/api/auth/login", testsupport.StatusHandler(http.StatusOK))
	// Protected handler
	router.Mount("/api/users", testsupport.StatusHandler(http.StatusOK))
	router.Mount("/api/auth/logout", testsupport.StatusHandler(http.StatusNoContent))

	testsupport.RunCases(t, router, []testsupport.Case{
		{Name: "login public no token", Method: http.MethodPost, Path: "/api/auth/login", WantStatus: http.StatusOK},
		{Name: "users no token -> 401", Method: http.MethodGet, Path: "/api/users", WantStatus: http.StatusUnauthorized},
		{Name: "users invalid token -> 401", Method: http.MethodGet, Path: "/api/users", Token: "Bearer bad-token", WantStatus: http.StatusUnauthorized},
		{Name: "users valid token -> 200", Method: http.MethodGet, Path: "/api/users", Token: "Bearer user-token", WantStatus: http.StatusOK},
	})
}
