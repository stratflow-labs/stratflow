package wiring

import (
	"database/sql"
	"os"
	"strings"

	"github.com/stratflow-labs/stratflow/internal/authkit"
	"github.com/stratflow-labs/stratflow/internal/authz"
	systemclock "github.com/stratflow-labs/stratflow/internal/foundation/clock/system"
	txpkg "github.com/stratflow-labs/stratflow/internal/foundation/tx"
	"github.com/stratflow-labs/stratflow/internal/foundation/tx/sqltx"
	identityclient "github.com/stratflow-labs/stratflow/services/identity/client"
	bybitauth "github.com/stratflow-labs/stratflow/services/strategy-registry/app/authverifier"
	strategyconnect "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/connect"
	strategygrpc "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/grpc"
	strategyhttp "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/http"
	strategypostgres "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/postgres"
	attribute "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attribute"
	attributeValue "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attributevalue"
	strategy "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/strategy"
	strategygraph "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/strategygraph"

	"github.com/samber/do"
)

// RegisterProviders wires all module providers into injector.
func RegisterProviders(inj *do.Injector) {
	registerShared(inj)
	registerStrategyServer(inj)
	registerHTTPGateway(inj)
}

func registerShared(inj *do.Injector) {
	do.Provide(inj, func(i *do.Injector) (txpkg.Manager, error) {
		db := do.MustInvoke[*sql.DB](i)
		return &sqltx.Manager{DB: db}, nil
	})

	do.Provide(inj, func(i *do.Injector) (authkit.AccessTokenVerifier, error) {
		identityGRPCTarget := strings.TrimSpace(os.Getenv("IDENTITY_GRPC_URL"))
		identityVerifier := bybitauth.VerifyTokenAdapter{Verifier: identityclient.NewTokenVerifier(identityGRPCTarget)}
		serviceToken := strings.TrimSpace(os.Getenv("MASTER_SESSION_TOKEN"))

		return bybitauth.ServiceTokenFallbackVerifier{
			ServiceToken:     serviceToken,
			IdentityVerifier: identityVerifier,
		}, nil
	})

	do.Provide(inj, func(i *do.Injector) (*authz.Engine, error) {
		cfg := authz.EngineConfig{
			ModelPath:  "internal/authz/model.conf",
			PolicyPath: "internal/authz/policy.csv",
			RoutesPath: "internal/authz/routes.yaml",
		}
		return authz.NewEngine(cfg)
	})
}

func registerStrategyServer(inj *do.Injector) {
	do.Provide(inj, func(i *do.Injector) (*strategygrpc.Handler, error) {
		db := do.MustInvoke[*sql.DB](i)
		strategyRepo := strategypostgres.NewStrategyRepository(db)
		attributeRepo := strategypostgres.NewAttributeRepository(db)
		attributeValueRepo := strategypostgres.NewAttributeValueRepository(db)
		txManager := do.MustInvoke[txpkg.Manager](i)
		clk := systemclock.New()

		strategyReadSvc := strategy.NewReadService(strategyRepo)
		strategyWriteSvc := strategy.NewWriteService(strategyRepo, strategyRepo, txManager, clk)
		strategySvc := strategy.NewService(strategyReadSvc, strategyWriteSvc)

		attributeReadSvc := attribute.NewReadService(attributeRepo, attributeValueRepo)
		attributeWriteSvc := attribute.NewWriteService(attributeRepo, txManager, clk)
		attributeSvc := attribute.NewService(attributeReadSvc, attributeWriteSvc)
		attributeValueSvc := attributeValue.NewService(attributeValueRepo, txManager, clk)
		strategyGraphSvc := strategygraph.NewService(
			strategyRepo,
			attributeRepo,
			attributeSvc,
			attributeValueRepo,
			txManager,
			clk,
		)

		return strategygrpc.NewHandler(strategygrpc.HandlerDependencies{
			Strategies:      strategySvc,
			Attributes:      attributeSvc,
			AttributeValues: attributeValueSvc,
			StrategyGraph:   strategyGraphSvc,
		}), nil
	})
}

func registerHTTPGateway(inj *do.Injector) {
	do.Provide(inj, func(i *do.Injector) (*strategyhttp.GatewayHandler, error) {
		handler := do.MustInvoke[*strategygrpc.Handler](i)
		return strategyhttp.NewGatewayHandler(handler), nil
	})

	do.Provide(inj, func(i *do.Injector) (*strategyconnect.Handler, error) {
		handler := do.MustInvoke[*strategygrpc.Handler](i)
		return strategyconnect.NewHandler(handler), nil
	})
}
