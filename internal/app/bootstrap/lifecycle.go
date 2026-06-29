package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/stratflow-labs/stratflow/internal/app/config"
	"github.com/stratflow-labs/stratflow/internal/foundation/logger"
)

func GracefulServe(ctx context.Context, server *http.Server, cfg *config.Config) error {
	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	sigCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-sigCtx.Done():
		logger.Info("shutdown signal received, stopping HTTP server")

		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(sigCtx), cfg.HTTP.ShutdownGrace)
		defer cancel()

		shutdownErr := server.Shutdown(shutdownCtx)
		listenErr := <-errCh

		if shutdownErr != nil && !errors.Is(shutdownErr, context.Canceled) {
			return fmt.Errorf("shutdown HTTP server: %w", shutdownErr)
		}
		if listenErr != nil {
			return fmt.Errorf("HTTP server stopped with error: %w", listenErr)
		}

		logger.Info("HTTP server stopped gracefully")
		return nil

	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("HTTP server failed: %w", err)
		}
		return nil
	}
}
