package users_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/grpc/testkit"
	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
	"github.com/stretchr/testify/require"
)

func TestGRPCUsersList(t *testing.T) {
	grpcsupport.Run(t, "list users supports search and bounded page size", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		adminToken := grpcsupport.LoginAdmin(tt, cfg)
		created := createIdentityGRPCE2EUser(tt, cfg, adminToken, "user")

		cfg.Step("search narrows results to the created user", func() {
			resp, err := client.Identity().ListUsers(grpcsupport.Context(cfg, adminToken), &identityv1.ListUsersRequest{
				Search: new(created.Email),
			})
			require.NoError(tt, err)
			require.NotNil(tt, resp.GetData())
			require.GreaterOrEqual(tt, len(resp.GetData().GetItems()), 1)

			found := false
			for _, item := range resp.GetData().GetItems() {
				if item.GetId() == created.ID {
					found = true
					break
				}
			}
			require.True(tt, found, "expected created user to be present in search results")
		})

		cfg.Step("oversized page size is clamped rather than failing", func() {
			resp, err := client.Identity().ListUsers(grpcsupport.Context(cfg, adminToken), &identityv1.ListUsersRequest{
				PageSize: new(int32(1000)),
			})
			require.NoError(tt, err)
			require.NotNil(tt, resp.GetData())
			require.LessOrEqual(tt, len(resp.GetData().GetItems()), 100)
		})

		cfg.Step("pagination splits a filtered result set across pages", func() {
			marker := "pagination-" + e2ecommon.UniqueSuffix("grpc")
			first := createIdentityGRPCE2EUser(tt, cfg, adminToken, "user")
			second := createIdentityGRPCE2EUser(tt, cfg, adminToken, "user")

			_, err := client.Identity().UpdateUser(grpcsupport.Context(cfg, adminToken), &identityv1.UpdateUserByIDRequest{
				UserId: first.ID,
				Patch: &identityv1.UpdateUserRequest{
					Name: new(marker + "-first"),
				},
			})
			require.NoError(tt, err)

			_, err = client.Identity().UpdateUser(grpcsupport.Context(cfg, adminToken), &identityv1.UpdateUserByIDRequest{
				UserId: second.ID,
				Patch: &identityv1.UpdateUserRequest{
					Name: new(marker + "-second"),
				},
			})
			require.NoError(tt, err)

			firstPage, err := client.Identity().ListUsers(grpcsupport.Context(cfg, adminToken), &identityv1.ListUsersRequest{
				Search:   new(marker),
				Page:     new(int32(1)),
				PageSize: new(int32(1)),
				Sort:     new("created_at ASC"),
			})
			require.NoError(tt, err)
			require.NotNil(tt, firstPage.GetData())
			require.Len(tt, firstPage.GetData().GetItems(), 1)
			require.GreaterOrEqual(tt, firstPage.GetData().GetTotal(), int64(2))

			secondPage, err := client.Identity().ListUsers(grpcsupport.Context(cfg, adminToken), &identityv1.ListUsersRequest{
				Search:   new(marker),
				Page:     new(int32(2)),
				PageSize: new(int32(1)),
				Sort:     new("created_at ASC"),
			})
			require.NoError(tt, err)
			require.NotNil(tt, secondPage.GetData())
			require.Len(tt, secondPage.GetData().GetItems(), 1)
			require.GreaterOrEqual(tt, secondPage.GetData().GetTotal(), int64(2))
			require.NotEqual(tt, firstPage.GetData().GetItems()[0].GetId(), secondPage.GetData().GetItems()[0].GetId())
		})
	})
}
