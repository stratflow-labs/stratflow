package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	testsupport "github.com/stratflow-labs/stratflow/services/strategy-registry/tests/e2e/http/testkit"
)

func TestStrategiesCRUDGraph(t *testing.T) {
	testsupport.Run(t, "admin manages strategy graph through public http api", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		graph := testsupport.CreateStrategyGraph(tt, admin)

		cfg.Step("read and list created graph", func() {
			testsupport.RequireStatus(tt, admin.Get(tt, testsupport.StrategyPath(graph.StrategyID), nil), http.StatusOK)
			testsupport.RequireStatus(tt, admin.Get(tt, testsupport.AttributePath(graph.StrategyID, graph.AttributeID), nil), http.StatusOK)
			testsupport.RequireListTotalAtLeast(tt, admin.Get(tt, testsupport.StrategyPath(graph.StrategyID)+"/attributes", nil), 1)
			testsupport.RequireListTotal(tt, admin.Get(tt, "/strategies", url.Values{
				"search":   []string{graph.StrategySlug},
				"page":     []string{"1"},
				"pageSize": []string{"10"},
				"sort":     []string{"created_at_desc"},
			}), 1)
			testsupport.RequireListTotal(tt, admin.Get(tt, testsupport.StrategyPath(graph.StrategyID)+"/attributes", url.Values{
				"search":   []string{graph.AttributeSlug},
				"page":     []string{"1"},
				"pageSize": []string{"10"},
				"sort":     []string{"created_at_desc"},
			}), 1)
			testsupport.RequireListTotal(tt, admin.Get(tt, testsupport.AttributePath(graph.StrategyID, graph.AttributeID)+"/values", url.Values{
				"search":   []string{graph.ValueSlug},
				"page":     []string{"1"},
				"pageSize": []string{"10"},
				"sort":     []string{"created_at_desc"},
			}), 1)
			testsupport.RequireListTotalAtLeast(tt, admin.Get(tt, "/strategies/all", nil), 1)
		})

		cfg.Step("update strategy, attribute and value", func() {
			updatedStrategy := admin.Patch(tt, testsupport.StrategyPath(graph.StrategyID), map[string]any{"name": "E2E Strategy v2"})
			testsupport.RequireStatus(tt, updatedStrategy, http.StatusOK)
			testsupport.RequireEqualString(tt, testsupport.RequireDataString(tt, updatedStrategy, "name"), "E2E Strategy v2", "name")

			updatedAttribute := admin.Patch(tt, testsupport.AttributePath(graph.StrategyID, graph.AttributeID), map[string]any{"name": "E2E Attribute v2"})
			testsupport.RequireStatus(tt, updatedAttribute, http.StatusOK)
			testsupport.RequireEqualString(tt, testsupport.RequireDataString(tt, updatedAttribute, "name"), "E2E Attribute v2", "name")

			const xssLikeValue = `<script>alert("strategy")</script>`
			updatedValue := admin.Patch(tt, testsupport.ValuePath(graph.StrategyID, graph.AttributeID, graph.ValueID), map[string]any{
				"value":     xssLikeValue,
				"relations": []any{},
			})
			testsupport.RequireStatus(tt, updatedValue, http.StatusOK)
			testsupport.RequireEqualString(tt, testsupport.RequireDataString(tt, updatedValue, "value"), xssLikeValue, "value")
		})

		cfg.Step("batch graph action updates the graph", func() {
			resp := admin.Patch(tt, testsupport.GraphPath(graph.StrategyID), map[string]any{
				"actions": []map[string]any{
					{
						"createAttribute": map[string]any{
							"slug":        graph.AttributeSlug + "-batch",
							"name":        "Batch Attribute",
							"description": "created by batch graph e2e",
						},
					},
				},
			})
			testsupport.RequireStatus(tt, resp, http.StatusOK)
			data := testsupport.RequireDataMap(tt, resp)
			if _, ok := data["parameters"].([]any); !ok {
				t.Fatalf("expected graph response parameters array, got %v", data)
			}
		})

		cfg.Step("clone graph and reject invalid list query", func() {
			clone := admin.Post(tt, "/strategies/clone", map[string]any{
				"items": []map[string]any{{
					"sourceStrategyId": graph.StrategyID,
					"slug":             graph.StrategySlug + "-clone",
				}},
			})
			testsupport.RequireStatus(tt, clone, http.StatusCreated)
			testsupport.RequireErrorCode(tt, admin.Get(tt, "/strategies", url.Values{"pageSize": []string{"1000"}}), http.StatusBadRequest)
			testsupport.RequireErrorCode(tt, admin.Get(tt, "/strategies", url.Values{"sort": []string{"unknown"}}), http.StatusBadRequest)
		})

		cfg.Step("delete value, attribute and strategy", func() {
			testsupport.RequireStatus(tt, admin.Delete(tt, testsupport.ValuePath(graph.StrategyID, graph.AttributeID, graph.ValueID)), http.StatusNoContent)
			testsupport.RequireStatus(tt, admin.Delete(tt, testsupport.AttributePath(graph.StrategyID, graph.AttributeID)), http.StatusNoContent)
			testsupport.RequireStatus(tt, admin.Delete(tt, testsupport.StrategyPath(graph.StrategyID)), http.StatusNoContent)
			testsupport.RequireErrorCode(tt, admin.Get(tt, testsupport.StrategyPath(graph.StrategyID), nil), http.StatusNotFound)
		})
	})
}
