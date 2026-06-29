package strategies_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	"github.com/google/uuid"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/strategy-registry/tests/e2e/grpc/testkit"
	"google.golang.org/grpc/codes"
)

func TestGRPCStrategiesNotFound(t *testing.T) {
	grpcsupport.Run(t, "missing strategy resources return not found across grpc mutations", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		adminToken := grpcsupport.LoginAdmin(tt, cfg)
		missingStrategyID := uuid.NewString()
		missingAttributeID := uuid.NewString()
		missingValueID := uuid.NewString()

		cfg.Step("get missing strategy returns not found", func() {
			_, err := client.StrategyRegistry().GetStrategy(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.GetStrategyRequest{
				StrategyRef: missingStrategyID,
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.NotFound, "strategy.notFound")
		})

		cfg.Step("update and delete missing strategy return not found", func() {
			_, err := client.StrategyRegistry().UpdateStrategy(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.UpdateStrategyByRefRequest{
				StrategyRef: missingStrategyID,
				Patch:       &strategyregistryv1.UpdateStrategyRequest{Name: new("Ghost")},
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.NotFound, "strategy.notFound")

			_, err = client.StrategyRegistry().DeleteStrategy(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.DeleteStrategyRequest{
				StrategyRef: missingStrategyID,
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.NotFound, "strategy.notFound")
		})

		cfg.Step("missing attribute and value return not found", func() {
			graph := grpcsupport.CreateStrategyGraph(tt, cfg, adminToken)

			_, err := client.StrategyRegistry().GetAttribute(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.GetAttributeRequest{
				StrategyRef:  graph.StrategyID,
				AttributeRef: missingAttributeID,
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.NotFound, "attribute.notFound")

			_, err = client.StrategyRegistry().DeleteAttributeValue(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.DeleteAttributeValueRequest{
				StrategyRef:  graph.StrategyID,
				AttributeRef: graph.AttributeID,
				ValueRef:     missingValueID,
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.NotFound, "attributeValue.notFound")
		})
	})
}
