package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	testsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/http/testkit"
)

func TestUsersDelete(t *testing.T) {
	testsupport.Run(t, "delete user removes it from reads and filtered lists", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		created := testsupport.CreateIdentityE2EUser(tt, admin, "user")
		manager := testsupport.LoginCreatedUser(tt, env, new(testsupport.CreateIdentityE2EUser(tt, admin, "manager")))
		anonymous := testsupport.NewClient(env, "")
		invalid := testsupport.NewClient(env, "invalid-token")
		managerClient := testsupport.NewClient(env, manager.Token)

		cfg.Step("delete user rejects anonymous and forged tokens", func() {
			testsupport.RequireErrorCode(tt, anonymous.Delete(tt, testsupport.UserPath(created.ID)), http.StatusUnauthorized)
			testsupport.RequireErrorCode(tt, invalid.Delete(tt, testsupport.UserPath(created.ID)), http.StatusUnauthorized)
		})

		cfg.Step("delete user rejects non-admin foreign deletion", func() {
			testsupport.RequireErrorCode(tt, managerClient.Delete(tt, testsupport.UserPath(created.ID)), http.StatusForbidden)
		})

		cfg.Step("delete user by id succeeds", func() {
			resp := admin.Delete(tt, testsupport.UserPath(created.ID))
			testsupport.RequireStatus(tt, resp, http.StatusNoContent)
		})

		cfg.Step("subsequent get returns not found", func() {
			resp := admin.Get(tt, testsupport.UserPath(created.ID), nil)
			testsupport.RequireStatus(tt, resp, http.StatusNotFound)
		})

		cfg.Step("filtered list no longer includes the deleted user", func() {
			resp := admin.Get(tt, "/users", url.Values{"search": []string{created.Email}})
			testsupport.RequireStatus(tt, resp, http.StatusOK)
			data := testsupport.RequireDataMap(tt, resp)
			items, ok := data["items"].([]any)
			if !ok {
				t.Fatalf("expected items array, got %v", data["items"])
			}
			if len(items) != 0 {
				t.Fatalf("expected no items after delete, got %v", items)
			}
		})
	})
}
