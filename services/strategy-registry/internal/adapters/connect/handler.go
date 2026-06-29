package strategyconnect

import (
	"context"
	"net/http"

	"github.com/stratflow-labs/stratflow/internal/httpserver"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	strategyregistryv1connect "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1/strategyregistryv1connect"

	connect "connectrpc.com/connect"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pattern string
	handler http.Handler
}

func NewHandler(server strategyregistryv1.StrategyRegistryServiceServer) *Handler {
	pattern, handler := strategyregistryv1connect.NewStrategyRegistryServiceHandler(serviceAdapter{server: server})
	return &Handler{
		pattern: "/connect" + pattern,
		handler: http.StripPrefix("/connect", handler),
	}
}

func (h *Handler) RegisterRoutes(r httpserver.Router) {
	r.Mount(h.pattern, h.handler)
}

type serviceAdapter struct {
	server strategyregistryv1.StrategyRegistryServiceServer
}

func (a serviceAdapter) ListStrategies(ctx context.Context, req *connect.Request[strategyregistryv1.ListStrategiesRequest]) (*connect.Response[strategyregistryv1.StrategiesListResponse], error) {
	return unary(ctx, req, a.server.ListStrategies)
}

func (a serviceAdapter) CreateStrategy(ctx context.Context, req *connect.Request[strategyregistryv1.CreateStrategyRequest]) (*connect.Response[strategyregistryv1.StrategyResponse], error) {
	return unary(ctx, req, a.server.CreateStrategy)
}

func (a serviceAdapter) GetStrategy(ctx context.Context, req *connect.Request[strategyregistryv1.GetStrategyRequest]) (*connect.Response[strategyregistryv1.StrategyResponse], error) {
	return unary(ctx, req, a.server.GetStrategy)
}

func (a serviceAdapter) UpdateStrategy(ctx context.Context, req *connect.Request[strategyregistryv1.UpdateStrategyByRefRequest]) (*connect.Response[strategyregistryv1.StrategyResponse], error) {
	return unary(ctx, req, a.server.UpdateStrategy)
}

func (a serviceAdapter) DeleteStrategy(ctx context.Context, req *connect.Request[strategyregistryv1.DeleteStrategyRequest]) (*connect.Response[emptypb.Empty], error) {
	return unary(ctx, req, a.server.DeleteStrategy)
}

func (a serviceAdapter) BatchActionStrategyGraph(ctx context.Context, req *connect.Request[strategyregistryv1.BatchActionStrategyGraphRequest]) (*connect.Response[strategyregistryv1.StrategyGraphResponse], error) {
	return unary(ctx, req, a.server.BatchActionStrategyGraph)
}

func (a serviceAdapter) CloneStrategies(ctx context.Context, req *connect.Request[strategyregistryv1.CloneStrategiesRequest]) (*connect.Response[strategyregistryv1.StrategiesListResponse], error) {
	return unary(ctx, req, a.server.CloneStrategies)
}

func (a serviceAdapter) ListAttributes(ctx context.Context, req *connect.Request[strategyregistryv1.ListAttributesRequest]) (*connect.Response[strategyregistryv1.AttributesListResponse], error) {
	return unary(ctx, req, a.server.ListAttributes)
}

func (a serviceAdapter) CreateAttribute(ctx context.Context, req *connect.Request[strategyregistryv1.CreateAttributeRequest]) (*connect.Response[strategyregistryv1.AttributeResponse], error) {
	return unary(ctx, req, a.server.CreateAttribute)
}

func (a serviceAdapter) GetAttribute(ctx context.Context, req *connect.Request[strategyregistryv1.GetAttributeRequest]) (*connect.Response[strategyregistryv1.AttributeResponse], error) {
	return unary(ctx, req, a.server.GetAttribute)
}

func (a serviceAdapter) UpdateAttribute(ctx context.Context, req *connect.Request[strategyregistryv1.UpdateAttributeRequest]) (*connect.Response[strategyregistryv1.AttributeResponse], error) {
	return unary(ctx, req, a.server.UpdateAttribute)
}

func (a serviceAdapter) DeleteAttribute(ctx context.Context, req *connect.Request[strategyregistryv1.DeleteAttributeRequest]) (*connect.Response[emptypb.Empty], error) {
	return unary(ctx, req, a.server.DeleteAttribute)
}

func (a serviceAdapter) CreateAttributeValue(ctx context.Context, req *connect.Request[strategyregistryv1.CreateAttributeValueRequest]) (*connect.Response[strategyregistryv1.AttributeValueResponse], error) {
	return unary(ctx, req, a.server.CreateAttributeValue)
}

func (a serviceAdapter) UpdateAttributeValue(ctx context.Context, req *connect.Request[strategyregistryv1.UpdateAttributeValueRequest]) (*connect.Response[strategyregistryv1.AttributeValueResponse], error) {
	return unary(ctx, req, a.server.UpdateAttributeValue)
}

func (a serviceAdapter) DeleteAttributeValue(ctx context.Context, req *connect.Request[strategyregistryv1.DeleteAttributeValueRequest]) (*connect.Response[emptypb.Empty], error) {
	return unary(ctx, req, a.server.DeleteAttributeValue)
}

func unary[Req any, Resp any](ctx context.Context, req *connect.Request[Req], fn func(context.Context, *Req) (*Resp, error)) (*connect.Response[Resp], error) {
	msg, err := fn(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(msg), nil
}
