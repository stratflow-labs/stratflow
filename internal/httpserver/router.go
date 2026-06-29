package httpserver

import (
	"net/http"
)

// Router is the minimal router interface required by the application.
type Router interface {
	http.Handler
	Use(middlewares ...func(http.Handler) http.Handler)
	Mount(pattern string, h http.Handler)
}

// Mux is a thin wrapper around http.ServeMux with middleware chain support.
type Mux struct {
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

// New returns a new router.
func New() *Mux {
	return &Mux{mux: http.NewServeMux()}
}

// Use adds middleware to the chain in registration order.
func (m *Mux) Use(middlewares ...func(http.Handler) http.Handler) {
	m.middlewares = append(m.middlewares, middlewares...)
}

// Mount registers a handler for the given pattern and wraps it in the middleware chain.
func (m *Mux) Mount(pattern string, h http.Handler) {
	m.mux.Handle(pattern, h)
}

// ServeHTTP implements http.Handler.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := http.Handler(m.mux)
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		handler = m.middlewares[i](handler)
	}
	handler.ServeHTTP(w, r)
}
