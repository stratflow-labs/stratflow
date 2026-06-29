package attribute

import (
	"context"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	attributedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type CreateInput struct {
	StrategyID  uuid.UUID
	Slug        string
	Name        string
	Description string
}

func (s *WriteService) Create(ctx context.Context, input CreateInput) (AttributeView, error) {
	if input.StrategyID == uuid.Nil {
		return AttributeView{}, apperr.NotFoundError[attributedomain.Strategy]()
	}

	now := s.clock.Now()
	id := uuid.New()
	attribute, err := attributedomain.NewAttribute(
		id,
		input.StrategyID,
		input.Slug,
		input.Name,
		input.Description,
		now,
	)
	if err != nil {
		return AttributeView{}, apperr.Invalid[attributedomain.Attribute](
			"create",
			"validation",
			"attribute validation failed",
			err,
		)
	}

	created, err := s.attributeRepo.Create(ctx, attribute)
	if err != nil {
		return AttributeView{}, fmt.Errorf("create attribute: %w", err)
	}

	return attributeViewFromDomain(created), nil
}
