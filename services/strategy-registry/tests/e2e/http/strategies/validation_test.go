package tests

import (
	"net/http"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	"github.com/google/uuid"
	testsupport "github.com/stratflow-labs/stratflow/services/strategy-registry/tests/e2e/http/testkit"
)

func TestStrategiesValidation(t *testing.T) {
	testsupport.Run(t, "strategy registry validation returns stable client errors", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		graph := testsupport.CreateStrategyGraph(tt, admin)

		cfg.Step("create strategy rejects empty normalized fields", func() {
			resp := admin.Post(tt, "/strategies", map[string]any{
				"slug":        "   ",
				"name":        "Invalid",
				"description": "Invalid",
			})
			testsupport.RequireErrorCode(tt, resp, http.StatusBadRequest)
		})

		cfg.Step("create attribute under missing strategy returns not found", func() {
			resp := admin.Post(tt, testsupport.StrategyPath(uuid.NewString())+"/attributes", map[string]any{
				"slug":        "missing-strategy-attribute",
				"name":        "Missing Strategy",
				"description": "Missing Strategy",
			})
			testsupport.RequireErrorCode(tt, resp, http.StatusNotFound)
		})

		cfg.Step("relation mismatch is a client error", func() {
			other := testsupport.CreateStrategyGraph(tt, admin)
			resp := admin.Patch(tt, testsupport.ValuePath(graph.StrategyID, graph.AttributeID, graph.ValueID), map[string]any{
				"relations": []map[string]any{{
					"toAttributeId": other.AttributeID,
					"toValueId":     other.ValueID,
				}},
			})
			testsupport.RequireErrorCode(tt, resp, http.StatusBadRequest)
		})
	})
}
