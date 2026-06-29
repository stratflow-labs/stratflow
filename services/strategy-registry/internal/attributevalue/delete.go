package attributevalue

import (
	"context"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	attributeValuedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type DeleteInput struct {
	ID          uuid.UUID
	StrategyID  uuid.UUID
	AttributeID uuid.UUID
}

func (s *Service) Delete(ctx context.Context, input DeleteInput) error {
	if input.ID == uuid.Nil {
		return apperr.NotFoundError[attributeValuedomain.AttributeValue]()
	}
	if input.StrategyID == uuid.Nil {
		return apperr.NotFoundError[attributeValuedomain.Strategy]()
	}
	if input.AttributeID == uuid.Nil {
		return apperr.NotFoundError[attributeValuedomain.Attribute]()
	}

	if err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		current, err := s.attributeValueRepo.GetByID(txCtx, input.ID)
		if err != nil {
			return fmt.Errorf("load attributeValue: %w", err)
		}
		if current.StrategyID != input.StrategyID {
			return apperr.NotFoundError[attributeValuedomain.AttributeValue]()
		}
		if current.AttributeID != input.AttributeID {
			return apperr.NotFoundError[attributeValuedomain.AttributeValue]()
		}

		if err := s.attributeValueRepo.Delete(txCtx, input.ID); err != nil {
			return fmt.Errorf("delete attributeValue: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("delete attributeValue: %w", err)
	}
	return nil
}
