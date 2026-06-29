package main

import (
	"context"
	"os"

	"github.com/stratflow-labs/stratflow/internal/foundation/logger"
	"github.com/stratflow-labs/stratflow/services/strategy-registry/app/bootstrap"
)

func main() {
	if err := bootstrap.Run(context.Background(), os.Args[1:]); err != nil {
		logger.Err("application exited with error", err)
		os.Exit(1)
	}
}
