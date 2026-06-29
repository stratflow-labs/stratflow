package strategygrpc

import (
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	attribute "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attribute"
	attributevalue "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attributevalue"
	strategy "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/strategy"
	strategygraph "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/strategygraph"
)

type HandlerDependencies struct {
	Strategies      *strategy.Service
	Attributes      *attribute.Service
	AttributeValues *attributevalue.Service
	StrategyGraph   *strategygraph.Service
}

type Handler struct {
	strategyregistryv1.UnimplementedStrategyRegistryServiceServer

	strategies      *strategy.Service
	attributes      *attribute.Service
	attributeValues *attributevalue.Service
	strategyGraph   *strategygraph.Service
}

func NewHandler(deps HandlerDependencies) *Handler {
	return &Handler{
		strategies:      deps.Strategies,
		attributes:      deps.Attributes,
		attributeValues: deps.AttributeValues,
		strategyGraph:   deps.StrategyGraph,
	}
}
