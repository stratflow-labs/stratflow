package grpcserver

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/stratflow-labs/stratflow/internal/authkit"
)

type Config struct {
	AuthVerifier           authkit.AccessTokenVerifier
	AuthSkipper            AuthSkipper
	Metrics                MetricsRecorder
	EnableReflection       bool
	ExtraUnaryInterceptors []grpc.UnaryServerInterceptor
}

func New(cfg Config) *grpc.Server {
	interceptors := []grpc.UnaryServerInterceptor{
		RecoveryUnaryInterceptor(),
		RequestMetadataUnaryInterceptor(),
		LoggingUnaryInterceptor(),
		MetricsUnaryInterceptor(cfg.Metrics),
		AuthUnaryInterceptor(cfg.AuthVerifier, cfg.AuthSkipper),
	}
	interceptors = append(interceptors, cfg.ExtraUnaryInterceptors...)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
	)

	if cfg.EnableReflection {
		reflection.Register(server)
	}

	return server
}
