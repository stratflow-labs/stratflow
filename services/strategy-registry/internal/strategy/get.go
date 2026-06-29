package strategy

import (
	"context"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	strategydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

func (s *ReadService) Get(ctx context.Context, id uuid.UUID) (StrategyView, error) {
	if id == uuid.Nil {
		return StrategyView{}, apperr.NotFoundError[strategydomain.Strategy]()
	}

	strategy, err := s.strategyRepo.GetByID(ctx, id)
	if err != nil {
		return StrategyView{}, fmt.Errorf("get strategy: %w", err)
	}

	return strategyViewFromDomain(strategy), nil
}
