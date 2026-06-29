package strategies_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/strategy-registry/tests/e2e/grpc/testkit"
	"google.golang.org/grpc/codes"
)

func TestGRPCStrategiesAccess(t *testing.T) {
	grpcsupport.Run(t, "strategy registry grpc protected methods require authorization metadata", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)

		cfg.Step("list strategies rejects missing authorization metadata", func() {
			_, err := client.StrategyRegistry().ListStrategies(grpcsupport.Context(cfg, ""), &strategyregistryv1.ListStrategiesRequest{})
			grpcsupport.RequireGRPCCode(tt, err, codes.Unauthenticated)
		})

		cfg.Step("create strategy rejects invalid bearer token", func() {
			_, err := client.StrategyRegistry().CreateStrategy(grpcsupport.Context(cfg, "invalid-token"), &strategyregistryv1.CreateStrategyRequest{
				Slug:        "invalid-token-create",
				Name:        "Invalid",
				Description: "Invalid",
			})
			grpcsupport.RequireGRPCCode(tt, err, codes.Unauthenticated)
		})
	})
}
