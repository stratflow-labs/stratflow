package wiring

import (
	apphttp "github.com/stratflow-labs/stratflow/internal/app/http"
	identityconnect "github.com/stratflow-labs/stratflow/services/identity/internal/adapters/connect"
	identityhttp "github.com/stratflow-labs/stratflow/services/identity/internal/adapters/http"

	"github.com/samber/do"
)

func Register(inj *do.Injector) {
	RegisterProviders(inj)
	registerHTTPRoutes(inj)
}

func registerHTTPRoutes(inj *do.Injector) {
	do.Provide(inj, func(i *do.Injector) ([]apphttp.RouteRegistrar, error) {
		gatewayHandler := do.MustInvoke[*identityhttp.GatewayHandler](i)
		connectHandler := do.MustInvoke[*identityconnect.Handler](i)
		return []apphttp.RouteRegistrar{gatewayHandler, connectHandler}, nil
	})
}
