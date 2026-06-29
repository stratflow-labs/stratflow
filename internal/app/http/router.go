package http

import (
	"github.com/stratflow-labs/stratflow/internal/httpserver"

	"github.com/samber/do"
)

type RouteRegistrar interface {
	RegisterRoutes(r httpserver.Router)
}

type RouteRegistrarFunc func(r httpserver.Router)

func (f RouteRegistrarFunc) RegisterRoutes(r httpserver.Router) {
	f(r)
}

func RegisterRoutes(r httpserver.Router, injector *do.Injector) {
	registrars := do.MustInvoke[[]RouteRegistrar](injector)
	for _, registrar := range registrars {
		registrar.RegisterRoutes(r)
	}
}
