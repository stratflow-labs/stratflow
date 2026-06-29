package testkit

import (
	"fmt"
	"net/http"
	"testing"

	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
)

type IdentityUser struct {
	ID       string
	Login    string
	Email    string
	Password string
	Role     string
	Token    string
}

func Login(t *testing.T, env Env, login, password string) string {
	t.Helper()

	resp := NewClient(env, "").IdentityPost(t, "/auth/login", map[string]any{
		"login":    login,
		"password": password,
	})
	RequireStatus(t, resp, http.StatusOK)
	return RequireStringField(t, RequireDataMap(t, resp), "accessToken")
}

func LoginAdmin(t *testing.T, env Env) Client {
	t.Helper()
	return NewClient(env, Login(t, env, env.AdminEmail, env.AdminPassword))
}

func CreateIdentityUser(t *testing.T, env Env, admin Client, role string) IdentityUser {
	t.Helper()

	user := e2ecommon.NewUserFixture("strategy-" + role)
	out := IdentityUser{
		Login:    user.Login,
		Email:    user.Email,
		Password: user.Password,
		Role:     role,
	}

	resp := admin.IdentityPost(t, "/users", map[string]any{
		"login":    out.Login,
		"email":    out.Email,
		"password": out.Password,
		"name":     fmt.Sprintf("Strategy %s", role),
		"lastName": "E2E",
		"gender":   1,
		"role":     role,
	})
	RequireStatus(t, resp, http.StatusCreated)
	out.ID = RequireDataString(t, resp, "id")
	out.Token = Login(t, env, out.Email, out.Password)

	t.Cleanup(func() {
		cleanup := admin.IdentityDelete(t, "/users/"+out.ID)
		if cleanup.StatusCode != http.StatusNoContent && cleanup.StatusCode != http.StatusNotFound {
			t.Fatalf("cleanup identity user %s: status=%d body=%s", out.ID, cleanup.StatusCode, cleanup.Body)
		}
	})

	return out
}
