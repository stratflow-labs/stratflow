package testkit

import (
	"net/http"
	"testing"
)

func Login(t *testing.T, env Env, login, password string) string {
	t.Helper()

	resp := NewClient(env, "").Post(t, "/auth/login", map[string]any{
		"login":    login,
		"password": password,
	})
	RequireStatus(t, resp, http.StatusOK)

	data := RequireDataMap(t, resp)
	token := RequireStringField(t, data, "accessToken")
	return token
}

func LoginAdmin(t *testing.T, env Env) Client {
	t.Helper()
	token := Login(t, env, env.AdminEmail, env.AdminPassword)
	return NewClient(env, token)
}
