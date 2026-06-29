package http

import (
	"github.com/stratflow-labs/stratflow/internal/app/config"
	"github.com/stratflow-labs/stratflow/internal/app/middlewares"
	"github.com/stratflow-labs/stratflow/internal/authkit"
	"github.com/stratflow-labs/stratflow/internal/authz"
	"github.com/stratflow-labs/stratflow/internal/httpserver"
	hsMiddleware "github.com/stratflow-labs/stratflow/internal/httpserver/middleware"

	"github.com/samber/do"
)

func NewServer(inj *do.Injector, _ *config.Config) *httpserver.Mux {
	r := httpserver.New()
	r.Use(
		hsMiddleware.RequestID,
		hsMiddleware.Logger,
		hsMiddleware.Recoverer,
		middlewares.CORS(),
	)

	// Optional auth middleware: enabled when AccessTokenVerifier is registered.
	if verifier, err := do.Invoke[authkit.AccessTokenVerifier](inj); err == nil && verifier != nil {
		skip := authkit.PathPrefixSkipper(
			"/api/auth/login",
			"/connect/stratflow.identity.v1.IdentityService/Login",
			"/connect/stratflow.identity.v1.IdentityService/VerifyToken",
		)
		r.Use(authkit.AuthenticateWithSkipper(verifier, skip))

		// Authorization (runs after authentication).
		if engine, err := do.Invoke[*authz.Engine](inj); err == nil && engine != nil {
			r.Use(authz.Middleware(engine))
		}
	}

	RegisterRoutes(r, inj)
	return r
}
