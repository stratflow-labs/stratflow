package grpcserver

import (
	"context"
	"errors"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/stratflow-labs/stratflow/internal/foundation/logger"
)

func Serve(ctx context.Context, server *grpc.Server, listener net.Listener, shutdownGrace time.Duration) error {
	if ctx == nil {
		return errors.New("grpc server context is nil")
	}
	if server == nil {
		return errors.New("grpc server is nil")
	}
	if listener == nil {
		return errors.New("grpc listener is nil")
	}

	errCh := make(chan error, 1)

	go func() {
		err := server.Serve(listener)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		logger.Info("grpc shutdown signal received")
		stopGracefully(server, shutdownGrace)
		_ = listener.Close()
		return <-errCh
	case err := <-errCh:
		_ = listener.Close()
		return err
	}
}

func stopGracefully(server *grpc.Server, shutdownGrace time.Duration) {
	done := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(done)
	}()

	if shutdownGrace <= 0 {
		<-done
		return
	}

	select {
	case <-done:
	case <-time.After(shutdownGrace):
		logger.Warn("grpc graceful shutdown timed out, forcing stop")
		server.Stop()
		<-done
	}
}
