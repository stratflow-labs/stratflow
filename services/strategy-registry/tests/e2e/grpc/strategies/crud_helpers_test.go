package strategies_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/strategy-registry/tests/e2e/grpc/testkit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func runStrategiesCRUDGraph(t *testing.T, role string) {
	t.Helper()

	grpcsupport.Run(t, role+" manages strategy graph over grpc", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		adminToken := grpcsupport.LoginAdmin(tt, cfg)

		actorToken := adminToken
		if role != "admin" {
			user := grpcsupport.CreateIdentityUser(tt, cfg, adminToken, role)
			actorToken = user.Token
		}

		graph := grpcsupport.CreateStrategyGraph(tt, cfg, actorToken)

		cfg.Step("read and list created graph", func() {
			strategyResp, err := client.StrategyRegistry().GetStrategy(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.GetStrategyRequest{
				StrategyRef: graph.StrategyID,
			})
			require.NoError(tt, err)
			require.Equal(tt, graph.StrategySlug, strategyResp.GetData().GetSlug())

			attributeResp, err := client.StrategyRegistry().GetAttribute(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.GetAttributeRequest{
				StrategyRef:  graph.StrategyID,
				AttributeRef: graph.AttributeID,
			})
			require.NoError(tt, err)
			require.Equal(tt, graph.AttributeSlug, attributeResp.GetData().GetSlug())
			require.NotEmpty(tt, attributeResp.GetData().GetValues())

			listResp, err := client.StrategyRegistry().ListStrategies(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.ListStrategiesRequest{
				Search:   new(graph.StrategySlug),
				Page:     new(int32(1)),
				PageSize: new(int32(10)),
				Sort:     new("created_at_desc"),
			})
			require.NoError(tt, err)
			require.EqualValues(tt, 1, listResp.GetData().GetTotal())

			attrsResp, err := client.StrategyRegistry().ListAttributes(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.ListAttributesRequest{
				StrategyRef: graph.StrategyID,
				Search:      new(graph.AttributeSlug),
				Page:        new(int32(1)),
				PageSize:    new(int32(10)),
				Sort:        new("created_at_desc"),
			})
			require.NoError(tt, err)
			require.EqualValues(tt, 1, attrsResp.GetData().GetTotal())
		})

		cfg.Step("update strategy, attribute and value", func() {
			updatedStrategy, err := client.StrategyRegistry().UpdateStrategy(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.UpdateStrategyByRefRequest{
				StrategyRef: graph.StrategyID,
				Patch: &strategyregistryv1.UpdateStrategyRequest{
					Name: new("E2E Strategy v2"),
				},
			})
			require.NoError(tt, err)
			require.Equal(tt, "E2E Strategy v2", updatedStrategy.GetData().GetName())

			updatedAttribute, err := client.StrategyRegistry().UpdateAttribute(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.UpdateAttributeRequest{
				StrategyRef:  graph.StrategyID,
				AttributeRef: graph.AttributeID,
				Patch: &strategyregistryv1.UpdateAttributePatch{
					Name: new("E2E Attribute v2"),
				},
			})
			require.NoError(tt, err)
			require.Equal(tt, "E2E Attribute v2", updatedAttribute.GetData().GetName())

			const xssLikeValue = `<script>alert("strategy")</script>`
			updatedValue, err := client.StrategyRegistry().UpdateAttributeValue(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.UpdateAttributeValueRequest{
				StrategyRef:  graph.StrategyID,
				AttributeRef: graph.AttributeID,
				ValueRef:     graph.ValueID,
				Patch: &strategyregistryv1.UpdateAttributeValuePatch{
					Value:     new(xssLikeValue),
					Relations: []*strategyregistryv1.UpdateAttributeValueRelationInput{},
				},
			})
			require.NoError(tt, err)
			require.Equal(tt, xssLikeValue, updatedValue.GetData().GetValue())
		})

		cfg.Step("batch graph action updates the graph", func() {
			resp, err := client.StrategyRegistry().BatchActionStrategyGraph(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.BatchActionStrategyGraphRequest{
				StrategyRef: graph.StrategyID,
				Actions: []*strategyregistryv1.StrategyGraphAction{
					{
						Action: &strategyregistryv1.StrategyGraphAction_CreateAttribute{
							CreateAttribute: &strategyregistryv1.GraphActionCreateAttribute{
								Slug:        graph.AttributeSlug + "-batch",
								Name:        "Batch Attribute",
								Description: "created by strategy-registry grpc e2e",
							},
						},
					},
				},
			})
			require.NoError(tt, err)
			require.NotNil(tt, resp.GetData())
			require.NotEmpty(tt, resp.GetData().GetParameters())
		})

		cfg.Step("clone graph and reject invalid list query", func() {
			cloneResp, err := client.StrategyRegistry().CloneStrategies(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.CloneStrategiesRequest{
				Items: []*strategyregistryv1.CloneStrategyItemInput{{
					SourceStrategyId: graph.StrategyID,
					Slug:             graph.StrategySlug + "-clone",
				}},
			})
			require.NoError(tt, err)
			require.EqualValues(tt, 1, cloneResp.GetData().GetTotal())

			_, err = client.StrategyRegistry().ListStrategies(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.ListStrategiesRequest{
				PageSize: new(int32(1000)),
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.InvalidArgument, "strategy.pageOutOfRange")

			_, err = client.StrategyRegistry().ListStrategies(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.ListStrategiesRequest{
				Sort: new("unknown"),
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.InvalidArgument, "strategy.sortInvalid")
		})

		cfg.Step("delete value, attribute and strategy", func() {
			_, err := client.StrategyRegistry().DeleteAttributeValue(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.DeleteAttributeValueRequest{
				StrategyRef:  graph.StrategyID,
				AttributeRef: graph.AttributeID,
				ValueRef:     graph.ValueID,
			})
			require.NoError(tt, err)

			_, err = client.StrategyRegistry().DeleteAttribute(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.DeleteAttributeRequest{
				StrategyRef:  graph.StrategyID,
				AttributeRef: graph.AttributeID,
			})
			require.NoError(tt, err)

			_, err = client.StrategyRegistry().DeleteStrategy(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.DeleteStrategyRequest{
				StrategyRef: graph.StrategyID,
			})
			require.NoError(tt, err)

			_, err = client.StrategyRegistry().GetStrategy(grpcsupport.Context(cfg, actorToken), &strategyregistryv1.GetStrategyRequest{
				StrategyRef: graph.StrategyID,
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.NotFound, "strategy.notFound")
		})
	})
}
