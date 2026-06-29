package strategies_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/strategy-registry/tests/e2e/grpc/testkit"
	"google.golang.org/grpc/codes"
)

func TestGRPCStrategiesUniqueness(t *testing.T) {
	grpcsupport.Run(t, "strategy graph creation enforces unique slugs over grpc", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		adminToken := grpcsupport.LoginAdmin(tt, cfg)
		graph := grpcsupport.CreateStrategyGraph(tt, cfg, adminToken)

		cfg.Step("duplicate strategy slug is rejected", func() {
			_, err := client.StrategyRegistry().CreateStrategy(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.CreateStrategyRequest{
				Slug:        graph.StrategySlug,
				Name:        "Duplicate Strategy",
				Description: "Duplicate Strategy",
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.AlreadyExists, "strategy.alreadyExists")
		})

		cfg.Step("duplicate attribute slug under the same strategy is rejected", func() {
			_, err := client.StrategyRegistry().CreateAttribute(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.CreateAttributeRequest{
				StrategyRef: graph.StrategyID,
				Slug:        graph.AttributeSlug,
				Name:        "Duplicate Attribute",
				Description: "Duplicate Attribute",
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.AlreadyExists, "attribute.alreadyExists")
		})

		cfg.Step("duplicate value slug under the same attribute is rejected", func() {
			_, err := client.StrategyRegistry().CreateAttributeValue(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.CreateAttributeValueRequest{
				StrategyRef:  graph.StrategyID,
				AttributeRef: graph.AttributeID,
				Slug:         graph.ValueSlug,
				Value:        "duplicate",
				Relations:    []*strategyregistryv1.CreateAttributeValueRelationInput{},
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.AlreadyExists, "attributeValue.alreadyExists")
		})

		cfg.Step("clone request requires unique destination slugs in one payload", func() {
			cloneSlug := "duplicate-clone-" + e2ecommon.UniqueSuffix("strategy")
			_, err := client.StrategyRegistry().CloneStrategies(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.CloneStrategiesRequest{
				Items: []*strategyregistryv1.CloneStrategyItemInput{
					{SourceStrategyId: graph.StrategyID, Slug: cloneSlug},
					{SourceStrategyId: graph.StrategyID, Slug: cloneSlug},
				},
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.InvalidArgument, "strategy.cloneDuplicateSlug")
		})
	})
}
