package db

import (
	strategyregistrydbsqlc "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/postgres/sqlc/gen"
	attributeValuedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

func attributeValueToDomain(model *strategyregistrydbsqlc.StrategyAttributeValue, strategyID uuid.UUID) attributeValuedomain.AttributeValue {
	if model == nil {
		return attributeValuedomain.AttributeValue{}
	}

	return attributeValuedomain.AttributeValue{
		ID:          model.ID,
		StrategyID:  strategyID,
		AttributeID: model.StrategyAttributeID,
		Slug:        model.Slug,
		Value:       model.Value,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func attributeValueToCreateParams(attributeValue *attributeValuedomain.AttributeValue) strategyregistrydbsqlc.CreateStrategyAttributeValueParams {
	return strategyregistrydbsqlc.CreateStrategyAttributeValueParams{
		ID:                  attributeValue.ID,
		StrategyAttributeID: attributeValue.AttributeID,
		Slug:                attributeValue.Slug,
		Value:               attributeValue.Value,
		CreatedAt:           nonZeroTime(attributeValue.CreatedAt),
		UpdatedAt:           nonZeroTime(attributeValue.UpdatedAt, attributeValue.CreatedAt),
	}
}

func attributeValueToUpdateParams(attributeValue *attributeValuedomain.AttributeValue) strategyregistrydbsqlc.UpdateStrategyAttributeValueParams {
	return strategyregistrydbsqlc.UpdateStrategyAttributeValueParams{
		ID:        attributeValue.ID,
		Slug:      attributeValue.Slug,
		Value:     attributeValue.Value,
		UpdatedAt: nonZeroTime(attributeValue.UpdatedAt),
	}
}
