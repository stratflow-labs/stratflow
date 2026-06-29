package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

type e2eIdentityUser struct {
	ID       string
	Login    string
	Email    string
	Password string
	Role     string
	Token    string
}

func loginIdentity(t *testing.T, env strategyE2EEnv, email, password string) string {
	t.Helper()
	anonymous := newStrategyE2EClient(env, "")
	resp := anonymous.identityPost(t, env, "/auth/login", map[string]any{
		"login":    email,
		"password": password,
	})
	requireStatus(t, resp, http.StatusOK)
	data := requireData(t, resp)
	token, ok := data["accessToken"].(string)
	if !ok || token == "" {
		t.Fatalf("expected accessToken in %v", data)
	}
	return token
}

func adminIdentityClient(t *testing.T, env strategyE2EEnv) strategyE2EClient {
	t.Helper()
	return newStrategyE2EClient(env, loginIdentity(t, env, env.AdminEmail, env.AdminPassword))
}

func createIdentityUser(t *testing.T, env strategyE2EEnv, admin strategyE2EClient, role string) e2eIdentityUser {
	t.Helper()
	suffix := uuid.NewString()
	user := e2eIdentityUser{
		Login:    fmt.Sprintf("%.10s-%s", "strategy-"+role, suffix[:12]),
		Email:    fmt.Sprintf("strategy.%s.%s@example.test", role, suffix),
		Password: "E2e-password-12345",
		Role:     role,
	}
	resp := admin.identityPost(t, env, "/users", map[string]any{
		"login":    user.Login,
		"email":    user.Email,
		"password": user.Password,
		"name":     "Strategy " + role,
		"lastName": "E2E",
		"gender":   1,
		"role":     role,
	})
	requireStatus(t, resp, http.StatusCreated)
	user.ID = requireDataString(t, resp, "id")
	user.Token = loginIdentity(t, env, user.Email, user.Password)
	t.Cleanup(func() {
		cleanup := admin.identityDelete(t, env, "/users/"+user.ID)
		if cleanup.StatusCode != http.StatusNoContent && cleanup.StatusCode != http.StatusNotFound {
			t.Fatalf("cleanup identity user %s: status=%d body=%s", user.ID, cleanup.StatusCode, cleanup.Body)
		}
	})
	return user
}
