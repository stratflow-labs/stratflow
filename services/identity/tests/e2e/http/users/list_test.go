package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	testsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/http/testkit"
	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
)

func TestUsersList(t *testing.T) {
	testsupport.Run(t, "users list supports search and bounded page size over http", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		created := testsupport.CreateIdentityE2EUser(tt, admin, "user")

		cfg.Step("search narrows results to the created user", func() {
			q := url.Values{"search": []string{created.Email}}
			resp := admin.Get(tt, "/users", q)
			testsupport.RequireStatus(tt, resp, http.StatusOK)
			data := testsupport.RequireDataMap(tt, resp)
			items, ok := data["items"].([]any)
			if !ok || len(items) == 0 {
				t.Fatalf("expected non-empty items array, got %v", data["items"])
			}
		})

		cfg.Step("oversized pageSize is rejected by the public http contract", func() {
			q := url.Values{"pageSize": []string{"1000"}}
			resp := admin.Get(tt, "/users", q)
			testsupport.RequireStatus(tt, resp, http.StatusBadRequest)
		})

		cfg.Step("pagination splits a filtered result set across pages", func() {
			marker := "pagination-" + e2ecommon.UniqueSuffix("http")
			first := testsupport.CreateIdentityE2EUser(tt, admin, "user")
			second := testsupport.CreateIdentityE2EUser(tt, admin, "user")

			firstUpdate := admin.Put(tt, testsupport.UserPath(first.ID), map[string]any{"name": marker + "-first"})
			testsupport.RequireStatus(tt, firstUpdate, http.StatusOK)

			secondUpdate := admin.Put(tt, testsupport.UserPath(second.ID), map[string]any{"name": marker + "-second"})
			testsupport.RequireStatus(tt, secondUpdate, http.StatusOK)

			firstPage := admin.Get(tt, "/users", url.Values{
				"search":   []string{marker},
				"page":     []string{"1"},
				"pageSize": []string{"1"},
				"sort":     []string{"created_at ASC"},
			})
			testsupport.RequireStatus(tt, firstPage, http.StatusOK)
			firstData := testsupport.RequireDataMap(tt, firstPage)
			firstItems, ok := firstData["items"].([]any)
			if !ok || len(firstItems) != 1 {
				t.Fatalf("expected first page to contain exactly one item, got %v", firstData["items"])
			}

			secondPage := admin.Get(tt, "/users", url.Values{
				"search":   []string{marker},
				"page":     []string{"2"},
				"pageSize": []string{"1"},
				"sort":     []string{"created_at ASC"},
			})
			testsupport.RequireStatus(tt, secondPage, http.StatusOK)
			secondData := testsupport.RequireDataMap(tt, secondPage)
			secondItems, ok := secondData["items"].([]any)
			if !ok || len(secondItems) != 1 {
				t.Fatalf("expected second page to contain exactly one item, got %v", secondData["items"])
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
				t.Fatalf("expected different user ids across pages, got identical id %v", firstItem["id"])
			}
		})
	})
}
