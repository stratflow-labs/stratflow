package auth_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/grpc/testkit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestGRPCAuthLogout(t *testing.T) {
	grpcsupport.Run(t, "logout invalidates the issued grpc token", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		token := loginAdminGRPCE2E(tt, cfg)

		cfg.Step("logout with the issued token", func() {
			_, err := client.Identity().Logout(grpcsupport.Context(cfg, token), &identityv1.LogoutRequest{})
			require.NoError(tt, err)
		})

		cfg.Step("reusing the same token is rejected", func() {
			_, err := client.Identity().GetCurrentUser(grpcsupport.Context(cfg, token), &identityv1.GetCurrentUserRequest{})
			grpcsupport.RequireGRPCCode(tt, err, codes.Unauthenticated)
		})
	})
}
