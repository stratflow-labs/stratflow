package strategies_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	"github.com/google/uuid"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/strategy-registry/tests/e2e/grpc/testkit"
	"google.golang.org/grpc/codes"
)

func TestGRPCStrategiesValidation(t *testing.T) {
	grpcsupport.Run(t, "strategy registry grpc validation returns stable domain reasons", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		adminToken := grpcsupport.LoginAdmin(tt, cfg)
		graph := grpcsupport.CreateStrategyGraph(tt, cfg, adminToken)

		cfg.Step("create strategy rejects empty normalized fields", func() {
			_, err := client.StrategyRegistry().CreateStrategy(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.CreateStrategyRequest{
				Slug:        "   ",
				Name:        "Invalid",
				Description: "Invalid",
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.InvalidArgument, "registry.slugEmpty")
		})

		cfg.Step("create attribute under missing strategy returns not found", func() {
			_, err := client.StrategyRegistry().CreateAttribute(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.CreateAttributeRequest{
				StrategyRef: uuid.NewString(),
				Slug:        "missing-strategy-attribute",
				Name:        "Missing Strategy",
				Description: "Missing Strategy",
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.NotFound, "strategy.notFound")
		})

		cfg.Step("empty updates are rejected", func() {
			_, err := client.StrategyRegistry().UpdateStrategy(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.UpdateStrategyByRefRequest{
				StrategyRef: graph.StrategyID,
				Patch:       &strategyregistryv1.UpdateStrategyRequest{},
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.InvalidArgument, "strategy.updateEmpty")

			_, err = client.StrategyRegistry().UpdateAttribute(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.UpdateAttributeRequest{
				StrategyRef:  graph.StrategyID,
				AttributeRef: graph.AttributeID,
				Patch:        &strategyregistryv1.UpdateAttributePatch{},
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.InvalidArgument, "attribute.updateEmpty")

			_, err = client.StrategyRegistry().UpdateAttributeValue(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.UpdateAttributeValueRequest{
				StrategyRef:  graph.StrategyID,
				AttributeRef: graph.AttributeID,
				ValueRef:     graph.ValueID,
				Patch:        &strategyregistryv1.UpdateAttributeValuePatch{},
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.InvalidArgument, "attributeValue.updateEmpty")
		})

		cfg.Step("relation mismatch is a client error", func() {
			other := grpcsupport.CreateStrategyGraph(tt, cfg, adminToken)
			_, err := client.StrategyRegistry().UpdateAttributeValue(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.UpdateAttributeValueRequest{
				StrategyRef:  graph.StrategyID,
				AttributeRef: graph.AttributeID,
				ValueRef:     graph.ValueID,
				Patch: &strategyregistryv1.UpdateAttributeValuePatch{
					Relations: []*strategyregistryv1.UpdateAttributeValueRelationInput{{
						ToAttributeId: other.AttributeID,
						ToValueId:     other.ValueID,
					}},
				},
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.InvalidArgument, "attributeValue.relationCombinationNotFound")
		})

		cfg.Step("duplicate and self-reference relations are rejected", func() {
			other := grpcsupport.CreateStrategyGraph(tt, cfg, adminToken)

			_, err := client.StrategyRegistry().UpdateAttributeValue(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.UpdateAttributeValueRequest{
				StrategyRef:  graph.StrategyID,
				AttributeRef: graph.AttributeID,
				ValueRef:     graph.ValueID,
				Patch: &strategyregistryv1.UpdateAttributeValuePatch{
					Relations: []*strategyregistryv1.UpdateAttributeValueRelationInput{
						{
							ToAttributeId: other.AttributeID,
							ToValueId:     other.ValueID,
						},
						{
							ToAttributeId: other.AttributeID,
							ToValueId:     other.ValueID,
						},
					},
				},
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.InvalidArgument, "attributeValue.relationDuplicate")

			_, err = client.StrategyRegistry().UpdateAttributeValue(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.UpdateAttributeValueRequest{
				StrategyRef:  graph.StrategyID,
				AttributeRef: graph.AttributeID,
				ValueRef:     graph.ValueID,
				Patch: &strategyregistryv1.UpdateAttributeValuePatch{
					Relations: []*strategyregistryv1.UpdateAttributeValueRelationInput{
						{
							ToAttributeId: graph.AttributeID,
							ToValueId:     graph.ValueID,
						},
					},
				},
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.InvalidArgument, "attributeValue.relationSelfReference")
		})
	})
}
