package testkit

import (
	"fmt"
	"net/http"
	"testing"

	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
)

type StrategyGraph struct {
	StrategyID    string
	StrategySlug  string
	AttributeID   string
	AttributeSlug string
	ValueID       string
	ValueSlug     string
}

func CreateStrategyGraph(t *testing.T, admin Client) StrategyGraph {
	t.Helper()

	suffix := e2ecommon.UniqueSuffix("strategy")
	strategySlug := fmt.Sprintf("e2e-strategy-%s", suffix)
	attributeSlug := fmt.Sprintf("e2e-param-%s", suffix)
	valueSlug := fmt.Sprintf("e2e-value-%s", suffix)

	strategy := admin.Post(t, "/strategies", map[string]any{
		"slug":        strategySlug,
		"name":        "E2E Strategy",
		"description": "created by strategy-registry e2e",
	})
	RequireStatus(t, strategy, http.StatusCreated)
	strategyID := RequireDataString(t, strategy, "id")
	t.Cleanup(func() {
		cleanup := admin.Delete(t, StrategyPath(strategyID))
		if cleanup.StatusCode != http.StatusNoContent && cleanup.StatusCode != http.StatusNotFound {
			t.Fatalf("cleanup strategy %s: status=%d body=%s", strategyID, cleanup.StatusCode, cleanup.Body)
		}
	})

	attribute := admin.Post(t, StrategyPath(strategyID)+"/attributes", map[string]any{
		"slug":        attributeSlug,
		"name":        "E2E Attribute",
		"description": "created by strategy-registry e2e",
	})
	RequireStatus(t, attribute, http.StatusCreated)
	attributeID := RequireDataString(t, attribute, "id")

	value := admin.Post(t, AttributePath(strategyID, attributeID)+"/values", map[string]any{
		"slug":      valueSlug,
		"value":     "42",
		"relations": []any{},
	})
	RequireStatus(t, value, http.StatusCreated)

	return StrategyGraph{
		StrategyID:    strategyID,
		StrategySlug:  strategySlug,
		AttributeID:   attributeID,
		AttributeSlug: attributeSlug,
		ValueID:       RequireDataString(t, value, "id"),
		ValueSlug:     valueSlug,
	}
}
