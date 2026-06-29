package users_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/grpc/testkit"
	"github.com/stretchr/testify/require"
)

func TestGRPCUsersUpdate(t *testing.T) {
	grpcsupport.Run(t, "partial update preserves untouched user fields", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		adminToken := grpcsupport.LoginAdmin(tt, cfg)
		created := createIdentityGRPCE2EUser(tt, cfg, adminToken, "user")

		var original *identityv1.User

		cfg.Step("load original user state", func() {
			resp, err := client.Identity().GetUser(grpcsupport.Context(cfg, adminToken), &identityv1.GetUserRequest{
				UserId: created.ID,
			})
			require.NoError(tt, err)
			original = resp.GetData()
			require.NotNil(tt, original)
		})

		cfg.Step("update only the display name", func() {
			resp, err := client.Identity().UpdateUser(grpcsupport.Context(cfg, adminToken), &identityv1.UpdateUserByIDRequest{
				UserId: created.ID,
				Patch: &identityv1.UpdateUserRequest{
					Name: new("Patched Name"),
				},
			})
			require.NoError(tt, err)
			require.Equal(tt, "Patched Name", resp.GetData().GetName())
			require.Equal(tt, original.GetLastName(), resp.GetData().GetLastName())
			require.Equal(tt, original.GetEmail(), resp.GetData().GetEmail())
			require.Equal(tt, original.GetGender(), resp.GetData().GetGender())
		})
	})
}
