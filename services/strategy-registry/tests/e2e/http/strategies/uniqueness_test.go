package tests

import (
	"net/http"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
	testsupport "github.com/stratflow-labs/stratflow/services/strategy-registry/tests/e2e/http/testkit"
)

func TestStrategiesUniqueness(t *testing.T) {
	testsupport.Run(t, "strategy graph creation enforces unique slugs over http", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		graph := testsupport.CreateStrategyGraph(tt, admin)

		cfg.Step("duplicate strategy slug is rejected", func() {
			resp := admin.Post(tt, "/strategies", map[string]any{
				"slug":        graph.StrategySlug,
				"name":        "Duplicate Strategy",
				"description": "Duplicate Strategy",
			})
			testsupport.RequireStatus(tt, resp, http.StatusConflict)
			testsupport.RequireErrorCode(tt, resp, http.StatusConflict)
		})

		cfg.Step("duplicate attribute slug under the same strategy is rejected", func() {
			resp := admin.Post(tt, testsupport.StrategyPath(graph.StrategyID)+"/attributes", map[string]any{
				"slug":        graph.AttributeSlug,
				"name":        "Duplicate Attribute",
				"description": "Duplicate Attribute",
			})
			testsupport.RequireStatus(tt, resp, http.StatusConflict)
			testsupport.RequireErrorCode(tt, resp, http.StatusConflict)
		})

		cfg.Step("duplicate value slug under the same attribute is rejected", func() {
			resp := admin.Post(tt, testsupport.AttributePath(graph.StrategyID, graph.AttributeID)+"/values", map[string]any{
				"slug":      graph.ValueSlug,
				"value":     "duplicate",
				"relations": []any{},
			})
			testsupport.RequireStatus(tt, resp, http.StatusConflict)
			testsupport.RequireErrorCode(tt, resp, http.StatusConflict)
		})

		cfg.Step("clone request requires unique destination slugs in one payload", func() {
			cloneSlug := "duplicate-clone-" + e2ecommon.UniqueSuffix("strategy")
			resp := admin.Post(tt, "/strategies/clone", map[string]any{
				"items": []map[string]any{
					{"sourceStrategyId": graph.StrategyID, "slug": cloneSlug},
					{"sourceStrategyId": graph.StrategyID, "slug": cloneSlug},
				},
			})
			testsupport.RequireStatus(tt, resp, http.StatusBadRequest)
			testsupport.RequireErrorCode(tt, resp, http.StatusBadRequest)
		})
	})
}
