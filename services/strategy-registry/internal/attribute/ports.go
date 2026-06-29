package attribute

import (
	"context"

	"github.com/stratflow-labs/stratflow/internal/foundation/clock"
	attributevalue "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attributevalue"
	attributedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type AttributeRepository interface {
	Create(ctx context.Context, attribute *attributedomain.Attribute) (attributedomain.Attribute, error)
	GetByID(ctx context.Context, id uuid.UUID) (attributedomain.Attribute, error)
	GetBySlug(ctx context.Context, strategyID uuid.UUID, slug string) (attributedomain.Attribute, error)
	List(ctx context.Context, filter ListFilter) ([]attributedomain.Attribute, int64, error)
	Update(ctx context.Context, attribute *attributedomain.Attribute) (attributedomain.Attribute, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type AttributeValueLookup interface {
	ListValueGraphByAttributeIDs(ctx context.Context, strategyID uuid.UUID, attributeIDs []uuid.UUID) (map[uuid.UUID][]attributevalue.AttributeValueView, error)
}

type Clock interface {
	clock.Clock
}

type AttributeSort = attributedomain.AttributeSort

const (
	AttributeSortCreatedAtDesc = attributedomain.AttributeSortCreatedAtDesc
	AttributeSortCreatedAtAsc  = attributedomain.AttributeSortCreatedAtAsc
)

type ListFilter struct {
	StrategyID uuid.UUID
	Search     string
	Page       int
	PageSize   int
	Sort       AttributeSort
}
