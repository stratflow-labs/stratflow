package integration_test

import (
	"net/http"
	"testing"

	"github.com/stratflow-labs/stratflow/internal/authkit"
	"github.com/stratflow-labs/stratflow/internal/authz"
	"github.com/stratflow-labs/stratflow/internal/httpserver"
	testsupport "github.com/stratflow-labs/stratflow/services/identity/tests/integration/testkit"
)

// User module integration: verify role-based access and own-resource rules.
func TestUsersAuthz(t *testing.T) {
	engine := testsupport.LoadEngineFromRepo(t)

	router := httpserver.New()
	router.Use(authkit.AuthenticateWithSkipper(testsupport.FakeVerifier{}, authkit.PathPrefixSkipper("/api/auth/login")))
	router.Use(authz.Middleware(engine))

	// Stub handlers (authn/authz decide access).
	router.Mount("/api/users", testsupport.StatusHandler(http.StatusOK))
	router.Mount("/api/users/user-1", testsupport.StatusHandler(http.StatusOK))
	router.Mount("/api/users/other", testsupport.StatusHandler(http.StatusOK))
	router.Mount("/api/users/any", testsupport.StatusHandler(http.StatusNoContent))

	token := func(v string) string { return "Bearer " + v }

	testsupport.RunCases(t, router, []testsupport.Case{
		// list
		{Name: "user list users -> 403", Method: http.MethodGet, Path: "/api/users", Token: token("user-token"), WantStatus: http.StatusForbidden},
		{Name: "manager list users -> 200", Method: http.MethodGet, Path: "/api/users", Token: token("manager-token"), WantStatus: http.StatusOK},
		{Name: "admin list users -> 200", Method: http.MethodGet, Path: "/api/users", Token: token("admin-token"), WantStatus: http.StatusOK},
		// get by id
		{Name: "user get self -> 200", Method: http.MethodGet, Path: "/api/users/user-1", Token: token("user-token"), WantStatus: http.StatusOK},
		{Name: "user get other -> 403", Method: http.MethodGet, Path: "/api/users/other", Token: token("user-token"), WantStatus: http.StatusForbidden},
		{Name: "manager get other -> 200", Method: http.MethodGet, Path: "/api/users/other", Token: token("manager-token"), WantStatus: http.StatusOK},
		{Name: "admin get other -> 200", Method: http.MethodGet, Path: "/api/users/other", Token: token("admin-token"), WantStatus: http.StatusOK},
		// delete
		{Name: "user delete user -> 403", Method: http.MethodDelete, Path: "/api/users/any", Token: token("user-token"), WantStatus: http.StatusForbidden},
		{Name: "manager delete user -> 403", Method: http.MethodDelete, Path: "/api/users/any", Token: token("manager-token"), WantStatus: http.StatusForbidden},
		{Name: "admin delete user -> 204", Method: http.MethodDelete, Path: "/api/users/any", Token: token("admin-token"), WantStatus: http.StatusNoContent},
	})
}
