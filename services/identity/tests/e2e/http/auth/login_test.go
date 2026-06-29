package tests

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	httptestkit "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/http/testkit"
)

func TestAuthLogin(t *testing.T) {
	httptestkit.Run(t, "admin can login through public HTTP API", func(cfg *axiom.Config, env httptestkit.Env) {
		tt := httptestkit.CaseT(t, cfg)
		var client httptestkit.Client

		cfg.Step("login as seeded admin", func() {
			client = httptestkit.LoginAdmin(tt, env)
		})

		cfg.Step("read current user with issued token", func() {
			resp := client.Get(tt, "/users/me", nil)
			httptestkit.RequireStatus(tt, resp, 200)
			httptestkit.RequireEqualString(tt, httptestkit.RequireStringField(tt, httptestkit.RequireDataMap(tt, resp), "email"), env.AdminEmail, "email")
		})
	})
}

func TestAuthLoginInvalidCredentials(t *testing.T) {
	httptestkit.Run(t, "invalid login credentials are rejected", func(cfg *axiom.Config, env httptestkit.Env) {
		tt := httptestkit.CaseT(t, cfg)
		cfg.Step("try login with wrong password", func() {
			// The negative scenario locks in the external contract: invalid credentials must not issue an access token.
			resp := httptestkit.NewClient(env, "").Post(tt, "/auth/login", map[string]any{
				"login":    env.AdminEmail,
				"password": "definitely-wrong-password",
			})
			httptestkit.RequireErrorCode(tt, resp, 401)
		})
	})
}
