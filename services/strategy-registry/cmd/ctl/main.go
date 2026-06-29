package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/stratflow-labs/stratflow/internal/foundation/logger"
	"github.com/stratflow-labs/stratflow/services/strategy-registry/app/bootstrap"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		logger.Err("db command failed", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("expected subcommand: migrate or seed")
	}

	switch args[0] {
	case "migrate":
		return runMigrate(ctx)
	case "seed":
		return runSeed(ctx)
	default:
		return fmt.Errorf("unknown subcommand %q", args[0])
	}
}

func runMigrate(ctx context.Context) error {
	return bootstrap.Migrate(ctx)
}

func runSeed(ctx context.Context) error {
	return bootstrap.Seed(ctx)
}
