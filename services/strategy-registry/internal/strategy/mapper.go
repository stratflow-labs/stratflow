package strategy

import (
	"time"

	strategydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type StrategyView struct {
	ID          uuid.UUID
	Slug        string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func strategyViewFromDomain(s strategydomain.Strategy) StrategyView {
	return StrategyView{
		ID:          s.ID,
		Slug:        s.Slug,
		Name:        s.Name,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func strategyViewsFromDomain(items []strategydomain.Strategy) []StrategyView {
	out := make([]StrategyView, len(items))
	for i := range items {
		out[i] = strategyViewFromDomain(items[i])
	}
	return out
}
