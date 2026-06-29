package users_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/grpc/testkit"
	"google.golang.org/grpc/codes"
)

func TestGRPCUsersValidation(t *testing.T) {
	grpcsupport.Run(t, "grpc user validation returns stable domain reasons", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		adminToken := grpcsupport.LoginAdmin(tt, cfg)

		cfg.Step("create user rejects invalid email", func() {
			_, err := client.Identity().CreateUser(grpcsupport.Context(cfg, adminToken), &identityv1.CreateUserRequest{
				Login:    uniqueGRPCE2ELogin("invalid"),
				Email:    "not-an-email",
				Password: "E2e-password-12345",
				Name:     "Invalid",
				LastName: "Email",
				Role:     "user",
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.InvalidArgument, "auth.emailInvalid")
		})

		cfg.Step("get user rejects malformed identifiers as not found", func() {
			_, err := client.Identity().GetUser(grpcsupport.Context(cfg, adminToken), &identityv1.GetUserRequest{
				UserId: "not-a-uuid",
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.NotFound, "user.notFound")
		})
	})
}
