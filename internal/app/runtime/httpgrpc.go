package runtime

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/stratflow-labs/stratflow/internal/foundation/logger"
	"github.com/stratflow-labs/stratflow/internal/grpcserver"
)

type HTTPAndGRPCConfig struct {
	ServiceName   string
	HTTPServer    *http.Server
	GRPCServer    *grpc.Server
	GRPCListener  net.Listener
	ShutdownGrace time.Duration
}

func ServeHTTPAndGRPC(ctx context.Context, cfg HTTPAndGRPCConfig) error {
	if ctx == nil {
		return errors.New("runtime context is nil")
	}
	if cfg.HTTPServer == nil {
		return errors.New("http server is nil")
	}
	if cfg.GRPCServer == nil {
		return errors.New("grpc server is nil")
	}
	if cfg.GRPCListener == nil {
		return errors.New("grpc listener is nil")
	}

	errCh := make(chan error, 2)
	runCtx, cancelRun := context.WithCancel(ctx)
	defer cancelRun()

	sigCtx, stop := signal.NotifyContext(runCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := cfg.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	go func() {
		errCh <- grpcserver.Serve(sigCtx, cfg.GRPCServer, cfg.GRPCListener, cfg.ShutdownGrace)
	}()

	select {
	case <-sigCtx.Done():
		logger.Info("shutdown signal received, stopping " + cfg.ServiceName + " servers")

		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(sigCtx), cfg.ShutdownGrace)
		defer cancel()

		httpErr := cfg.HTTPServer.Shutdown(shutdownCtx)
		first := <-errCh
		second := <-errCh

		if httpErr != nil && !errors.Is(httpErr, context.Canceled) {
			return httpErr
		}
		if first != nil {
			return first
		}
		if second != nil {
			return second
		}

		logger.Info(cfg.ServiceName + " http and grpc servers stopped gracefully")
		return nil

	case err := <-errCh:
		cancelRun()
		_ = cfg.HTTPServer.Close()
		if err != nil {
			return err
		}
		if other := <-errCh; other != nil {
			return other
		}
		return nil
	}
}
