package strategygrpc

import (
	"strings"

	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	attribute "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attribute"
	attributevalue "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attributevalue"
	strategy "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/strategy"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapStrategy(item strategy.StrategyView) *strategyregistryv1.Strategy {
	return &strategyregistryv1.Strategy{
		Id:          item.ID.String(),
		Slug:        item.Slug,
		Name:        item.Name,
		Description: item.Description,
		CreatedAt:   timestamppb.New(item.CreatedAt),
		UpdatedAt:   timestamppb.New(item.UpdatedAt),
	}
}

func mapStrategyList(items []strategy.StrategyView) []*strategyregistryv1.Strategy {
	out := make([]*strategyregistryv1.Strategy, len(items))
	for i := range items {
		out[i] = mapStrategy(items[i])
	}
	return out
}

func mapStrategyWithParameters(item strategy.StrategyView, attributes []attribute.AttributeView) *strategyregistryv1.StrategyWithParameters {
	parameters := make([]*strategyregistryv1.AttributeWithValues, len(attributes))
	for i := range attributes {
		parameters[i] = mapAttribute(attributes[i])
	}

	return &strategyregistryv1.StrategyWithParameters{
		Id:          item.ID.String(),
		Slug:        item.Slug,
		Name:        item.Name,
		Description: item.Description,
		Parameters:  parameters,
		CreatedAt:   timestamppb.New(item.CreatedAt),
		UpdatedAt:   timestamppb.New(item.UpdatedAt),
	}
}

func mapAttribute(item attribute.AttributeView) *strategyregistryv1.AttributeWithValues {
	values := make([]*strategyregistryv1.AttributeValueInline, len(item.Values))
	for i := range item.Values {
		values[i] = mapAttributeValueInline(item.Values[i])
	}

	return &strategyregistryv1.AttributeWithValues{
		Id:          item.ID.String(),
		StrategyId:  item.StrategyID.String(),
		Slug:        item.Slug,
		Name:        item.Name,
		Description: item.Description,
		Values:      values,
		CreatedAt:   timestamppb.New(item.CreatedAt),
		UpdatedAt:   timestamppb.New(item.UpdatedAt),
	}
}

func mapAttributeList(items []attribute.AttributeView) []*strategyregistryv1.AttributeWithValues {
	out := make([]*strategyregistryv1.AttributeWithValues, len(items))
	for i := range items {
		out[i] = mapAttribute(items[i])
	}
	return out
}

func mapAttributeValue(item attributevalue.AttributeValueView) *strategyregistryv1.AttributeValue {
	return &strategyregistryv1.AttributeValue{
		Id:          item.ID.String(),
		AttributeId: item.AttributeID.String(),
		Slug:        item.Slug,
		Value:       item.Value,
		Relations:   mapAttributeValueRelations(item.Relations),
		CreatedAt:   timestamppb.New(item.CreatedAt),
		UpdatedAt:   timestamppb.New(item.UpdatedAt),
	}
}

func mapAttributeValueInline(item attribute.AttributeValueView) *strategyregistryv1.AttributeValueInline {
	return &strategyregistryv1.AttributeValueInline{
		Id:          item.ID.String(),
		AttributeId: item.AttributeID.String(),
		Slug:        item.Slug,
		Value:       item.Value,
		Relations:   mapAttributeRelations(item.Relations),
		CreatedAt:   timestamppb.New(item.CreatedAt),
		UpdatedAt:   timestamppb.New(item.UpdatedAt),
	}
}

func mapAttributeValueRelations(items []attributevalue.AttributeValueRelationView) []*strategyregistryv1.AttributeValueRelationInline {
	out := make([]*strategyregistryv1.AttributeValueRelationInline, len(items))
	for i := range items {
		out[i] = mapRelation(
			items[i].FromAttributeID.String(),
			items[i].FromValueID.String(),
			items[i].ToAttributeID.String(),
			items[i].ToValueID.String(),
			items[i].ToAttributeSlug,
			items[i].ToValueSlug,
		)
	}
	return out
}

func mapAttributeRelations(items []attribute.AttributeValueRelationView) []*strategyregistryv1.AttributeValueRelationInline {
	out := make([]*strategyregistryv1.AttributeValueRelationInline, len(items))
	for i := range items {
		out[i] = mapRelation(
			items[i].FromAttributeID.String(),
			items[i].FromValueID.String(),
			items[i].ToAttributeID.String(),
			items[i].ToValueID.String(),
			items[i].ToAttributeSlug,
			items[i].ToValueSlug,
		)
	}
	return out
}

func mapRelation(fromAttributeID, fromValueID, toAttributeID, toValueID, toAttributeSlug, toValueSlug string) *strategyregistryv1.AttributeValueRelationInline {
	out := &strategyregistryv1.AttributeValueRelationInline{
		FromAttributeId: fromAttributeID,
		FromValueId:     fromValueID,
		ToAttributeId:   toAttributeID,
		ToValueId:       toValueID,
	}
	if strings.TrimSpace(toAttributeSlug) != "" {
		out.ToAttributeSlug = stringPtr(toAttributeSlug)
	}
	if strings.TrimSpace(toValueSlug) != "" {
		out.ToValueSlug = stringPtr(toValueSlug)
	}
	return out
}

func stringPtr(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	out := strings.TrimSpace(value)
	return &out
}
