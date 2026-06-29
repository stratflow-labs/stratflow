package attribute

import (
	"context"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	attributedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type DeleteInput struct {
	ID         uuid.UUID
	StrategyID uuid.UUID
}

func (s *WriteService) Delete(ctx context.Context, input DeleteInput) error {
	if input.ID == uuid.Nil {
		return apperr.NotFoundError[attributedomain.Attribute]()
	}
	if input.StrategyID == uuid.Nil {
		return apperr.NotFoundError[attributedomain.Strategy]()
	}

	if err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		current, err := s.attributeRepo.GetByID(txCtx, input.ID)
		if err != nil {
			return fmt.Errorf("load attribute: %w", err)
		}
		if current.StrategyID != input.StrategyID {
			return apperr.NotFoundError[attributedomain.Attribute]()
		}

		if err := s.attributeRepo.Delete(txCtx, input.ID); err != nil {
			return fmt.Errorf("delete attribute: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("delete attribute: %w", err)
	}
	return nil
}
