package tests

import (
	"net/http"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	testsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/http/testkit"
)

func TestUsersUniqueness(t *testing.T) {
	testsupport.Run(t, "create user enforces unique login and email over http", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		created := testsupport.CreateIdentityE2EUser(tt, admin, "user")

		cfg.Step("duplicate login is rejected", func() {
			resp := admin.Post(tt, "/users", map[string]any{
				"login":    created.Login,
				"email":    testsupport.UniqueE2ELogin("dup") + "@example.test",
				"password": "E2e-password-12345",
				"name":     "Duplicate",
				"lastName": "Login",
				"gender":   1,
				"role":     "user",
			})
			testsupport.RequireStatus(tt, resp, http.StatusConflict)
			testsupport.RequireErrorCode(tt, resp, http.StatusConflict)
		})

		cfg.Step("duplicate email is rejected", func() {
			resp := admin.Post(tt, "/users", map[string]any{
				"login":    testsupport.UniqueE2ELogin("dup"),
				"email":    created.Email,
				"password": "E2e-password-12345",
				"name":     "Duplicate",
				"lastName": "Email",
				"gender":   1,
				"role":     "user",
			})
			testsupport.RequireStatus(tt, resp, http.StatusConflict)
			testsupport.RequireErrorCode(tt, resp, http.StatusConflict)
		})
	})
}
