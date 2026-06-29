package strategygraph

import (
	"context"

	attribute "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attribute"
	attributevalue "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attributevalue"
	registrydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type StrategyRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (registrydomain.Strategy, error)
}

type AttributeLister interface {
	List(context.Context, attribute.ListInput) (attribute.ListOutput, error)
}

type AttributeRepository interface {
	Create(ctx context.Context, item *registrydomain.Attribute) (registrydomain.Attribute, error)
	GetByID(ctx context.Context, id uuid.UUID) (registrydomain.Attribute, error)
	GetBySlug(ctx context.Context, strategyID uuid.UUID, slug string) (registrydomain.Attribute, error)
	Update(ctx context.Context, item *registrydomain.Attribute) (registrydomain.Attribute, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type AttributeValueRepository interface {
	Create(ctx context.Context, item *registrydomain.AttributeValue) (registrydomain.AttributeValue, error)
	GetAttributeByID(ctx context.Context, id uuid.UUID) (attributevalue.AttributeRef, error)
	GetAttributeBySlug(ctx context.Context, strategyID uuid.UUID, slug string) (attributevalue.AttributeRef, error)
	GetByID(ctx context.Context, id uuid.UUID) (registrydomain.AttributeValue, error)
	GetBySlug(ctx context.Context, attributeID uuid.UUID, slug string) (registrydomain.AttributeValue, error)
	Update(ctx context.Context, item *registrydomain.AttributeValue) (registrydomain.AttributeValue, error)
	ReplaceRelations(ctx context.Context, input attributevalue.ReplaceAttributeValueRelationsInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}
