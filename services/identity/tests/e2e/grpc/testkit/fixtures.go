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
		return nil, nil, errors.New("identity grpc e2e env is not configured")
	}
	return env, nil, nil
}

func ConnFixture(cfg *axiom.Config) (any, func(), error) {
	env := axiom.GetFixture[Env](cfg, "env")

	conn, err := grpc.NewClient(
		env.Target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = conn.Close()
	}

	return conn, cleanup, nil
}
