// strategy/ports.go
package strategy

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	"github.com/stratflow-labs/stratflow/internal/foundation/clock"
	strategydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
)

type StrategyRepository interface {
	Create(ctx context.Context, strategy *strategydomain.Strategy) (strategydomain.Strategy, error)
	GetByID(ctx context.Context, id uuid.UUID) (strategydomain.Strategy, error)
	GetBySlug(ctx context.Context, slug string) (strategydomain.Strategy, error)
	List(ctx context.Context, filter ListFilter) ([]strategydomain.Strategy, int64, error)
	Update(ctx context.Context, strategy *strategydomain.Strategy) (strategydomain.Strategy, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type StrategyCloneRepository interface {
	CloneBatch(ctx context.Context, items []CloneStrategySpec) ([]strategydomain.Strategy, error)
}

type Clock interface {
	clock.Clock
}

type ListFilter struct {
	Search   string
	Page     int
	PageSize int
	Sort     strategydomain.StrategySort
}

type CloneStrategySpec struct {
	SourceStrategyID uuid.UUID
	Slug             string
}

func (s *WriteService) resolveStrategy(ctx context.Context, ref strategydomain.EntityRef) (strategydomain.Strategy, error) {
	if ref.IsZero() {
		return strategydomain.Strategy{}, apperr.Validation[strategydomain.Strategy](
			"resolve",
			[]apperr.FieldViolation{apperr.RefRequired("strategyRef", "strategy")},
		)
	}

	if ref.HasID() {
		item, err := s.strategyRepo.GetByID(ctx, ref.ID)
		if err != nil {
			return strategydomain.Strategy{}, fmt.Errorf("get strategy by id: %w", err)
		}
		return item, nil
	}

	item, err := s.strategyRepo.GetBySlug(ctx, ref.NormalizedSlug())
	if err != nil {
		return strategydomain.Strategy{}, fmt.Errorf("get strategy by slug: %w", err)
	}
	return item, nil
}
