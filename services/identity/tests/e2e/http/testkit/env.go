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
	AdminEmail    string
	AdminPassword string
	HTTPClient    *http.Client
}

var identityE2ERunner = axiom.NewRunner(
	axiom.WithRunnerMeta(
		axiom.WithMetaEpic("identity"),
		axiom.WithMetaFeature("black-box e2e"),
		axiom.WithMetaLayer("e2e"),
	),
)

func LoadEnv(t *testing.T) Env {
	t.Helper()

	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("IDENTITY_E2E_BASE_URL")), "/")
	admin, ok := e2ecommon.LoadAdminCredentials()
	if baseURL == "" || !ok {
		t.Skip("identity e2e disabled: set IDENTITY_E2E_BASE_URL, IDENTITY_E2E_ADMIN_EMAIL and IDENTITY_E2E_ADMIN_PASSWORD for a test environment")
	}

	return Env{
		BaseURL:       baseURL,
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
			axiom.WithMetaSuite("identity e2e"),
			axiom.WithMetaTag("identity"),
		),
	)

	identityE2ERunner.RunCase(t, c, func(cfg *axiom.Config) {
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
