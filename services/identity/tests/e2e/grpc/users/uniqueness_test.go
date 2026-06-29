package users_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/grpc/testkit"
	"google.golang.org/grpc/codes"
)

func TestGRPCUsersUniqueness(t *testing.T) {
	grpcsupport.Run(t, "create user enforces unique login and email", func(cfg *axiom.Config) {
		tt := grpcsupport.CaseT(t, cfg)
		client := grpcsupport.ClientFromFixture(cfg)
		adminToken := grpcsupport.LoginAdmin(tt, cfg)
		created := createIdentityGRPCE2EUser(tt, cfg, adminToken, "user")

		cfg.Step("duplicate login is rejected", func() {
			_, err := client.Identity().CreateUser(grpcsupport.Context(cfg, adminToken), &identityv1.CreateUserRequest{
				Login:    created.Login,
				Email:    uniqueGRPCE2ELogin("email") + "@example.test",
				Password: "E2e-password-12345",
				Name:     "Duplicate",
				LastName: "Login",
				Role:     "user",
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.AlreadyExists, "auth.loginAlreadyUsed")
		})

		cfg.Step("duplicate email is rejected", func() {
			_, err := client.Identity().CreateUser(grpcsupport.Context(cfg, adminToken), &identityv1.CreateUserRequest{
				Login:    uniqueGRPCE2ELogin("login"),
				Email:    created.Email,
				Password: "E2e-password-12345",
				Name:     "Duplicate",
				LastName: "Email",
				Role:     "user",
			})
			grpcsupport.RequireGRPCReason(tt, err, codes.AlreadyExists, "auth.emailAlreadyUsed")
		})
	})
}
