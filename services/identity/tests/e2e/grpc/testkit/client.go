package testkit

import (
	"context"
	"strings"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	client identityv1.IdentityServiceClient
}

func ClientFromFixture(cfg *axiom.Config) Client {
	conn := axiom.GetFixture[*grpc.ClientConn](cfg, "grpc_conn")
	return Client{
		client: identityv1.NewIdentityServiceClient(conn),
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

func (c Client) Identity() identityv1.IdentityServiceClient {
	return c.client
}
