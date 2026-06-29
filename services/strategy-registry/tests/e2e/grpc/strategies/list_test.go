package strategies_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/strategy-registry/tests/e2e/grpc/testkit"
	"github.com/stretchr/testify/require"
)

func TestGRPCStrategiesList(t *testing.T) {
	grpcsupport.Run(t, "list strategies supports search and pagination over grpc", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		adminToken := grpcsupport.LoginAdmin(tt, cfg)
		graph := grpcsupport.CreateStrategyGraph(tt, cfg, adminToken)

		cfg.Step("search narrows results to the created strategy", func() {
			resp, err := client.StrategyRegistry().ListStrategies(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.ListStrategiesRequest{
				Search: new(graph.StrategySlug),
			})
			require.NoError(tt, err)
			require.NotNil(tt, resp.GetData())
			require.GreaterOrEqual(tt, len(resp.GetData().GetItems()), 1)
		})

		cfg.Step("pagination splits a filtered result set across pages", func() {
			marker := "pagination-" + e2ecommon.UniqueSuffix("strategy-grpc")

			firstResp, err := client.StrategyRegistry().CreateStrategy(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.CreateStrategyRequest{
				Slug:        marker + "-first",
				Name:        marker + "-first",
				Description: "pagination marker",
			})
			require.NoError(tt, err)

			secondResp, err := client.StrategyRegistry().CreateStrategy(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.CreateStrategyRequest{
				Slug:        marker + "-second",
				Name:        marker + "-second",
				Description: "pagination marker",
			})
			require.NoError(tt, err)

			t.Cleanup(func() {
				_, _ = client.StrategyRegistry().DeleteStrategy(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.DeleteStrategyRequest{StrategyRef: firstResp.GetData().GetId()})
			})
			t.Cleanup(func() {
				_, _ = client.StrategyRegistry().DeleteStrategy(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.DeleteStrategyRequest{StrategyRef: secondResp.GetData().GetId()})
			})

			firstPage, err := client.StrategyRegistry().ListStrategies(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.ListStrategiesRequest{
				Search:   new(marker),
				Page:     new(int32(1)),
				PageSize: new(int32(1)),
				Sort:     new("created_at_asc"),
			})
			require.NoError(tt, err)
			require.NotNil(tt, firstPage.GetData())
			require.Len(tt, firstPage.GetData().GetItems(), 1)
			require.GreaterOrEqual(tt, firstPage.GetData().GetTotal(), int64(2))

			secondPage, err := client.StrategyRegistry().ListStrategies(grpcsupport.Context(cfg, adminToken), &strategyregistryv1.ListStrategiesRequest{
				Search:   new(marker),
				Page:     new(int32(2)),
				PageSize: new(int32(1)),
				Sort:     new("created_at_asc"),
			})
			require.NoError(tt, err)
			require.NotNil(tt, secondPage.GetData())
			require.Len(tt, secondPage.GetData().GetItems(), 1)
			require.GreaterOrEqual(tt, secondPage.GetData().GetTotal(), int64(2))
			require.NotEqual(tt, firstPage.GetData().GetItems()[0].GetId(), secondPage.GetData().GetItems()[0].GetId())
		})
	})
}
