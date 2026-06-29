package tests

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Nikita-Filonov/axiom"
)

type strategyE2EEnv struct {
	BaseURL       string
	IdentityURL   string
	AdminEmail    string
	AdminPassword string
	HTTPClient    *http.Client
}

var strategyE2ERunner = axiom.NewRunner(
	axiom.WithRunnerMeta(
		axiom.WithMetaEpic("strategy-registry"),
		axiom.WithMetaFeature("black-box e2e"),
		axiom.WithMetaLayer("e2e"),
	),
)

func loadStrategyE2EEnv(t *testing.T) strategyE2EEnv {
	t.Helper()

	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("STRATEGY_REGISTRY_E2E_BASE_URL")), "/")
	identityURL := strings.TrimRight(strings.TrimSpace(os.Getenv("IDENTITY_E2E_BASE_URL")), "/")
	adminEmail := strings.TrimSpace(os.Getenv("IDENTITY_E2E_ADMIN_EMAIL"))
	adminPassword := strings.TrimSpace(os.Getenv("IDENTITY_E2E_ADMIN_PASSWORD"))
	if baseURL == "" || identityURL == "" || adminEmail == "" || adminPassword == "" {
		t.Skip("strategy-registry e2e disabled: set STRATEGY_REGISTRY_E2E_BASE_URL, IDENTITY_E2E_BASE_URL, IDENTITY_E2E_ADMIN_EMAIL and IDENTITY_E2E_ADMIN_PASSWORD")
	}

	return strategyE2EEnv{
		BaseURL:       baseURL,
		IdentityURL:   identityURL,
		AdminEmail:    adminEmail,
		AdminPassword: adminPassword,
		HTTPClient:    &http.Client{Timeout: 10 * time.Second},
	}
}

func runStrategyE2E(t *testing.T, name string, fn func(*axiom.Config, strategyE2EEnv)) {
	t.Helper()

	env := loadStrategyE2EEnv(t)
	c := axiom.NewCase(
		axiom.WithCaseName(name),
		axiom.WithCaseMeta(
			axiom.WithMetaSuite("strategy-registry e2e"),
			axiom.WithMetaTag("strategy-registry"),
		),
	)
	strategyE2ERunner.RunCase(t, c, func(cfg *axiom.Config) {
		fn(cfg, env)
	})
}

func e2ET(t *testing.T, cfg *axiom.Config) *testing.T {
	t.Helper()
	if cfg != nil && cfg.SubT != nil {
		return cfg.SubT
	}
	return t
}
