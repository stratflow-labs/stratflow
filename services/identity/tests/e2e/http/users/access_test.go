package tests

import (
	"net/http"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	testsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/http/testkit"
)

func TestUsersAccess(t *testing.T) {
	testsupport.Run(t, "users endpoint role matrix", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		createdUser := testsupport.CreateIdentityE2EUser(tt, admin, "user")
		createdManager := testsupport.CreateIdentityE2EUser(tt, admin, "manager")
		user := testsupport.LoginCreatedUser(tt, env, &createdUser)
		manager := testsupport.LoginCreatedUser(tt, env, &createdManager)

		anonymous := testsupport.NewClient(env, "")
		invalid := testsupport.NewClient(env, "invalid-token")
		userClient := testsupport.NewClient(env, user.Token)
		managerClient := testsupport.NewClient(env, manager.Token)

		cfg.Step("protected list rejects anonymous and invalid tokens", func() {
			// The protected user list must not reveal data without valid authorization.
			testsupport.RequireErrorCode(tt, anonymous.Get(tt, "/users", nil), http.StatusUnauthorized)
			testsupport.RequireErrorCode(tt, invalid.Get(tt, "/users", nil), http.StatusUnauthorized)
		})

		cfg.Step("list users follows role permissions", func() {
			// user has no list permission; manager/admin can read the list according to policy.csv.
			testsupport.RequireErrorCode(tt, userClient.Get(tt, "/users", nil), http.StatusForbidden)
			testsupport.RequireStatus(tt, managerClient.Get(tt, "/users", nil), http.StatusOK)
			testsupport.RequireStatus(tt, admin.Get(tt, "/users", nil), http.StatusOK)
		})

		cfg.Step("admin-only create is forbidden for user and manager", func() {
			payload := map[string]any{
				"login":    "forbidden-create",
				"email":    "forbidden.create@example.test",
				"password": "E2e-password-12345",
				"name":     "Forbidden",
				"lastName": "Create",
				"gender":   1,
				"role":     "user",
			}
			testsupport.RequireErrorCode(tt, userClient.Post(tt, "/users", payload), http.StatusForbidden)
			testsupport.RequireErrorCode(tt, managerClient.Post(tt, "/users", payload), http.StatusForbidden)
		})

		cfg.Step("foreign user access is denied for user and allowed for manager/admin", func() {
			// IDOR regression: a regular user must not read another user's userId through the path parameter.
			testsupport.RequireStatus(tt, userClient.Get(tt, testsupport.UserPath(user.ID), nil), http.StatusOK)
			testsupport.RequireErrorCode(tt, userClient.Get(tt, testsupport.UserPath(manager.ID), nil), http.StatusForbidden)
			testsupport.RequireStatus(tt, managerClient.Get(tt, testsupport.UserPath(user.ID), nil), http.StatusOK)
			testsupport.RequireStatus(tt, admin.Get(tt, testsupport.UserPath(user.ID), nil), http.StatusOK)
		})

		cfg.Step("foreign update and delete are restricted", func() {
			// Verify broken access control on modifying and deleting someone else's data.
			testsupport.RequireErrorCode(tt, userClient.Put(tt, testsupport.UserPath(manager.ID), map[string]any{"name": "Stolen"}), http.StatusForbidden)
			testsupport.RequireErrorCode(tt, managerClient.Delete(tt, testsupport.UserPath(user.ID)), http.StatusForbidden)
		})
	})
}
