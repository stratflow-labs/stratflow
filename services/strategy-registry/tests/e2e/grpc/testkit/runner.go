package testkit

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
)

var strategyRegistryGRPCE2ERunner = axiom.NewRunner(
	axiom.WithRunnerMeta(
		axiom.WithMetaEpic("strategy-registry"),
		axiom.WithMetaFeature("grpc e2e"),
		axiom.WithMetaLayer("e2e"),
	),
	axiom.WithRunnerFixture("env", EnvFixture),
	axiom.WithRunnerFixture("strategy_grpc_conn", StrategyConnFixture),
	axiom.WithRunnerFixture("identity_grpc_conn", IdentityConnFixture),
)

func Run(t *testing.T, name string, fn func(*axiom.Config)) {
	t.Helper()

	if _, ok := LoadEnv(); !ok {
		t.Skip("strategy-registry grpc e2e disabled: set STRATEGY_REGISTRY_GRPC_URL, IDENTITY_GRPC_URL, IDENTITY_E2E_ADMIN_EMAIL and IDENTITY_E2E_ADMIN_PASSWORD")
	}

	c := axiom.NewCase(
		axiom.WithCaseName(name),
		axiom.WithCaseMeta(
			axiom.WithMetaSuite("strategy-registry grpc e2e"),
			axiom.WithMetaTag("strategy-registry"),
			axiom.WithMetaTag("grpc"),
		),
	)

	strategyRegistryGRPCE2ERunner.RunCase(t, c, fn)
}

func CaseT(t *testing.T, cfg *axiom.Config) *testing.T {
	t.Helper()
	if cfg != nil && cfg.SubT != nil {
		return cfg.SubT
	}
	return t
}
