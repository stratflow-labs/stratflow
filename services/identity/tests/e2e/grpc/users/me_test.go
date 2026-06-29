package users_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/grpc/testkit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestGRPCUsersMe(t *testing.T) {
	grpcsupport.Run(t, "current user lifecycle works through grpc api", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		adminToken := grpcsupport.LoginAdmin(tt, cfg)
		created := createIdentityGRPCE2EUser(tt, cfg, adminToken, "user")
		user := loginCreatedUserGRPCE2E(tt, cfg, &created)
		client := grpcsupport.ClientFromFixture(cfg)

		cfg.Step("get current user resolves subject from access token", func() {
			resp, err := client.Identity().GetCurrentUser(grpcsupport.Context(cfg, user.Token), &identityv1.GetCurrentUserRequest{})
			require.NoError(tt, err)
			require.Equal(tt, user.ID, resp.GetData().GetId())
		})

		cfg.Step("update current user applies self-service profile changes", func() {
			resp, err := client.Identity().UpdateCurrentUser(grpcsupport.Context(cfg, user.Token), &identityv1.UpdateCurrentUserRequest{
				Patch: &identityv1.UpdateUserRequest{
					Name:     new("Self Updated"),
					LastName: new("Profile"),
					Gender:   new(int32(2)),
				},
			})
			require.NoError(tt, err)
			require.Equal(tt, "Self Updated", resp.GetData().GetName())
			require.Equal(tt, "Profile", resp.GetData().GetLastName())
		})

		cfg.Step("delete current user revokes further resolution of the token subject", func() {
			_, err := client.Identity().DeleteCurrentUser(grpcsupport.Context(cfg, user.Token), &identityv1.DeleteCurrentUserRequest{})
			require.NoError(tt, err)

			_, err = client.Identity().GetCurrentUser(grpcsupport.Context(cfg, user.Token), &identityv1.GetCurrentUserRequest{})
			grpcsupport.RequireGRPCCode(tt, err, codes.NotFound)
		})
	})
}

func TestGRPCUsersMeMissingOrInvalidToken(t *testing.T) {
	grpcsupport.Run(t, "current user methods reject missing or forged tokens", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)

		cfg.Step("missing authorization metadata is rejected", func() {
			_, err := client.Identity().GetCurrentUser(grpcsupport.Context(cfg, ""), &identityv1.GetCurrentUserRequest{})
			grpcsupport.RequireGRPCCode(tt, err, codes.Unauthenticated)
		})

		cfg.Step("forged authorization metadata is rejected", func() {
			_, err := client.Identity().GetCurrentUser(grpcsupport.Context(cfg, "invalid-token"), &identityv1.GetCurrentUserRequest{})
			grpcsupport.RequireGRPCCode(tt, err, codes.Unauthenticated)
		})
	})
}
