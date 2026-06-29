package testkit

import (
	"context"
	"strings"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	strategy strategyregistryv1.StrategyRegistryServiceClient
	identity identityv1.IdentityServiceClient
}

func ClientFromFixture(cfg *axiom.Config) Client {
	strategyConn := axiom.GetFixture[*grpc.ClientConn](cfg, "strategy_grpc_conn")
	identityConn := axiom.GetFixture[*grpc.ClientConn](cfg, "identity_grpc_conn")
	return Client{
		strategy: strategyregistryv1.NewStrategyRegistryServiceClient(strategyConn),
		identity: identityv1.NewIdentityServiceClient(identityConn),
	}
}

func Context(cfg *axiom.Config, token string) context.Context {
	_ = axiom.GetFixture[Env](cfg, "env")
	ctx := context.Background()
	if strings.TrimSpace(token) == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+strings.TrimSpace(token))
}

func (c Client) StrategyRegistry() strategyregistryv1.StrategyRegistryServiceClient {
	return c.strategy
}

func (c Client) Identity() identityv1.IdentityServiceClient {
	return c.identity
}
