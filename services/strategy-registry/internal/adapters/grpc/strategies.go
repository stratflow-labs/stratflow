package strategygrpc

import (
	"context"
	"net/http"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	registrydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
	strategy "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/strategy"
	strategygraph "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/strategygraph"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) ListStrategies(ctx context.Context, req *strategyregistryv1.ListStrategiesRequest) (*strategyregistryv1.StrategiesListResponse, error) {
	out, err := h.strategies.List(ctx, strategy.ListInput{
		Search:   req.GetSearch(),
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Sort:     req.GetSort(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &strategyregistryv1.StrategiesListResponse{
		Message: "strategies listed",
		Data: &strategyregistryv1.StrategiesListData{
			Items: mapStrategyList(out.Strategies),
			Total: out.Total,
		},
	}, nil
}

func (h *Handler) CreateStrategy(ctx context.Context, req *strategyregistryv1.CreateStrategyRequest) (*strategyregistryv1.StrategyResponse, error) {
	out, err := h.strategies.Create(ctx, strategy.CreateInput{
		Slug:        req.GetSlug(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	setHTTPStatus(ctx, http.StatusCreated)
	return &strategyregistryv1.StrategyResponse{
		Message: "strategy created",
		Data:    mapStrategy(out),
	}, nil
}

func (h *Handler) GetStrategy(ctx context.Context, req *strategyregistryv1.GetStrategyRequest) (*strategyregistryv1.StrategyResponse, error) {
	strategyID, err := parseStrategyID(req.GetStrategyRef())
	if err != nil {
		return nil, mapError(err)
	}

	out, err := h.strategies.Get(ctx, strategyID)
	if err != nil {
		return nil, mapError(err)
	}

	return &strategyregistryv1.StrategyResponse{
		Message: "strategy retrieved",
		Data:    mapStrategy(out),
	}, nil
}

func (h *Handler) UpdateStrategy(ctx context.Context, req *strategyregistryv1.UpdateStrategyByRefRequest) (*strategyregistryv1.StrategyResponse, error) {
	strategyID, err := parseStrategyID(req.GetStrategyRef())
	if err != nil {
		return nil, mapError(err)
	}

	patch := req.GetPatch()
	if patch == nil {
		patch = &strategyregistryv1.UpdateStrategyRequest{}
	}

	out, err := h.strategies.Update(ctx, strategy.UpdateInput{
		Ref:         registrydomain.RefByID(strategyID),
		Slug:        patch.Slug,
		Name:        patch.Name,
		Description: patch.Description,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &strategyregistryv1.StrategyResponse{
		Message: "strategy updated",
		Data:    mapStrategy(out),
	}, nil
}

func (h *Handler) DeleteStrategy(ctx context.Context, req *strategyregistryv1.DeleteStrategyRequest) (*emptypb.Empty, error) {
	strategyID, err := parseStrategyID(req.GetStrategyRef())
	if err != nil {
		return nil, mapError(err)
	}

	if err := h.strategies.Delete(ctx, strategyID); err != nil {
		return nil, mapError(err)
	}

	return noContent(ctx), nil
}

func (h *Handler) BatchActionStrategyGraph(ctx context.Context, req *strategyregistryv1.BatchActionStrategyGraphRequest) (*strategyregistryv1.StrategyGraphResponse, error) {
	strategyID, err := parseStrategyID(req.GetStrategyRef())
	if err != nil {
		return nil, mapError(err)
	}

	actions := make([]strategygraph.Action, len(req.GetActions()))
	for i := range req.GetActions() {
		action, err := mapStrategyGraphAction(req.GetActions()[i])
		if err != nil {
			return nil, err
		}
		actions[i] = action
	}

	out, err := h.strategyGraph.BatchAction(ctx, strategygraph.BatchActionInput{
		StrategyID: strategyID,
		Actions:    actions,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &strategyregistryv1.StrategyGraphResponse{
		Message: "strategy graph batch action applied",
		Data:    mapStrategyWithParameters(out.Strategy, out.Attributes),
	}, nil
}

func mapStrategyGraphAction(req *strategyregistryv1.StrategyGraphAction) (strategygraph.Action, error) {
	switch action := req.GetAction().(type) {
	case *strategyregistryv1.StrategyGraphAction_CreateAttribute:
		return strategygraph.Action{
			CreateAttribute: &strategygraph.CreateInput{
				Slug:        action.CreateAttribute.GetSlug(),
				Name:        action.CreateAttribute.GetName(),
				Description: action.CreateAttribute.GetDescription(),
			},
		}, nil
	case *strategyregistryv1.StrategyGraphAction_UpdateAttribute:
		ref, err := mapGraphAttributeRef(action.UpdateAttribute.GetAttributeRef())
		if err != nil {
			return strategygraph.Action{}, invalidArgument("request.invalid", "invalid attributeRef")
		}
		return strategygraph.Action{
			UpdateAttribute: &strategygraph.UpdateInput{
				AttributeRef: ref,
				Slug:         action.UpdateAttribute.Slug,
				Name:         action.UpdateAttribute.Name,
				Description:  action.UpdateAttribute.Description,
			},
		}, nil
	case *strategyregistryv1.StrategyGraphAction_DeleteAttribute:
		ref, err := mapGraphAttributeRef(action.DeleteAttribute.GetAttributeRef())
		if err != nil {
			return strategygraph.Action{}, invalidArgument("request.invalid", "invalid attributeRef")
		}
		return strategygraph.Action{
			DeleteAttribute: &strategygraph.DeleteInput{AttributeRef: ref},
		}, nil
	case *strategyregistryv1.StrategyGraphAction_CreateValue:
		attrRef, err := mapGraphAttributeRef(action.CreateValue.GetAttributeRef())
		if err != nil {
			return strategygraph.Action{}, invalidArgument("request.invalid", "invalid attributeRef")
		}
		return strategygraph.Action{
			CreateValue: &strategygraph.CreateValueInput{
				AttributeRef: attrRef,
				Slug:         action.CreateValue.GetSlug(),
				Value:        action.CreateValue.GetValue(),
			},
		}, nil
	case *strategyregistryv1.StrategyGraphAction_UpdateValue:
		attrRef, err := mapGraphAttributeRef(action.UpdateValue.GetAttributeRef())
		if err != nil {
			return strategygraph.Action{}, invalidArgument("request.invalid", "invalid attributeRef")
		}
		valueRef, err := mapGraphValueRef(action.UpdateValue.GetValueRef())
		if err != nil {
			return strategygraph.Action{}, invalidArgument("request.invalid", "invalid valueRef")
		}
		return strategygraph.Action{
			UpdateValue: &strategygraph.UpdateValueInput{
				AttributeRef: attrRef,
				ValueRef:     valueRef,
				Slug:         action.UpdateValue.Slug,
				Value:        action.UpdateValue.Value,
			},
		}, nil
	case *strategyregistryv1.StrategyGraphAction_DeleteValue:
		attrRef, err := mapGraphAttributeRef(action.DeleteValue.GetAttributeRef())
		if err != nil {
			return strategygraph.Action{}, invalidArgument("request.invalid", "invalid attributeRef")
		}
		valueRef, err := mapGraphValueRef(action.DeleteValue.GetValueRef())
		if err != nil {
			return strategygraph.Action{}, invalidArgument("request.invalid", "invalid valueRef")
		}
		return strategygraph.Action{
			DeleteValue: &strategygraph.DeleteValueInput{
				AttributeRef: attrRef,
				ValueRef:     valueRef,
			},
		}, nil
	case *strategyregistryv1.StrategyGraphAction_ReplaceRelations:
		attrRef, err := mapGraphAttributeRef(action.ReplaceRelations.GetAttributeRef())
		if err != nil {
			return strategygraph.Action{}, invalidArgument("request.invalid", "invalid attributeRef")
		}
		valueRef, err := mapGraphValueRef(action.ReplaceRelations.GetValueRef())
		if err != nil {
			return strategygraph.Action{}, invalidArgument("request.invalid", "invalid valueRef")
		}
		relations := make([]strategygraph.RelationTargetInput, len(action.ReplaceRelations.GetRelations()))
		for i := range action.ReplaceRelations.GetRelations() {
			targetAttrRef, err := mapGraphAttributeRef(action.ReplaceRelations.GetRelations()[i].GetAttributeRef())
			if err != nil {
				return strategygraph.Action{}, invalidArgument("request.invalid", "invalid relation attributeRef")
			}
			targetValueRef, err := mapGraphValueRef(action.ReplaceRelations.GetRelations()[i].GetValueRef())
			if err != nil {
				return strategygraph.Action{}, invalidArgument("request.invalid", "invalid relation valueRef")
			}
			relations[i] = strategygraph.RelationTargetInput{
				AttributeRef: targetAttrRef,
				ValueRef:     targetValueRef,
			}
		}
		return strategygraph.Action{
			ReplaceRelations: &strategygraph.ReplaceRelationsInput{
				AttributeRef: attrRef,
				ValueRef:     valueRef,
				Relations:    relations,
			},
		}, nil
	default:
		return strategygraph.Action{}, invalidArgument("request.invalid", "graph action is required")
	}
}

func mapGraphAttributeRef(ref *strategyregistryv1.GraphAttributeRef) (strategygraph.EntityRef, error) {
	if ref == nil {
		return strategygraph.EntityRef{}, apperr.NotFoundError[registrydomain.Attribute]()
	}
	mapped, err := parseOptionalGraphRef(ref.Id, ref.Slug, apperr.NotFoundError[registrydomain.Attribute]())
	if err != nil {
		return strategygraph.EntityRef{}, err
	}
	return strategygraph.EntityRef{
		ID:   mapped.ID,
		Slug: mapped.Slug,
	}, nil
}

func mapGraphValueRef(ref *strategyregistryv1.GraphValueRef) (strategygraph.EntityRef, error) {
	if ref == nil {
		return strategygraph.EntityRef{}, apperr.NotFoundError[registrydomain.AttributeValue]()
	}
	mapped, err := parseOptionalGraphRef(ref.Id, ref.Slug, apperr.NotFoundError[registrydomain.AttributeValue]())
	if err != nil {
		return strategygraph.EntityRef{}, err
	}
	return strategygraph.EntityRef{
		ID:   mapped.ID,
		Slug: mapped.Slug,
	}, nil
}

func (h *Handler) CloneStrategies(ctx context.Context, req *strategyregistryv1.CloneStrategiesRequest) (*strategyregistryv1.StrategiesListResponse, error) {
	items := make([]strategy.CloneStrategyItemInput, len(req.GetItems()))
	for i := range req.GetItems() {
		sourceID, err := parseUUID(req.GetItems()[i].GetSourceStrategyId(), apperr.NotFoundError[registrydomain.Strategy]())
		if err != nil {
			return nil, invalidArgument("request.invalid", "invalid sourceStrategyId")
		}
		items[i] = strategy.CloneStrategyItemInput{
			SourceStrategyID: sourceID,
			Slug:             req.GetItems()[i].GetSlug(),
		}
	}

	out, err := h.strategies.Clone(ctx, strategy.CloneInput{Items: items})
	if err != nil {
		return nil, mapError(err)
	}

	setHTTPStatus(ctx, http.StatusCreated)
	return &strategyregistryv1.StrategiesListResponse{
		Message: "strategies cloned",
		Data: &strategyregistryv1.StrategiesListData{
			Items: mapStrategyList(out.Strategies),
			Total: out.Total,
		},
	}, nil
}
