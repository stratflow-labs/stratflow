package attributevalue

import (
	"context"

	"github.com/stratflow-labs/stratflow/internal/foundation/clock"
	attributeValuedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type AttributeValueRepository interface {
	Create(ctx context.Context, attributeValue *attributeValuedomain.AttributeValue) (attributeValuedomain.AttributeValue, error)
	GetAttributeByID(ctx context.Context, id uuid.UUID) (AttributeRef, error)
	GetAttributeBySlug(ctx context.Context, strategyID uuid.UUID, slug string) (AttributeRef, error)
	GetByID(ctx context.Context, id uuid.UUID) (attributeValuedomain.AttributeValue, error)
	GetBySlug(ctx context.Context, attributeID uuid.UUID, slug string) (attributeValuedomain.AttributeValue, error)
	List(ctx context.Context, filter ListFilter) ([]attributeValuedomain.AttributeValue, int64, error)
	ListRelationsByFromValueIDsForAttributeValues(ctx context.Context, strategyID uuid.UUID, fromValueIDs []uuid.UUID) (map[uuid.UUID][]AttributeValueRelationView, error)
	Update(ctx context.Context, attributeValue *attributeValuedomain.AttributeValue) (attributeValuedomain.AttributeValue, error)
	ReplaceRelations(ctx context.Context, input ReplaceAttributeValueRelationsInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type AttributeRef struct {
	ID         uuid.UUID
	StrategyID uuid.UUID
	Slug       string
}

type Clock interface {
	clock.Clock
}

type AttributeValueSort = attributeValuedomain.AttributeValueSort

const (
	AttributeValueSortCreatedAtDesc = attributeValuedomain.AttributeValueSortCreatedAtDesc
	AttributeValueSortCreatedAtAsc  = attributeValuedomain.AttributeValueSortCreatedAtAsc
)

type ListFilter struct {
	StrategyID  uuid.UUID
	AttributeID uuid.UUID
	Search      string
	Page        int
	PageSize    int
	Sort        AttributeValueSort
}

type AttributeValueRelationInput struct {
	ToAttributeID uuid.UUID
	ToValueID     uuid.UUID
}

type ReplaceAttributeValueRelationsInput struct {
	StrategyID      uuid.UUID
	FromAttributeID uuid.UUID
	FromValueID     uuid.UUID
	Relations       []AttributeValueRelationInput
}
