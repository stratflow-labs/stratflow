package attributevalue

import (
	"time"

	attributeValuedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type AttributeValueView struct {
	ID          uuid.UUID
	StrategyID  uuid.UUID
	AttributeID uuid.UUID
	Slug        string
	Value       string
	Relations   []AttributeValueRelationView
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AttributeValueRelationView struct {
	FromAttributeID uuid.UUID
	FromValueID     uuid.UUID
	ToAttributeID   uuid.UUID
	ToValueID       uuid.UUID
	ToAttributeSlug string
	ToValueSlug     string
}

func attributeValueViewFromDomain(s attributeValuedomain.AttributeValue) AttributeValueView {
	return AttributeValueView{
		ID:          s.ID,
		StrategyID:  s.StrategyID,
		AttributeID: s.AttributeID,
		Slug:        s.Slug,
		Value:       s.Value,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func attributeValueViewsFromDomain(items []attributeValuedomain.AttributeValue) []AttributeValueView {
	out := make([]AttributeValueView, len(items))
	for i := range items {
		out[i] = attributeValueViewFromDomain(items[i])
	}
	return out
}
