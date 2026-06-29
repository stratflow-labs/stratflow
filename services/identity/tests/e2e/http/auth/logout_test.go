package tests

import (
	"net/http"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	httptestkit "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/http/testkit"
)

func TestAuthLogout(t *testing.T) {
	httptestkit.Run(t, "logout invalidates the issued HTTP token", func(cfg *axiom.Config, env httptestkit.Env) {
		tt := httptestkit.CaseT(t, cfg)
		client := httptestkit.LoginAdmin(tt, env)

		cfg.Step("logout with the issued token", func() {
			httptestkit.RequireStatus(tt, client.Post(tt, "/auth/logout", nil), http.StatusNoContent)
		})

		cfg.Step("reusing the same token is rejected", func() {
			httptestkit.RequireErrorCode(tt, client.Get(tt, "/users/me", nil), http.StatusUnauthorized)
		})
	})
}
