package attribute

import (
	"context"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	attributedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type GetInput struct {
	ID         uuid.UUID
	StrategyID uuid.UUID
}

func (s *ReadService) Get(ctx context.Context, input GetInput) (AttributeView, error) {
	if input.ID == uuid.Nil {
		return AttributeView{}, apperr.NotFoundError[attributedomain.Attribute]()
	}
	if input.StrategyID == uuid.Nil {
		return AttributeView{}, apperr.NotFoundError[attributedomain.Strategy]()
	}
	attribute, err := s.attributeRepo.GetByID(ctx, input.ID)
	if err != nil {
		return AttributeView{}, fmt.Errorf("get attribute by id: %w", err)
	}
	if attribute.StrategyID != input.StrategyID {
		return AttributeView{}, apperr.NotFoundError[attributedomain.Attribute]()
	}

	view := attributeViewFromDomain(attribute)
	if s.attributeValueRepo == nil {
		return view, nil
	}

	valuesByAttributeID, err := s.attributeValueRepo.ListValueGraphByAttributeIDs(ctx, input.StrategyID, []uuid.UUID{input.ID})
	if err != nil {
		return AttributeView{}, fmt.Errorf("get attribute value graph: %w", err)
	}
	if values, ok := valuesByAttributeID[input.ID]; ok {
		view.Values = values
	}

	return view, nil
}
