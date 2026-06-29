package strategy

import (
	"context"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	strategydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

func (s *WriteService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return apperr.NotFoundError[strategydomain.Strategy]()
	}

	if err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		if err := s.strategyRepo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete strategy: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("delete strategy transaction: %w", err)
	}

	return nil
}
