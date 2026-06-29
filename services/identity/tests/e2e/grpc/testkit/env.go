package testkit

import (
	"os"
	"strings"
	"time"

	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
)

type Env struct {
	Target        string
	AdminEmail    string
	AdminPassword string
	Timeout       time.Duration
}

func LoadEnv() (Env, bool) {
	target := strings.TrimSpace(os.Getenv("IDENTITY_GRPC_URL"))
	admin, ok := e2ecommon.LoadAdminCredentials()
	if target == "" || !ok {
		return Env{}, false
	}

	return Env{
		Target:        target,
		AdminEmail:    admin.Email,
		AdminPassword: admin.Password,
		Timeout:       e2ecommon.DefaultTimeout,
	}, true
}
