package wiring

import (
	apphttp "github.com/stratflow-labs/stratflow/internal/app/http"
	strategyconnect "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/connect"
	strategyhttp "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/http"

	"github.com/samber/do"
)

func Register(inj *do.Injector) {
	RegisterProviders(inj)
	registerHTTPRoutes(inj)
}

func registerHTTPRoutes(inj *do.Injector) {
	do.Provide(inj, func(i *do.Injector) ([]apphttp.RouteRegistrar, error) {
		gatewayHandler := do.MustInvoke[*strategyhttp.GatewayHandler](i)
		connectHandler := do.MustInvoke[*strategyconnect.Handler](i)
		return []apphttp.RouteRegistrar{gatewayHandler, connectHandler}, nil
	})
}
