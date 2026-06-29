package tests

import (
	"net/http"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	testsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/http/testkit"
)

func TestUsersPayload(t *testing.T) {
	testsupport.Run(t, "admin update preserves text payload as data", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		user := testsupport.CreateIdentityE2EUser(tt, admin, "user")

		cfg.Step("update user with xss-like display name", func() {
			// The API accepts user-provided text; the e2e test locks in that the payload remains public-contract data.
			const displayName = `<script>alert("e2e")</script>`
			resp := admin.Put(tt, testsupport.UserPath(user.ID), map[string]any{"name": displayName})
			testsupport.RequireStatus(tt, resp, http.StatusOK)
			testsupport.RequireEqualString(tt, testsupport.RequireStringField(tt, testsupport.RequireDataMap(tt, resp), "name"), displayName, "name")
		})
	})
}
