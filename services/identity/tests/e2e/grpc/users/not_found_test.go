package users_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	"github.com/google/uuid"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/grpc/testkit"
	"google.golang.org/grpc/codes"
)

func TestGRPCUsersByIDNotFound(t *testing.T) {
	grpcsupport.Run(t, "missing users return not found across get update delete", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		adminToken := grpcsupport.LoginAdmin(tt, cfg)
		missingID := uuid.NewString()

		cfg.Step("get missing user returns not found", func() {
			_, err := client.Identity().GetUser(grpcsupport.Context(cfg, adminToken), &identityv1.GetUserRequest{
				UserId: missingID,
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.NotFound, "user.notFound")
		})

		cfg.Step("update missing user returns not found", func() {
			_, err := client.Identity().UpdateUser(grpcsupport.Context(cfg, adminToken), &identityv1.UpdateUserByIDRequest{
				UserId: missingID,
				Patch: &identityv1.UpdateUserRequest{
					Name: new("Ghost"),
				},
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.NotFound, "user.notFound")
		})

		cfg.Step("delete missing user returns not found", func() {
			_, err := client.Identity().DeleteUser(grpcsupport.Context(cfg, adminToken), &identityv1.DeleteUserRequest{
				UserId: missingID,
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.NotFound, "user.notFound")
		})
	})
}
