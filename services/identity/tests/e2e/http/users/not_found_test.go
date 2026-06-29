package tests

import (
	"net/http"
	"testing"

	"github.com/Nikita-Filonov/axiom"
	"github.com/google/uuid"
	testsupport "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/http/testkit"
)

func TestUsersByIDNotFound(t *testing.T) {
	testsupport.Run(t, "missing users return not found over http", func(cfg *axiom.Config, env testsupport.Env) {
		tt := testsupport.CaseT(t, cfg)
		admin := testsupport.LoginAdmin(tt, env)
		missingPath := testsupport.UserPath(uuid.NewString())

		cfg.Step("get missing user returns not found", func() {
			resp := admin.Get(tt, missingPath, nil)
			testsupport.RequireStatus(tt, resp, http.StatusNotFound)
		})

		cfg.Step("update missing user returns not found", func() {
			resp := admin.Put(tt, missingPath, map[string]any{"name": "Ghost"})
			testsupport.RequireStatus(tt, resp, http.StatusNotFound)
		})

		cfg.Step("delete missing user returns not found", func() {
			resp := admin.Delete(tt, missingPath)
			testsupport.RequireStatus(tt, resp, http.StatusNotFound)
		})
	})
}
