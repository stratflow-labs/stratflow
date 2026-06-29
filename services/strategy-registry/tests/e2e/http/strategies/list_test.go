package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
	testsupport "github.com/stratflow-labs/stratflow/services/strategy-registry/tests/e2e/http/testkit"
)

func TestStrategiesList(t *testing.T) {
	testsupport.Run(t, "strategies list supports search and pagination over http", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		graph := testsupport.CreateStrategyGraph(tt, admin)

		cfg.Step("search narrows results to the created strategy", func() {
			resp := admin.Get(tt, "/strategies", url.Values{"search": []string{graph.StrategySlug}})
			testsupport.RequireStatus(tt, resp, http.StatusOK)
			data := testsupport.RequireDataMap(tt, resp)
			items, ok := data["items"].([]any)
			if !ok || len(items) == 0 {
				t.Fatalf("expected non-empty items array, got %v", data["items"])
			}
		})

		cfg.Step("pagination splits a filtered result set across pages", func() {
			marker := "pagination-" + e2ecommon.UniqueSuffix("strategy-http")
			first := admin.Post(tt, "/strategies", map[string]any{
				"slug":        marker + "-first",
				"name":        marker + "-first",
				"description": "pagination marker",
			})
			testsupport.RequireStatus(tt, first, http.StatusCreated)

			second := admin.Post(tt, "/strategies", map[string]any{
				"slug":        marker + "-second",
				"name":        marker + "-second",
				"description": "pagination marker",
			})
			testsupport.RequireStatus(tt, second, http.StatusCreated)

			firstID := testsupport.RequireDataString(tt, first, "id")
			secondID := testsupport.RequireDataString(tt, second, "id")
			t.Cleanup(func() { _ = admin.Delete(tt, testsupport.StrategyPath(firstID)) })
			t.Cleanup(func() { _ = admin.Delete(tt, testsupport.StrategyPath(secondID)) })

			firstPage := admin.Get(tt, "/strategies", url.Values{
				"search":   []string{marker},
				"page":     []string{"1"},
				"pageSize": []string{"1"},
				"sort":     []string{"created_at_asc"},
			})
			testsupport.RequireStatus(tt, firstPage, http.StatusOK)
			firstItems, ok := testsupport.RequireDataMap(tt, firstPage)["items"].([]any)
			if !ok || len(firstItems) != 1 {
				t.Fatalf("expected first page to contain exactly one item, got %v", testsupport.RequireDataMap(tt, firstPage)["items"])
			}

			secondPage := admin.Get(tt, "/strategies", url.Values{
				"search":   []string{marker},
				"page":     []string{"2"},
				"pageSize": []string{"1"},
				"sort":     []string{"created_at_asc"},
			})
			testsupport.RequireStatus(tt, secondPage, http.StatusOK)
			secondItems, ok := testsupport.RequireDataMap(tt, secondPage)["items"].([]any)
			if !ok || len(secondItems) != 1 {
				t.Fatalf("expected second page to contain exactly one item, got %v", testsupport.RequireDataMap(tt, secondPage)["items"])
			}

			firstItem, ok := firstItems[0].(map[string]any)
			if !ok {
				t.Fatalf("expected first page item object, got %T", firstItems[0])
			}
			secondItem, ok := secondItems[0].(map[string]any)
			if !ok {
				t.Fatalf("expected second page item object, got %T", secondItems[0])
			}
			if firstItem["id"] == secondItem["id"] {
				t.Fatalf("expected different strategy ids across pages, got identical id %v", firstItem["id"])
			}
		})
	})
}
