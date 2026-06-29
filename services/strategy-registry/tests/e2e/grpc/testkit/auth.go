package testkit

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IdentityUser struct {
	ID       string
	Login    string
	Email    string
	Password string
	Role     string
	Token    string
}

func Login(t *testing.T, cfg *axiom.Config, login, password string) string {
	t.Helper()

	resp, err := ClientFromFixture(cfg).Identity().Login(Context(cfg, ""), &identityv1.LoginRequest{
		Login:    login,
		Password: password,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())
	require.NotEmpty(t, resp.GetData().GetAccessToken())
	return resp.GetData().GetAccessToken()
}

func LoginAdmin(t *testing.T, cfg *axiom.Config) string {
	t.Helper()
	env := axiom.GetFixture[Env](cfg, "env")
	return Login(t, cfg, env.AdminEmail, env.AdminPassword)
}

func CreateIdentityUser(t *testing.T, cfg *axiom.Config, adminToken, role string) IdentityUser {
	t.Helper()

	client := ClientFromFixture(cfg)
	fixture := e2ecommon.NewUserFixture("strategy-" + role)
	resp, err := client.Identity().CreateUser(Context(cfg, adminToken), &identityv1.CreateUserRequest{
		Login:    fixture.Login,
		Email:    fixture.Email,
		Password: fixture.Password,
		Name:     "Strategy " + role,
		LastName: "E2E",
		Gender:   new(int32(1)),
		Role:     role,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.GetData())

	user := IdentityUser{
		ID:       resp.GetData().GetId(),
		Login:    fixture.Login,
		Email:    fixture.Email,
		Password: fixture.Password,
		Role:     role,
		Token:    Login(t, cfg, fixture.Email, fixture.Password),
	}

	t.Cleanup(func() {
		_, cleanupErr := client.Identity().DeleteUser(Context(cfg, adminToken), &identityv1.DeleteUserRequest{UserId: user.ID})
		if cleanupErr != nil && status.Code(cleanupErr) != codes.NotFound {
			t.Fatalf("cleanup identity user %s: %v", user.ID, cleanupErr)
		}
	})

	return user
}
