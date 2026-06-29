package tests

import (
	"net/http"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	testsupport "github.com/stratflow-labs/stratflow/services/strategy-registry/tests/e2e/http/testkit"
)

func TestStrategiesAccess(t *testing.T) {
	testsupport.Run(t, "strategy registry role access matrix", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		user := testsupport.CreateIdentityUser(tt, env, admin, "user")
		manager := testsupport.CreateIdentityUser(tt, env, admin, "manager")
		graph := testsupport.CreateStrategyGraph(tt, admin)

		anonymous := testsupport.NewClient(env, "")
		invalid := testsupport.NewClient(env, "invalid-token")
		userClient := testsupport.NewClient(env, user.Token)
		managerClient := testsupport.NewClient(env, manager.Token)

		cfg.Step("protected endpoints reject anonymous and invalid tokens", func() {
			testsupport.RequireErrorCode(tt, anonymous.Get(tt, "/strategies", nil), http.StatusUnauthorized)
			testsupport.RequireErrorCode(tt, invalid.Get(tt, "/strategies", nil), http.StatusUnauthorized)
		})

		cfg.Step("regular user has no strategy permissions", func() {
			testsupport.RequireErrorCodeOneOf(tt, userClient.Get(tt, "/strategies", nil), http.StatusUnauthorized, http.StatusForbidden)
			testsupport.RequireErrorCodeOneOf(tt, userClient.Get(tt, testsupport.StrategyPath(graph.StrategyID), nil), http.StatusUnauthorized, http.StatusForbidden)
			testsupport.RequireErrorCodeOneOf(tt, userClient.Post(tt, "/strategies", map[string]any{
				"slug":        "forbidden-user-create",
				"name":        "Forbidden",
				"description": "Forbidden",
			}), http.StatusUnauthorized, http.StatusForbidden)
		})

		cfg.Step("manager can read but cannot mutate", func() {
			testsupport.RequireStatus(tt, managerClient.Get(tt, "/strategies", nil), http.StatusOK)
			testsupport.RequireStatus(tt, managerClient.Get(tt, testsupport.StrategyPath(graph.StrategyID), nil), http.StatusOK)
			testsupport.RequireStatus(tt, managerClient.Get(tt, testsupport.StrategyPath(graph.StrategyID)+"/attributes", nil), http.StatusOK)
			testsupport.RequireErrorCode(tt, managerClient.Patch(tt, testsupport.StrategyPath(graph.StrategyID), map[string]any{"name": "forbidden"}), http.StatusForbidden)
			testsupport.RequireErrorCode(tt, managerClient.Patch(tt, testsupport.GraphPath(graph.StrategyID), map[string]any{
				"actions": []map[string]any{
					{"createAttribute": map[string]any{"slug": "forbidden", "name": "Forbidden", "description": "Forbidden"}},
				},
			}), http.StatusForbidden)
			testsupport.RequireErrorCode(tt, managerClient.Delete(tt, testsupport.StrategyPath(graph.StrategyID)), http.StatusForbidden)
		})
	})
}
