package tests

import (
	"net/http"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	testsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/http/testkit"
)

func TestUsersMe(t *testing.T) {
	testsupport.Run(t, "current user lifecycle through public HTTP API", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		created := testsupport.CreateIdentityE2EUser(tt, admin, "user")
		user := testsupport.LoginCreatedUser(tt, env, &created)
		userClient := testsupport.NewClient(env, user.Token)
		anonymous := testsupport.NewClient(env, "")
		invalid := testsupport.NewClient(env, "invalid-token")

		cfg.Step("users me rejects anonymous and invalid tokens", func() {
			// /users/me takes the subject from the bearer token, so a missing or forged token must return 401.
			testsupport.RequireErrorCode(tt, anonymous.Get(tt, "/users/me", nil), http.StatusUnauthorized)
			testsupport.RequireErrorCode(tt, invalid.Get(tt, "/users/me", nil), http.StatusUnauthorized)
		})

		cfg.Step("user can read and update own profile", func() {
			me := userClient.Get(tt, "/users/me", nil)
			testsupport.RequireStatus(tt, me, http.StatusOK)
			testsupport.RequireEqualString(tt, testsupport.RequireStringField(tt, testsupport.RequireDataMap(tt, me), "id"), user.ID, "id")

			updatedName := "Self " + testsupport.UniqueE2ESuffix("update")
			updated := userClient.Put(tt, "/users/me", map[string]any{
				"name":     updatedName,
				"lastName": "Updated",
				"gender":   2,
			})
			testsupport.RequireStatus(tt, updated, http.StatusOK)
			testsupport.RequireEqualString(tt, testsupport.RequireStringField(tt, testsupport.RequireDataMap(tt, updated), "name"), updatedName, "name")
		})

		cfg.Step("user can delete own profile and token stops resolving a user", func() {
			// Self-delete verifies the public boundary operation without direct access to the user repository.
			testsupport.RequireStatus(tt, userClient.Delete(tt, "/users/me"), http.StatusNoContent)
			testsupport.RequireErrorCode(tt, userClient.Get(tt, "/users/me", nil), http.StatusUnauthorized)
		})
	})
}
