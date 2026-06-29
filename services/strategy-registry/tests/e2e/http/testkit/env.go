package testkit

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
)

type Env struct {
	BaseURL       string
	IdentityURL   string
	AdminEmail    string
	AdminPassword string
	HTTPClient    *http.Client
}

var strategyRegistryE2ERunner = axiom.NewRunner(
	axiom.WithRunnerMeta(
		axiom.WithMetaEpic("strategy-registry"),
		axiom.WithMetaFeature("black-box e2e"),
		axiom.WithMetaLayer("e2e"),
	),
)

func LoadEnv(t *testing.T) Env {
	t.Helper()

	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("STRATEGY_REGISTRY_E2E_BASE_URL")), "/")
	identityURL := strings.TrimRight(strings.TrimSpace(os.Getenv("IDENTITY_E2E_BASE_URL")), "/")
	admin, ok := e2ecommon.LoadAdminCredentials()
	if baseURL == "" || identityURL == "" || !ok {
		t.Skip("strategy-registry e2e disabled: set STRATEGY_REGISTRY_E2E_BASE_URL, IDENTITY_E2E_BASE_URL, IDENTITY_E2E_ADMIN_EMAIL and IDENTITY_E2E_ADMIN_PASSWORD")
	}

	return Env{
		BaseURL:       baseURL,
		IdentityURL:   identityURL,
		AdminEmail:    admin.Email,
		AdminPassword: admin.Password,
		HTTPClient: &http.Client{
			Timeout: e2ecommon.DefaultTimeout,
		},
	}
}

func Run(t *testing.T, name string, fn func(*axiom.Config, Env)) {
	t.Helper()

	env := LoadEnv(t)
	c := axiom.NewCase(
		axiom.WithCaseName(name),
		axiom.WithCaseMeta(
			axiom.WithMetaSuite("strategy-registry e2e"),
			axiom.WithMetaTag("strategy-registry"),
		),
	)

	strategyRegistryE2ERunner.RunCase(t, c, func(cfg *axiom.Config) {
		fn(cfg, env)
	})
}

func CaseT(t *testing.T, cfg *axiom.Config) *testing.T {
	t.Helper()
	if cfg != nil && cfg.SubT != nil {
		return cfg.SubT
	}
	return t
}
