package attribute

import (
	"context"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	attributedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type UpdateInput struct {
	ID          uuid.UUID
	StrategyID  uuid.UUID
	Slug        *string
	Name        *string
	Description *string
}

func (in UpdateInput) HasChanges() bool {
	return in.Slug != nil || in.Name != nil || in.Description != nil
}

func (s *WriteService) Update(ctx context.Context, input UpdateInput) (AttributeView, error) {
	if input.ID == uuid.Nil {
		return AttributeView{}, apperr.NotFoundError[attributedomain.Attribute]()
	}

	if input.StrategyID == uuid.Nil {
		return AttributeView{}, apperr.NotFoundError[attributedomain.Strategy]()
	}

	if !input.HasChanges() {
		return AttributeView{}, apperr.UpdateEmptyError[attributedomain.Attribute]()
	}

	var updated attributedomain.Attribute
	err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		current, err := s.attributeRepo.GetByID(txCtx, input.ID)
		if err != nil {
			return fmt.Errorf("load attribute: %w", err)
		}
		if current.StrategyID != input.StrategyID {
			return apperr.NotFoundError[attributedomain.Attribute]()
		}

		if input.Slug != nil {
			current.Slug = attributedomain.SanitizeString(*input.Slug)
		}
		if input.Name != nil {
			current.Name = attributedomain.SanitizeString(*input.Name)
		}
		if input.Description != nil {
			current.Description = attributedomain.SanitizeString(*input.Description)
		}

		if fields := attributedomain.ValidateAttribute(current.Slug, current.Name, current.Description); len(fields) > 0 {
			return apperr.Validation[attributedomain.Attribute]("update", fields)
		}

		current.UpdatedAt = s.clock.Now()

		updated, err = s.attributeRepo.Update(txCtx, &current)
		if err != nil {
			return fmt.Errorf("update attribute: %w", err)
		}

		return nil
	})
	if err != nil {
		return AttributeView{}, fmt.Errorf("update attribute transaction: %w", err)
	}

	return attributeViewFromDomain(updated), nil
}
