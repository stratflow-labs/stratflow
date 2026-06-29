package testkit

import (
	"net/http"
	"testing"

	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
)

type IdentityE2EUser = e2ecommon.UserFixture

func CreateIdentityE2EUser(t *testing.T, admin Client, role string) IdentityE2EUser {
	t.Helper()

	user := e2ecommon.NewUserFixture(role)

	resp := admin.Post(t, "/users", map[string]any{
		"login":    user.Login,
		"email":    user.Email,
		"password": user.Password,
		"name":     "E2E " + role,
		"lastName": "User",
		"gender":   1,
		"role":     role,
	})
	RequireStatus(t, resp, http.StatusCreated)
	user.ID = RequireJSONContainsData(t, resp, "id")

	t.Cleanup(func() {
		cleanupResp := admin.Delete(t, UserPath(user.ID))
		if cleanupResp.StatusCode != http.StatusNoContent && cleanupResp.StatusCode != http.StatusNotFound {
			t.Fatalf("cleanup user %s: got %s, body = %s", user.ID, StatusName(cleanupResp.StatusCode), cleanupResp.Body)
		}
	})

	return user
}

func LoginCreatedUser(t *testing.T, env Env, user *IdentityE2EUser) IdentityE2EUser {
	t.Helper()
	user.Token = Login(t, env, user.Email, user.Password)
	return *user
}

func UniqueE2ESuffix(prefix string) string {
	return e2ecommon.UniqueSuffix(prefix)
}

func UniqueE2ELogin(prefix string) string {
	return e2ecommon.UniqueLogin(prefix)
}
