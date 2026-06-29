package users_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/grpc/testkit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestGRPCUsersDelete(t *testing.T) {
	grpcsupport.Run(t, "delete user removes it from reads and filtered lists", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		adminToken := grpcsupport.LoginAdmin(tt, cfg)
		created := createIdentityGRPCE2EUser(tt, cfg, adminToken, "user")
		manager := loginCreatedUserGRPCE2E(tt, cfg, new(createIdentityGRPCE2EUser(tt, cfg, adminToken, "manager")))

		cfg.Step("delete user rejects missing and forged authorization metadata", func() {
			_, err := client.Identity().DeleteUser(grpcsupport.Context(cfg, ""), &identityv1.DeleteUserRequest{
				UserId: created.ID,
			})
			grpcsupport.RequireGRPCCode(tt, err, codes.Unauthenticated)

			_, err = client.Identity().DeleteUser(grpcsupport.Context(cfg, "invalid-token"), &identityv1.DeleteUserRequest{
				UserId: created.ID,
			})
			grpcsupport.RequireGRPCCode(tt, err, codes.Unauthenticated)
		})

		cfg.Step("delete user rejects non-admin foreign deletion", func() {
			_, err := client.Identity().DeleteUser(grpcsupport.Context(cfg, manager.Token), &identityv1.DeleteUserRequest{
				UserId: created.ID,
			})
			grpcsupport.RequireGRPCCode(tt, err, codes.PermissionDenied)
		})

		cfg.Step("delete user by id succeeds", func() {
			_, err := client.Identity().DeleteUser(grpcsupport.Context(cfg, adminToken), &identityv1.DeleteUserRequest{
				UserId: created.ID,
			})
			require.NoError(tt, err)
		})

		cfg.Step("subsequent get returns not found", func() {
			_, err := client.Identity().GetUser(grpcsupport.Context(cfg, adminToken), &identityv1.GetUserRequest{
				UserId: created.ID,
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.NotFound, "user.notFound")
		})

		cfg.Step("filtered list no longer includes the deleted user", func() {
			resp, err := client.Identity().ListUsers(grpcsupport.Context(cfg, adminToken), &identityv1.ListUsersRequest{
				Search: new(created.Email),
			})
			require.NoError(tt, err)
			require.NotNil(tt, resp.GetData())
			require.Empty(tt, resp.GetData().GetItems())
		})
	})
}
