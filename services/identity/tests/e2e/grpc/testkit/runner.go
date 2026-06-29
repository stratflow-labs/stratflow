package testkit

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
)

var identityGRPCE2ERunner = axiom.NewRunner(
	axiom.WithRunnerMeta(
		axiom.WithMetaEpic("identity"),
		axiom.WithMetaFeature("grpc e2e"),
		axiom.WithMetaLayer("e2e"),
	),
	axiom.WithRunnerFixture("env", EnvFixture),
	axiom.WithRunnerFixture("grpc_conn", ConnFixture),
)

func Run(t *testing.T, name string, fn func(*axiom.Config)) {
	t.Helper()

	if _, ok := LoadEnv(); !ok {
		t.Skip("identity grpc e2e disabled: set IDENTITY_GRPC_URL, IDENTITY_E2E_ADMIN_EMAIL and IDENTITY_E2E_ADMIN_PASSWORD")
	}

	c := axiom.NewCase(
		axiom.WithCaseName(name),
		axiom.WithCaseMeta(
			axiom.WithMetaSuite("identity grpc e2e"),
			axiom.WithMetaTag("identity"),
			axiom.WithMetaTag("grpc"),
		),
	)

	identityGRPCE2ERunner.RunCase(t, c, fn)
}

func CaseT(t *testing.T, cfg *axiom.Config) *testing.T {
	t.Helper()
	if cfg != nil && cfg.SubT != nil {
		return cfg.SubT
	}
	return t
}
