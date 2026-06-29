package testkit

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stratflow-labs/stratflow/internal/authkit"
	"github.com/stratflow-labs/stratflow/internal/authz"
)

type FakeVerifier struct{}

func (FakeVerifier) Verify(_ context.Context, token string) (authkit.Claims, error) {
	switch strings.TrimSpace(token) {
	case "admin-token":
		return authkit.Claims{UserID: "admin-1", Role: "admin"}, nil
	case "manager-token":
		return authkit.Claims{UserID: "manager-1", Role: "manager"}, nil
	case "user-token":
		return authkit.Claims{UserID: "user-1", Role: "user"}, nil
	default:
		return authkit.Claims{}, errors.New("invalid token")
	}
}

func StatusHandler(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
	})
}

func LoadEngineFromRepo(t *testing.T) *authz.Engine {
	t.Helper()
	base := repoRootPath(t, "internal/authz")
	engine, err := authz.NewEngine(authz.EngineConfig{
		ModelPath:  filepath.Join(base, "model.conf"),
		PolicyPath: filepath.Join(base, "policy.csv"),
		RoutesPath: filepath.Join(base, "routes.yaml"),
	})
	if err != nil {
		t.Fatalf("load engine: %v", err)
	}
	return engine
}

func RepoRootPath(t *testing.T, subpath string) string {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	root := filepath.Clean(filepath.Join(dir, "..", "..", "..", "..", ".."))
	return filepath.Join(root, subpath)
}

func repoRootPath(t *testing.T, subpath string) string {
	t.Helper()
	return RepoRootPath(t, subpath)
}
