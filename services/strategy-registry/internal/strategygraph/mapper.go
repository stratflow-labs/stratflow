package strategygraph

import (
	registrydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
	strategy "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/strategy"
)

type StrategyView = strategy.StrategyView

func strategyViewFromDomain(item registrydomain.Strategy) StrategyView {
	return StrategyView{
		ID:          item.ID,
		Slug:        item.Slug,
		Name:        item.Name,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}
