package testkit

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	"github.com/stretchr/testify/require"
)

func Login(t *testing.T, cfg *axiom.Config, login, password string) string {
	t.Helper()

	resp, err := ClientFromFixture(cfg).Identity().Login(Context(cfg, ""), &identityv1.LoginRequest{
		Login:    login,
		Password: password,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	return RequireTokenEnvelope(t, &IdentityTokenEnvelopeView{
		AccessToken: resp.GetData().GetAccessToken(),
	})
}

func LoginAdmin(t *testing.T, cfg *axiom.Config) string {
	t.Helper()
	env := axiom.GetFixture[Env](cfg, "env")
	return Login(t, cfg, env.AdminEmail, env.AdminPassword)
}
