package wiring

import (
	"database/sql"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/app/config"
	"github.com/stratflow-labs/stratflow/internal/authkit"
	"github.com/stratflow-labs/stratflow/internal/authz"
	systemclock "github.com/stratflow-labs/stratflow/internal/foundation/clock/system"
	"github.com/stratflow-labs/stratflow/internal/foundation/password"
	txpkg "github.com/stratflow-labs/stratflow/internal/foundation/tx"
	"github.com/stratflow-labs/stratflow/internal/foundation/tx/sqltx"
	"github.com/stratflow-labs/stratflow/services/identity/app/authverifier"
	identityconnect "github.com/stratflow-labs/stratflow/services/identity/internal/adapters/connect"
	identitygrpc "github.com/stratflow-labs/stratflow/services/identity/internal/adapters/grpc"
	identityhttp "github.com/stratflow-labs/stratflow/services/identity/internal/adapters/http"
	identitypostgres "github.com/stratflow-labs/stratflow/services/identity/internal/adapters/postgres"
	auth "github.com/stratflow-labs/stratflow/services/identity/internal/auth"
	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"
	opaquetoken "github.com/stratflow-labs/stratflow/services/identity/internal/token/opaque"
	user "github.com/stratflow-labs/stratflow/services/identity/internal/user"

	"github.com/samber/do"
)

// RegisterProviders wires all module providers into injector.
func RegisterProviders(inj *do.Injector) {
	registerShared(inj)
	registerIdentity(inj)
}

func registerShared(inj *do.Injector) {
	do.Provide(inj, func(i *do.Injector) (*opaquetoken.Service, error) {
		db := do.MustInvoke[*sql.DB](i)
		cfg := do.MustInvoke[config.Config](i)

		if cfg.Security.TokenHashSecret == "" {
			return nil, fmt.Errorf("TOKEN_HASH_SECRET is required")
		}

		return opaquetoken.NewService(db, &opaquetoken.Config{
			TokenHashSecret: cfg.Security.TokenHashSecret,
		})
	})

	do.Provide(inj, func(i *do.Injector) (auth.TokenService, error) {
		return do.MustInvoke[*opaquetoken.Service](i), nil
	})

	do.Provide(inj, func(i *do.Injector) (auth.AccessTokenRevoker, error) {
		return do.MustInvoke[*opaquetoken.Service](i), nil
	})

	do.Provide(inj, func(i *do.Injector) (*auth.Service, error) {
		tokens := do.MustInvoke[auth.TokenService](i)
		return auth.NewService(nil, nil, tokens, nil, nil), nil
	})

	do.Provide(inj, func(i *do.Injector) (authkit.AccessTokenVerifier, error) {
		authService := do.MustInvoke[*auth.Service](i)
		return &authverifier.VerifyTokenAdapter{Auth: authService}, nil
	})

	do.Provide(inj, func(i *do.Injector) (*authz.Engine, error) {
		cfg := authz.EngineConfig{
			ModelPath:  "pkg/authz/model.conf",
			PolicyPath: "pkg/authz/policy.csv",
			RoutesPath: "pkg/authz/routes.yaml",
		}
		return authz.NewEngine(cfg)
	})

	do.Provide(inj, func(i *do.Injector) (*password.BcryptHasher, error) {
		hasher := password.NewBcryptHasher(0)
		return &hasher, nil
	})

	do.Provide(inj, func(i *do.Injector) (txpkg.Manager, error) {
		db := do.MustInvoke[*sql.DB](i)
		return &sqltx.Manager{DB: db}, nil
	})
}

func registerIdentity(inj *do.Injector) {
	do.Provide(inj, func(i *do.Injector) (*identitygrpc.Handler, error) {
		db := do.MustInvoke[*sql.DB](i)
		tokens := do.MustInvoke[auth.TokenService](i)
		accessRevoker := do.MustInvoke[auth.AccessTokenRevoker](i)
		passwords := do.MustInvoke[*password.BcryptHasher](i)
		txManager := do.MustInvoke[txpkg.Manager](i)
		clk := systemclock.New()
		validator := identitydomain.NewValidator()

		credentialRepo := identitypostgres.NewCredentialRepository(db)
		userRepo := identitypostgres.NewUserRepository(db)

		authService := auth.NewService(
			credentialRepo,
			passwords,
			tokens,
			accessRevoker,
			txManager,
		)
		users := user.NewService(
			userRepo,
			txManager,
			passwords,
			validator,
			clk,
		)

		return identitygrpc.NewHandler(identitygrpc.HandlerDependencies{
			Auth:  authService,
			Users: users,
		}), nil
	})

	do.Provide(inj, func(i *do.Injector) (*identityhttp.GatewayHandler, error) {
		handler := do.MustInvoke[*identitygrpc.Handler](i)
		return identityhttp.NewGatewayHandler(handler), nil
	})

	do.Provide(inj, func(i *do.Injector) (*identityconnect.Handler, error) {
		handler := do.MustInvoke[*identitygrpc.Handler](i)
		return identityconnect.NewHandler(handler), nil
	})
}
