package tests

import (
	"net/http"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	"github.com/google/uuid"
	testsupport "github.com/stratflow-labs/stratflow/services/strategy-registry/tests/e2e/http/testkit"
)

func TestStrategiesNotFound(t *testing.T) {
	testsupport.Run(t, "missing strategy resources return not found over http", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		missingStrategyID := uuid.NewString()
		missingAttributeID := uuid.NewString()
		missingValueID := uuid.NewString()

		cfg.Step("get missing strategy returns not found", func() {
			resp := admin.Get(tt, testsupport.StrategyPath(missingStrategyID), nil)
			testsupport.RequireStatus(tt, resp, http.StatusNotFound)
		})

		cfg.Step("update missing strategy returns not found", func() {
			resp := admin.Patch(tt, testsupport.StrategyPath(missingStrategyID), map[string]any{"name": "Ghost"})
			testsupport.RequireStatus(tt, resp, http.StatusNotFound)
		})

		cfg.Step("delete missing strategy returns not found", func() {
			resp := admin.Delete(tt, testsupport.StrategyPath(missingStrategyID))
			testsupport.RequireStatus(tt, resp, http.StatusNotFound)
		})

		cfg.Step("missing attribute and value return not found", func() {
			graph := testsupport.CreateStrategyGraph(tt, admin)

			attrResp := admin.Get(tt, testsupport.AttributePath(graph.StrategyID, missingAttributeID), nil)
			testsupport.RequireStatus(tt, attrResp, http.StatusNotFound)

			valueResp := admin.Delete(tt, testsupport.ValuePath(graph.StrategyID, graph.AttributeID, missingValueID))
			testsupport.RequireStatus(tt, valueResp, http.StatusNotFound)
		})
	})
}
