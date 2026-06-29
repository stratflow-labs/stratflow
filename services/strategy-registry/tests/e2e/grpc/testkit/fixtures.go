package testkit

import (
	"errors"

	"github.com/Nikita-Filonov/axiom"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func EnvFixture(_ *axiom.Config) (any, func(), error) {
	env, ok := LoadEnv()
	if !ok {
		return nil, nil, errors.New("strategy-registry grpc e2e env is not configured")
	}
	return env, nil, nil
}

func StrategyConnFixture(cfg *axiom.Config) (any, func(), error) {
	env := axiom.GetFixture[Env](cfg, "env")
	return newConn(env.StrategyTarget)
}

func IdentityConnFixture(cfg *axiom.Config) (any, func(), error) {
	env := axiom.GetFixture[Env](cfg, "env")
	return newConn(env.IdentityTarget)
}

func newConn(target string) (*grpc.ClientConn, func(), error) {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	return conn, func() { _ = conn.Close() }, nil
}
