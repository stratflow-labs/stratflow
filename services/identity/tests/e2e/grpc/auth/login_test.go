package auth_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/grpc/testkit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestGRPCAuthLogin(t *testing.T) {
	grpcsupport.Run(t, "admin can login through grpc api", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		var token string

		cfg.Step("login as seeded admin", func() {
			token = loginAdminGRPCE2E(tt, cfg)
		})

		cfg.Step("read current user with issued token", func() {
			resp, err := client.Identity().GetCurrentUser(grpcsupport.Context(cfg, token), &identityv1.GetCurrentUserRequest{})
			require.NoError(tt, err)
			require.Equal(tt, axiom.GetFixture[grpcsupport.Env](cfg, "env").AdminEmail, resp.GetData().GetEmail())
		})
	})
}

func TestGRPCAuthLoginInvalidCredentials(t *testing.T) {
	grpcsupport.Run(t, "invalid login credentials are rejected over grpc", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		env := axiom.GetFixture[grpcsupport.Env](cfg, "env")

		cfg.Step("try login with wrong password", func() {
			_, err := grpcsupport.ClientFromFixture(cfg).Identity().Login(grpcsupport.Context(cfg, ""), &identityv1.LoginRequest{
				Login:    env.AdminEmail,
				Password: "definitely-wrong-password",
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.Unauthenticated, "auth.invalidCredentials")
		})
	})
}
