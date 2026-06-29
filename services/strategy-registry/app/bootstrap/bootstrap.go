package bootstrap

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"

	appcli "github.com/stratflow-labs/stratflow/internal/app/cli"
	"github.com/stratflow-labs/stratflow/internal/app/config"
	apphttp "github.com/stratflow-labs/stratflow/internal/app/http"
	"github.com/stratflow-labs/stratflow/internal/app/migrations"
	appruntime "github.com/stratflow-labs/stratflow/internal/app/runtime"
	"github.com/stratflow-labs/stratflow/internal/authkit"
	"github.com/stratflow-labs/stratflow/internal/foundation/logger"
	"github.com/stratflow-labs/stratflow/internal/grpcserver"
	"github.com/stratflow-labs/stratflow/services/strategy-registry/app/wiring"
	"github.com/stratflow-labs/stratflow/services/strategy-registry/db"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	strategygrpc "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/grpc"

	"github.com/samber/do"

	"google.golang.org/grpc"
)

const usage = "usage: strategy-registry-service [serve|migrate|seed|help]"

// Run bootstraps the CLI application.
func Run(ctx context.Context, args []string) error {
	return appcli.Run(ctx, args, usage, appcli.CommandHandlers{
		Serve:   Serve,
		Migrate: Migrate,
		Seed:    Seed,
	})
}

// Serve starts the HTTP server.
func Serve(ctx context.Context) error {
	return appcli.WithConfigAndDB(ctx, LoadConfig, OpenDB, func(ctx context.Context, cfg *config.Config, dbConn *sql.DB) error {
		injector := wiring.BuildContainer(cfg, dbConn)
		wiring.Register(injector)
		startBackgroundProcesses(ctx, injector)

		httpHandler := apphttp.NewServer(injector, cfg)
		grpcServer, err := buildGRPCServer(injector)
		if err != nil {
			return err
		}

		httpServer := &http.Server{
			Addr:              cfg.HTTP.ListenAddr,
			Handler:           httpHandler,
			ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		}

		grpcListener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", cfg.GRPC.ListenAddr)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := grpcListener.Close(); cerr != nil && !errors.Is(cerr, net.ErrClosed) {
				logger.Err("close grpc listener", cerr)
			}
		}()

		return appruntime.ServeHTTPAndGRPC(ctx, appruntime.HTTPAndGRPCConfig{
			ServiceName:   "strategy registry",
			HTTPServer:    httpServer,
			GRPCServer:    grpcServer,
			GRPCListener:  grpcListener,
			ShutdownGrace: cfg.HTTP.ShutdownGrace,
		})
	})
}

// Migrate runs database migrations.
func Migrate(ctx context.Context) error {
	return appcli.Migrate(ctx, LoadConfig, OpenDB, RunMigrations)
}

// Seed runs database seeders.
func Seed(ctx context.Context) error {
	return appcli.SeedNoContext(ctx, LoadConfig, OpenDB, RunSeeders)
}

// OpenDB opens database connection.
func OpenDB(cfg *config.Config) (*sql.DB, func() error, error) {
	if cfg == nil {
		return nil, nil, errors.New("config is nil")
	}
	return db.NewDB(cfg)
}

// RunMigrations runs database migrations for strategy-registry service.
func RunMigrations(sqlDB *sql.DB) error {
	return migrations.RunServiceMigrations(sqlDB, "strategy-registry")
}

// RunSeeders runs database seeders for strategy-registry service.
func RunSeeders(sqlDB *sql.DB) error {
	return migrations.RunServiceSeeds(sqlDB, "strategy-registry")
}

// todo just check if need return nil, errors.New("injector is nil") or return nil, errors.New("resolve strategy grpc handler: nil handler") etc
// or samber already check it in do.MustInvoke
func buildGRPCServer(injector *do.Injector) (*grpc.Server, error) {
	if injector == nil {
		return nil, errors.New("injector is nil")
	}

	cfg := do.MustInvoke[config.Config](injector)
	authVerifier := do.MustInvoke[authkit.AccessTokenVerifier](injector)
	server := grpcserver.New(grpcserver.Config{
		AuthVerifier:     authVerifier,
		EnableReflection: cfg.AppEnv == "localhost",
	})

	handler := do.MustInvoke[*strategygrpc.Handler](injector)
	strategyregistryv1.RegisterStrategyRegistryServiceServer(server, handler)
	return server, nil
}

// startBackgroundProcesses initializes background workers (placeholder for future use).
func startBackgroundProcesses(_ context.Context, _ *do.Injector) {}
