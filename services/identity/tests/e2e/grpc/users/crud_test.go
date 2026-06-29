package users_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/grpc/testkit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestGRPCUsersCRUD(t *testing.T) {
	grpcsupport.Run(t, "admin can execute full user crud over grpc", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		adminToken := grpcsupport.LoginAdmin(tt, cfg)
		created := createIdentityGRPCE2EUser(tt, cfg, adminToken, "user")

		cfg.Step("list users includes the created user", func() {
			resp, err := client.Identity().ListUsers(grpcsupport.Context(cfg, adminToken), &identityv1.ListUsersRequest{
				Search: new(created.Email),
			})
			require.NoError(tt, err)
			require.NotNil(tt, resp.GetData())
			require.NotEmpty(tt, resp.GetData().GetItems())
		})

		cfg.Step("get user by id returns created profile", func() {
			resp, err := client.Identity().GetUser(grpcsupport.Context(cfg, adminToken), &identityv1.GetUserRequest{
				UserId: created.ID,
			})
			require.NoError(tt, err)
			require.Equal(tt, created.Email, resp.GetData().GetEmail())
			require.Equal(tt, created.Login, resp.GetData().GetLogin())
		})

		cfg.Step("update user by id persists profile changes", func() {
			resp, err := client.Identity().UpdateUser(grpcsupport.Context(cfg, adminToken), &identityv1.UpdateUserByIDRequest{
				UserId: created.ID,
				Patch: &identityv1.UpdateUserRequest{
					Name:     new("Updated"),
					LastName: new("Trader"),
					Gender:   new(int32(2)),
				},
			})
			require.NoError(tt, err)
			require.Equal(tt, "Updated", resp.GetData().GetName())
			require.Equal(tt, "Trader", resp.GetData().GetLastName())
			require.EqualValues(tt, 2, resp.GetData().GetGender())
		})

		cfg.Step("delete user removes it from subsequent reads", func() {
			_, err := client.Identity().DeleteUser(grpcsupport.Context(cfg, adminToken), &identityv1.DeleteUserRequest{
				UserId: created.ID,
			})
			require.NoError(tt, err)

			_, err = client.Identity().GetUser(grpcsupport.Context(cfg, adminToken), &identityv1.GetUserRequest{
				UserId: created.ID,
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.NotFound, "user.notFound")
		})
	})
}

func TestGRPCUsersCRUDMissingToken(t *testing.T) {
	grpcsupport.Run(t, "protected user methods reject missing authorization metadata", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)

		cfg.Step("list users requires authorization metadata", func() {
			_, err := client.Identity().ListUsers(grpcsupport.Context(cfg, ""), &identityv1.ListUsersRequest{})
			grpcsupport.RequireGRPCCode(tt, err, codes.Unauthenticated)
		})

		cfg.Step("create user requires authorization metadata", func() {
			_, err := client.Identity().CreateUser(grpcsupport.Context(cfg, ""), &identityv1.CreateUserRequest{
				Login:    "anon-user",
				Email:    "anon@example.test",
				Password: "E2e-password-12345",
				Name:     "Anon",
				LastName: "User",
				Role:     "user",
			})
			grpcsupport.RequireGRPCCode(tt, err, codes.Unauthenticated)
		})
	})
}
