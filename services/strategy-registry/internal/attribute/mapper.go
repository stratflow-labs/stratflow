package attribute

import (
	"time"

	attributevalue "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attributevalue"
	attributedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type AttributeView struct {
	ID          uuid.UUID
	StrategyID  uuid.UUID
	Slug        string
	Name        string
	Description string
	Values      []attributevalue.AttributeValueView
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AttributeValueView = attributevalue.AttributeValueView
type AttributeValueRelationView = attributevalue.AttributeValueRelationView

func attributeViewFromDomain(s attributedomain.Attribute) AttributeView {
	return AttributeView{
		ID:          s.ID,
		StrategyID:  s.StrategyID,
		Slug:        s.Slug,
		Name:        s.Name,
		Description: s.Description,
		Values:      make([]attributevalue.AttributeValueView, 0),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func attributeViewsFromDomain(items []attributedomain.Attribute) []AttributeView {
	out := make([]AttributeView, len(items))
	for i := range items {
		out[i] = attributeViewFromDomain(items[i])
	}
	return out
}
