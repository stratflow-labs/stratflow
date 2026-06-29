package auth_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	grpcsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/grpc/testkit"
)

func loginAdminGRPCE2E(t *testing.T, cfg *axiom.Config) string {
	t.Helper()
	return grpcsupport.LoginAdmin(t, cfg)
}
