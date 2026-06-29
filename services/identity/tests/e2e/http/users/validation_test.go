package tests

import (
	"net/http"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	testsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/http/testkit"
)

func TestUsersValidation(t *testing.T) {
	testsupport.Run(t, "users create rejects invalid input", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)

		cfg.Step("create user rejects invalid email and password payload", func() {
			// Validation regression: an invalid email must not create a user.
			resp := admin.Post(tt, "/users", map[string]any{
				"login":    "invalid-user",
				"email":    "not-an-email",
				"password": "short",
				"name":     "Invalid",
				"lastName": "User",
				"gender":   1,
				"role":     "user",
			})
			testsupport.RequireStatus(tt, resp, http.StatusBadRequest)
		})
	})
}
