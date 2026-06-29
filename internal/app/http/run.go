package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stratflow-labs/stratflow/internal/app/bootstrap"
	"github.com/stratflow-labs/stratflow/internal/app/config"

	"github.com/samber/do"
)

func Run(ctx context.Context, injector *do.Injector, cfg *config.Config) error {
	engine := NewServer(injector, cfg)
	server := &http.Server{
		Addr:              cfg.HTTP.ListenAddr,
		Handler:           engine,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	}

	if err := bootstrap.GracefulServe(ctx, server, cfg); err != nil {
		return fmt.Errorf("serve http server: %w", err)
	}
	return nil
}
