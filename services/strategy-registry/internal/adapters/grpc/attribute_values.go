package strategygrpc

import (
	"context"
	"net/http"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	attributevalue "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attributevalue"
	registrydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) CreateAttributeValue(ctx context.Context, req *strategyregistryv1.CreateAttributeValueRequest) (*strategyregistryv1.AttributeValueResponse, error) {
	strategyID, attributeID, err := parseStrategyAndAttributeIDs(req.GetStrategyRef(), req.GetAttributeRef())
	if err != nil {
		return nil, mapError(err)
	}

	relations := make([]attributevalue.CreateAttributeValueRelationInput, len(req.GetRelations()))
	for i := range req.GetRelations() {
		toAttributeID, _, err := parseOptionalUUID(req.GetRelations()[i].ToAttributeId)
		if err != nil {
			return nil, invalidArgument("request.invalid", "invalid toAttributeId")
		}
		toValueID, _, err := parseOptionalUUID(req.GetRelations()[i].ToValueId)
		if err != nil {
			return nil, invalidArgument("request.invalid", "invalid toValueId")
		}

		relations[i] = attributevalue.CreateAttributeValueRelationInput{
			ToAttributeID:   toAttributeID,
			ToAttributeSlug: req.GetRelations()[i].ToAttributeSlug,
			ToValueID:       toValueID,
			ToValueSlug:     req.GetRelations()[i].ToValueSlug,
		}
	}

	out, err := h.attributeValues.Create(ctx, &attributevalue.CreateInput{
		StrategyID:  strategyID,
		AttributeID: attributeID,
		Slug:        req.GetSlug(),
		Value:       req.GetValue(),
		Relations:   relations,
	})
	if err != nil {
		return nil, mapError(err)
	}

	setHTTPStatus(ctx, http.StatusCreated)
	return &strategyregistryv1.AttributeValueResponse{
		Message: "attribute value created",
		Data:    mapAttributeValue(out),
	}, nil
}

func (h *Handler) UpdateAttributeValue(ctx context.Context, req *strategyregistryv1.UpdateAttributeValueRequest) (*strategyregistryv1.AttributeValueResponse, error) {
	strategyID, attributeID, valueID, err := parseRefs(req.GetStrategyRef(), req.GetAttributeRef(), req.GetValueRef())
	if err != nil {
		return nil, mapError(err)
	}

	patch := req.GetPatch()
	if patch == nil {
		patch = &strategyregistryv1.UpdateAttributeValuePatch{}
	}

	var relations *[]attributevalue.AttributeValueRelationInput
	if patch.Relations != nil {
		items := make([]attributevalue.AttributeValueRelationInput, len(patch.Relations))
		for i := range patch.Relations {
			toAttributeID, err := parseUUID(patch.Relations[i].GetToAttributeId(), apperr.NotFoundError[registrydomain.Attribute]())
			if err != nil {
				return nil, invalidArgument("request.invalid", "invalid toAttributeId")
			}
			toValueID, err := parseUUID(patch.Relations[i].GetToValueId(), apperr.NotFoundError[registrydomain.AttributeValue]())
			if err != nil {
				return nil, invalidArgument("request.invalid", "invalid toValueId")
			}

			items[i] = attributevalue.AttributeValueRelationInput{
				ToAttributeID: toAttributeID,
				ToValueID:     toValueID,
			}
		}
		relations = &items
	}

	out, err := h.attributeValues.Update(ctx, attributevalue.UpdateInput{
		ID:          valueID,
		StrategyID:  strategyID,
		AttributeID: attributeID,
		Slug:        patch.Slug,
		Value:       patch.Value,
		Relations:   relations,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &strategyregistryv1.AttributeValueResponse{
		Message: "attribute value updated",
		Data:    mapAttributeValue(out),
	}, nil
}

func (h *Handler) DeleteAttributeValue(ctx context.Context, req *strategyregistryv1.DeleteAttributeValueRequest) (*emptypb.Empty, error) {
	strategyID, attributeID, valueID, err := parseRefs(req.GetStrategyRef(), req.GetAttributeRef(), req.GetValueRef())
	if err != nil {
		return nil, mapError(err)
	}

	if err := h.attributeValues.Delete(ctx, attributevalue.DeleteInput{
		ID:          valueID,
		StrategyID:  strategyID,
		AttributeID: attributeID,
	}); err != nil {
		return nil, mapError(err)
	}

	return noContent(ctx), nil
}
