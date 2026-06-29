package db

import (
	strategyregistrydbsqlc "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/postgres/sqlc/gen"
	attributedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
)

func attributeToDomain(model *strategyregistrydbsqlc.StrategyAttribute) attributedomain.Attribute {
	if model == nil {
		return attributedomain.Attribute{}
	}

	return attributedomain.Attribute{
		ID:          model.ID,
		StrategyID:  model.StrategyID,
		Slug:        model.Slug,
		Name:        model.Name,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func attributeToCreateParams(attribute *attributedomain.Attribute) strategyregistrydbsqlc.CreateStrategyAttributeParams {
	return strategyregistrydbsqlc.CreateStrategyAttributeParams{
		ID:          attribute.ID,
		StrategyID:  attribute.StrategyID,
		Slug:        attribute.Slug,
		Name:        attribute.Name,
		Description: attribute.Description,
		CreatedAt:   nonZeroTime(attribute.CreatedAt),
		UpdatedAt:   nonZeroTime(attribute.UpdatedAt, attribute.CreatedAt),
	}
}

func attributeToUpdateParams(attribute *attributedomain.Attribute) strategyregistrydbsqlc.UpdateStrategyAttributeParams {
	return strategyregistrydbsqlc.UpdateStrategyAttributeParams{
		ID:          attribute.ID,
		Slug:        attribute.Slug,
		Name:        attribute.Name,
		Description: attribute.Description,
		UpdatedAt:   nonZeroTime(attribute.UpdatedAt),
	}
}
