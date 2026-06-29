package testkit

import (
	"os"
	"strings"
	"time"

	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
)

type Env struct {
	StrategyTarget string
	IdentityTarget string
	AdminEmail     string
	AdminPassword  string
	Timeout        time.Duration
}

func LoadEnv() (Env, bool) {
	strategyTarget := strings.TrimSpace(os.Getenv("STRATEGY_REGISTRY_GRPC_URL"))
	identityTarget := strings.TrimSpace(os.Getenv("IDENTITY_GRPC_URL"))
	admin, ok := e2ecommon.LoadAdminCredentials()
	if strategyTarget == "" || identityTarget == "" || !ok {
		return Env{}, false
	}

	return Env{
		StrategyTarget: strategyTarget,
		IdentityTarget: identityTarget,
		AdminEmail:     admin.Email,
		AdminPassword:  admin.Password,
		Timeout:        e2ecommon.DefaultTimeout,
	}, true
}
