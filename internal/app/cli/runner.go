package cli

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/stratflow-labs/stratflow/internal/foundation/logger"
)

type CommandHandlers struct {
	Serve   func(context.Context) error
	Migrate func(context.Context) error
	Seed    func(context.Context) error
}

// Run routes CLI subcommands to handlers.
func Run(ctx context.Context, args []string, usage string, handlers CommandHandlers) error {
	command := "serve"
	if len(args) > 0 {
		command = strings.ToLower(args[0])
	}

	switch command {
	case "", "serve":
		if handlers.Serve == nil {
			return errors.New("serve handler is nil")
		}
		return handlers.Serve(ctx)
	case "migrate":
		if handlers.Migrate == nil {
			return errors.New("migrate handler is nil")
		}
		return handlers.Migrate(ctx)
	case "seed":
		if handlers.Seed == nil {
			return errors.New("seed handler is nil")
		}
		return handlers.Seed(ctx)
	case "help", "--help", "-h":
		logger.Info(usage)
		return nil
	default:
		return errors.New("unknown command " + strconv.Quote(command) + ". " + usage)
	}
}
