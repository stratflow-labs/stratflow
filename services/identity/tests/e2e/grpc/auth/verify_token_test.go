package auth_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/grpc/testkit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestGRPCAuthVerifyToken(t *testing.T) {
	grpcsupport.Run(t, "verify token reports claims for valid tokens", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		token := loginAdminGRPCE2E(tt, cfg)

		cfg.Step("verify valid access token", func() {
			resp, err := client.Identity().VerifyToken(grpcsupport.Context(cfg, ""), &identityv1.VerifyTokenRequest{
				AccessToken: token,
			})
			require.NoError(tt, err)
			require.NotEmpty(tt, resp.GetUserId())
			require.Equal(tt, identityv1.Role_ROLE_ADMIN, resp.GetRole())
		})

		cfg.Step("verify invalid access token", func() {
			_, err := client.Identity().VerifyToken(grpcsupport.Context(cfg, ""), &identityv1.VerifyTokenRequest{
				AccessToken: "invalid-token",
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.Unauthenticated, "auth.accessTokenInvalid")
		})
	})
}
