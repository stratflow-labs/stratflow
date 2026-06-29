package db

import (
	strategyregistrydbsqlc "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/postgres/sqlc/gen"
	strategydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
)

func strategyToDomain(model *strategyregistrydbsqlc.Strategy) strategydomain.Strategy {
	if model == nil {
		return strategydomain.Strategy{}
	}

	return strategydomain.Strategy{
		ID:          model.ID,
		Slug:        model.Slug,
		Name:        model.Name,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func strategyToCreateParams(s *strategydomain.Strategy) strategyregistrydbsqlc.CreateStrategyParams {
	return strategyregistrydbsqlc.CreateStrategyParams{
		ID:          s.ID,
		Slug:        s.Slug,
		Name:        s.Name,
		Description: s.Description,
		CreatedAt:   nonZeroTime(s.CreatedAt),
		UpdatedAt:   nonZeroTime(s.UpdatedAt, s.CreatedAt),
	}
}

func strategyToUpdateParams(strategy *strategydomain.Strategy) strategyregistrydbsqlc.UpdateStrategyParams {

	return strategyregistrydbsqlc.UpdateStrategyParams{
		ID:          strategy.ID,
		Slug:        strategy.Slug,
		Name:        strategy.Name,
		Description: strategy.Description,
		UpdatedAt:   nonZeroTime(strategy.UpdatedAt),
	}
}
