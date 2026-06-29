package strategygrpc

import (
	"context"
	"net/http"

	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	attribute "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attribute"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) ListAttributes(ctx context.Context, req *strategyregistryv1.ListAttributesRequest) (*strategyregistryv1.AttributesListResponse, error) {
	strategyID, err := parseStrategyID(req.GetStrategyRef())
	if err != nil {
		return nil, mapError(err)
	}

	out, err := h.attributes.List(ctx, attribute.ListInput{
		StrategyID: strategyID,
		Search:     req.GetSearch(),
		Page:       int(req.GetPage()),
		PageSize:   int(req.GetPageSize()),
		Sort:       req.GetSort(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &strategyregistryv1.AttributesListResponse{
		Message: "attributes listed",
		Data: &strategyregistryv1.AttributesListData{
			Items: mapAttributeList(out.Attributes),
			Total: out.Total,
		},
	}, nil
}

func (h *Handler) CreateAttribute(ctx context.Context, req *strategyregistryv1.CreateAttributeRequest) (*strategyregistryv1.AttributeResponse, error) {
	strategyID, err := parseStrategyID(req.GetStrategyRef())
	if err != nil {
		return nil, mapError(err)
	}

	out, err := h.attributes.Create(ctx, attribute.CreateInput{
		StrategyID:  strategyID,
		Slug:        req.GetSlug(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	setHTTPStatus(ctx, http.StatusCreated)
	return &strategyregistryv1.AttributeResponse{
		Message: "attribute created",
		Data:    mapAttribute(out),
	}, nil
}

func (h *Handler) GetAttribute(ctx context.Context, req *strategyregistryv1.GetAttributeRequest) (*strategyregistryv1.AttributeResponse, error) {
	strategyID, attributeID, err := parseStrategyAndAttributeIDs(req.GetStrategyRef(), req.GetAttributeRef())
	if err != nil {
		return nil, mapError(err)
	}

	out, err := h.attributes.Get(ctx, attribute.GetInput{
		ID:         attributeID,
		StrategyID: strategyID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &strategyregistryv1.AttributeResponse{
		Message: "attribute retrieved",
		Data:    mapAttribute(out),
	}, nil
}

func (h *Handler) UpdateAttribute(ctx context.Context, req *strategyregistryv1.UpdateAttributeRequest) (*strategyregistryv1.AttributeResponse, error) {
	strategyID, attributeID, err := parseStrategyAndAttributeIDs(req.GetStrategyRef(), req.GetAttributeRef())
	if err != nil {
		return nil, mapError(err)
	}

	patch := req.GetPatch()
	if patch == nil {
		patch = &strategyregistryv1.UpdateAttributePatch{}
	}

	out, err := h.attributes.Update(ctx, attribute.UpdateInput{
		ID:          attributeID,
		StrategyID:  strategyID,
		Slug:        patch.Slug,
		Name:        patch.Name,
		Description: patch.Description,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &strategyregistryv1.AttributeResponse{
		Message: "attribute updated",
		Data:    mapAttribute(out),
	}, nil
}

func (h *Handler) DeleteAttribute(ctx context.Context, req *strategyregistryv1.DeleteAttributeRequest) (*emptypb.Empty, error) {
	strategyID, attributeID, err := parseStrategyAndAttributeIDs(req.GetStrategyRef(), req.GetAttributeRef())
	if err != nil {
		return nil, mapError(err)
	}

	if err := h.attributes.Delete(ctx, attribute.DeleteInput{
		ID:         attributeID,
		StrategyID: strategyID,
	}); err != nil {
		return nil, mapError(err)
	}

	return noContent(ctx), nil
}
