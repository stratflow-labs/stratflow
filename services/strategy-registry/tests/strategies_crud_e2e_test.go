package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nikita-Filonov/axiom"
)

func TestStrategyRegistryStrategiesE2E_AdminCRUDGraph(t *testing.T) {
	runStrategyE2E(t, "admin manages strategy graph through public HTTP API", func(cfg *axiom.Config, env strategyE2EEnv) {
		tt := e2ET(t, cfg)
		admin := adminIdentityClient(tt, env)
		graph := createStrategyGraph(tt, admin)

		cfg.Step("read and list created graph", func() {
			requireStatus(tt, admin.get(tt, strategyPath(graph.StrategyID), nil), http.StatusOK)
			requireStatus(tt, admin.get(tt, attributePath(graph.StrategyID, graph.AttributeID), nil), http.StatusOK)
			requireListTotalAtLeast(tt, admin.get(tt, strategyPath(graph.StrategyID)+"/attributes", nil), 1)
			requireListTotal(tt, admin.get(tt, "/strategies", url.Values{"search": []string{graph.StrategySlug}, "page": []string{"1"}, "pageSize": []string{"10"}, "sort": []string{"created_at_desc"}}), 1)
			requireListTotal(tt, admin.get(tt, strategyPath(graph.StrategyID)+"/attributes", url.Values{"search": []string{graph.AttributeSlug}, "page": []string{"1"}, "pageSize": []string{"10"}, "sort": []string{"created_at_desc"}}), 1)
			requireListTotal(tt, admin.get(tt, attributePath(graph.StrategyID, graph.AttributeID)+"/values", url.Values{"search": []string{graph.ValueSlug}, "page": []string{"1"}, "pageSize": []string{"10"}, "sort": []string{"created_at_desc"}}), 1)
			requireListTotalAtLeast(tt, admin.get(tt, "/strategies/all", nil), 1)
		})

		cfg.Step("update strategy, attribute and value", func() {
			// Updates verify the external contract of patch endpoints and the risk of partial graph data loss.
			updatedStrategy := admin.patch(tt, strategyPath(graph.StrategyID), map[string]any{"name": "E2E Strategy v2"})
			requireStatus(tt, updatedStrategy, http.StatusOK)
			updatedAttribute := admin.patch(tt, attributePath(graph.StrategyID, graph.AttributeID), map[string]any{"name": "E2E Attribute v2"})
			requireStatus(tt, updatedAttribute, http.StatusOK)
			updatedValue := admin.patch(tt, valuePath(graph.StrategyID, graph.AttributeID, graph.ValueID), map[string]any{"value": `<script>alert("strategy")</script>`, "relations": []any{}})
			requireStatus(tt, updatedValue, http.StatusOK)
		})

		cfg.Step("clone graph and reject invalid query attributes", func() {
			clone := admin.post(tt, "/strategies/clone", map[string]any{
				"items": []map[string]any{{"sourceStrategyId": graph.StrategyID, "slug": graph.StrategySlug + "-clone"}},
			})
			requireStatus(tt, clone, http.StatusCreated)
			requireErrorStatus(tt, admin.get(tt, "/strategies", url.Values{"pageSize": []string{"1000"}}), http.StatusBadRequest)
			requireErrorStatus(tt, admin.get(tt, "/strategies", url.Values{"sort": []string{"unknown"}}), http.StatusBadRequest)
		})

		cfg.Step("delete value, attribute and strategy", func() {
			requireStatus(tt, admin.delete(tt, valuePath(graph.StrategyID, graph.AttributeID, graph.ValueID)), http.StatusNoContent)
			requireStatus(tt, admin.delete(tt, attributePath(graph.StrategyID, graph.AttributeID)), http.StatusNoContent)
			requireStatus(tt, admin.delete(tt, strategyPath(graph.StrategyID)), http.StatusNoContent)
			requireErrorStatus(tt, admin.get(tt, strategyPath(graph.StrategyID), nil), http.StatusNotFound)
		})
	})
}
