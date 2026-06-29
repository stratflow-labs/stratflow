package attributevalue

import (
	"context"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	attributeValuedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type GetInput struct {
	ID          uuid.UUID
	StrategyID  uuid.UUID
	AttributeID uuid.UUID
}

func (s *Service) Get(ctx context.Context, input GetInput) (AttributeValueView, error) {
	if input.ID == uuid.Nil {
		return AttributeValueView{}, apperr.NotFoundError[attributeValuedomain.AttributeValue]()
	}
	if input.StrategyID == uuid.Nil {
		return AttributeValueView{}, apperr.NotFoundError[attributeValuedomain.Strategy]()
	}
	if input.AttributeID == uuid.Nil {
		return AttributeValueView{}, apperr.NotFoundError[attributeValuedomain.Attribute]()
	}
	attributeValue, err := s.attributeValueRepo.GetByID(ctx, input.ID)
	if err != nil {
		return AttributeValueView{}, fmt.Errorf("get attributeValue by id: %w", err)
	}
	if attributeValue.StrategyID != input.StrategyID {
		return AttributeValueView{}, apperr.NotFoundError[attributeValuedomain.AttributeValue]()
	}
	if attributeValue.AttributeID != input.AttributeID {
		return AttributeValueView{}, apperr.NotFoundError[attributeValuedomain.AttributeValue]()
	}

	view := attributeValueViewFromDomain(attributeValue)
	relationsByFromValueID, err := s.attributeValueRepo.ListRelationsByFromValueIDsForAttributeValues(ctx, input.StrategyID, []uuid.UUID{input.ID})
	if err != nil {
		return AttributeValueView{}, fmt.Errorf("get attributeValue relations: %w", err)
	}
	if relations, ok := relationsByFromValueID[input.ID]; ok {
		view.Relations = relations
	}

	return view, nil
}
