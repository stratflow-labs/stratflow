package tests

import (
	"net/http"
	"testing"

	"github.com/Nikita-Filonov/axiom"
)

func TestStrategyRegistryStrategiesE2E_RoleAccessMatrix(t *testing.T) {
	runStrategyE2E(t, "strategy registry role access matrix", func(cfg *axiom.Config, env strategyE2EEnv) {
		tt := e2ET(t, cfg)
		identityAdmin := adminIdentityClient(tt, env)
		user := createIdentityUser(tt, env, identityAdmin, "user")
		manager := createIdentityUser(tt, env, identityAdmin, "manager")

		admin := newStrategyE2EClient(env, identityAdmin.token)
		userClient := newStrategyE2EClient(env, user.Token)
		managerClient := newStrategyE2EClient(env, manager.Token)
		anonymous := newStrategyE2EClient(env, "")
		invalid := newStrategyE2EClient(env, "invalid-token")
		graph := createStrategyGraph(tt, admin)

		cfg.Step("anonymous and invalid tokens cannot access protected endpoints", func() {
			// Protected strategy endpoints must not expose the catalog without a valid bearer token.
			requireErrorStatus(tt, anonymous.get(tt, "/strategies", nil), http.StatusUnauthorized)
			requireErrorStatus(tt, invalid.get(tt, "/strategies", nil), http.StatusUnauthorized)
		})

		cfg.Step("user has no strategy permissions", func() {
			// Identity-backed verification currently maps user tokens without strategy claims to 401;
			// if the token is verified, the policy is expected to return 403.
			requireErrorStatusOneOf(tt, userClient.get(tt, "/strategies", nil), http.StatusUnauthorized, http.StatusForbidden)
			requireErrorStatusOneOf(tt, userClient.get(tt, strategyPath(graph.StrategyID), nil), http.StatusUnauthorized, http.StatusForbidden)
			requireErrorStatusOneOf(tt, userClient.post(tt, "/strategies", map[string]any{"slug": "forbidden", "name": "Forbidden", "description": "Forbidden"}), http.StatusUnauthorized, http.StatusForbidden)
		})

		cfg.Step("manager can read but cannot mutate", func() {
			// The manager scope in policy.csv is read-only; the e2e test locks in the denial of admin-only mutations.
			requireStatus(tt, managerClient.get(tt, "/strategies", nil), http.StatusOK)
			requireStatus(tt, managerClient.get(tt, strategyPath(graph.StrategyID), nil), http.StatusOK)
			requireStatus(tt, managerClient.get(tt, strategyPath(graph.StrategyID)+"/attributes", nil), http.StatusOK)
			requireErrorStatus(tt, managerClient.patch(tt, strategyPath(graph.StrategyID), map[string]any{"name": "forbidden"}), http.StatusForbidden)
			requireErrorStatus(tt, managerClient.delete(tt, strategyPath(graph.StrategyID)), http.StatusForbidden)
		})
	})
}
