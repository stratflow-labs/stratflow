package users_test

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	grpcsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/grpc/testkit"
	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

type identityGRPCE2EUser = e2ecommon.UserFixture

func createIdentityGRPCE2EUser(t *testing.T, cfg *axiom.Config, adminToken, role string) identityGRPCE2EUser {
	t.Helper()

	user := e2ecommon.NewUserFixture(role)

	resp, err := grpcsupport.ClientFromFixture(cfg).Identity().CreateUser(grpcsupport.Context(cfg, adminToken), &identityv1.CreateUserRequest{
		Login:    user.Login,
		Email:    user.Email,
		Password: user.Password,
		Name:     "E2E " + role,
		LastName: "User",
		Gender:   new(int32(identityv1.Gender_GENDER_MALE)),
		Role:     role,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetData())

	user.ID = resp.GetData().GetId()
	require.NotEmpty(t, user.ID)

	t.Cleanup(func() {
		_, cleanupErr := grpcsupport.ClientFromFixture(cfg).Identity().DeleteUser(grpcsupport.Context(cfg, adminToken), &identityv1.DeleteUserRequest{
			UserId: user.ID,
		})
		if cleanupErr == nil {
			return
		}

		st := grpcsupport.RequireGRPCCode(t, cleanupErr, codes.NotFound)
		require.Equal(t, codes.NotFound, st.Code())
	})

	return user
}

func loginCreatedUserGRPCE2E(t *testing.T, cfg *axiom.Config, user *identityGRPCE2EUser) identityGRPCE2EUser {
	t.Helper()
	user.Token = grpcsupport.Login(t, cfg, user.Email, user.Password)
	return *user
}

func uniqueGRPCE2ELogin(prefix string) string {
	return e2ecommon.UniqueLogin(prefix)
}
