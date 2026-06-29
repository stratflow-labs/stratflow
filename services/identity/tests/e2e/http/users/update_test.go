package tests

import (
	"net/http"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	testsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/http/testkit"
)

func TestUsersUpdate(t *testing.T) {
	testsupport.Run(t, "partial user update preserves untouched fields over http", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		created := testsupport.CreateIdentityE2EUser(tt, admin, "user")

		var originalLastName string
		var originalEmail string

		cfg.Step("load original user state", func() {
			resp := admin.Get(tt, testsupport.UserPath(created.ID), nil)
			testsupport.RequireStatus(tt, resp, http.StatusOK)
			data := testsupport.RequireDataMap(tt, resp)
			originalLastName = testsupport.RequireStringField(tt, data, "lastName")
			originalEmail = testsupport.RequireStringField(tt, data, "email")
		})

		cfg.Step("update only the display name", func() {
			resp := admin.Put(tt, testsupport.UserPath(created.ID), map[string]any{
				"name": "Patched Name",
			})
			testsupport.RequireStatus(tt, resp, http.StatusOK)
			data := testsupport.RequireDataMap(tt, resp)
			testsupport.RequireEqualString(tt, testsupport.RequireStringField(tt, data, "name"), "Patched Name", "name")
			testsupport.RequireEqualString(tt, testsupport.RequireStringField(tt, data, "lastName"), originalLastName, "lastName")
			testsupport.RequireEqualString(tt, testsupport.RequireStringField(tt, data, "email"), originalEmail, "email")
		})
	})
}
