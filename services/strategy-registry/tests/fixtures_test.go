package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

type e2eStrategyGraph struct {
	StrategyID    string
	StrategySlug  string
	AttributeID   string
	AttributeSlug string
	ValueID       string
	ValueSlug     string
}

func createStrategyGraph(t *testing.T, admin strategyE2EClient) e2eStrategyGraph {
	t.Helper()

	suffix := uuid.NewString()
	strategySlug := fmt.Sprintf("e2e-strategy-%s", suffix[:12])
	attributeSlug := fmt.Sprintf("e2e-param-%s", suffix[:12])
	valueSlug := fmt.Sprintf("e2e-value-%s", suffix[:12])
	strategy := admin.post(t, "/strategies", map[string]any{
		"slug":        strategySlug,
		"name":        "E2E Strategy",
		"description": "created by black-box e2e",
	})
	requireStatus(t, strategy, http.StatusCreated)
	strategyID := requireDataString(t, strategy, "id")
	t.Cleanup(func() {
		cleanup := admin.delete(t, strategyPath(strategyID))
		if cleanup.StatusCode != http.StatusNoContent && cleanup.StatusCode != http.StatusNotFound {
			t.Fatalf("cleanup strategy %s: status=%d body=%s", strategyID, cleanup.StatusCode, cleanup.Body)
		}
	})

	attribute := admin.post(t, strategyPath(strategyID)+"/attributes", map[string]any{
		"slug":        attributeSlug,
		"name":        "E2E Attribute",
		"description": "created by black-box e2e",
	})
	requireStatus(t, attribute, http.StatusCreated)
	attributeID := requireDataString(t, attribute, "id")

	value := admin.post(t, attributePath(strategyID, attributeID)+"/values", map[string]any{
		"slug":      valueSlug,
		"value":     "42",
		"relations": []any{},
	})
	requireStatus(t, value, http.StatusCreated)

	return e2eStrategyGraph{
		StrategyID:    strategyID,
		StrategySlug:  strategySlug,
		AttributeID:   attributeID,
		AttributeSlug: attributeSlug,
		ValueID:       requireDataString(t, value, "id"),
		ValueSlug:     valueSlug,
	}
}
