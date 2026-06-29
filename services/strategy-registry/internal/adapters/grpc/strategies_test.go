package strategygrpc

import (
	"context"
	"testing"
	"time"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	attribute "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attribute"
	attributevalue "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attributevalue"
	registrydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
	strategygraph "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/strategygraph"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcGraphStrategyRepoStub struct {
	item registrydomain.Strategy
}

func (s *grpcGraphStrategyRepoStub) GetByID(_ context.Context, id uuid.UUID) (registrydomain.Strategy, error) {
	if id != s.item.ID {
		return registrydomain.Strategy{}, apperr.NotFoundError[registrydomain.Strategy]()
	}
	return s.item, nil
}

type grpcGraphAttributeRepoStub struct {
	items map[uuid.UUID]registrydomain.Attribute
}

func (s *grpcGraphAttributeRepoStub) Create(_ context.Context, item *registrydomain.Attribute) (registrydomain.Attribute, error) {
	item.ID = uuid.New()
	created := *item
	s.items[created.ID] = created
	return created, nil
}

func (s *grpcGraphAttributeRepoStub) GetByID(_ context.Context, id uuid.UUID) (registrydomain.Attribute, error) {
	item, ok := s.items[id]
	if !ok {
		return registrydomain.Attribute{}, apperr.NotFoundError[registrydomain.Attribute]()
	}
	return item, nil
}

func (s *grpcGraphAttributeRepoStub) GetBySlug(_ context.Context, strategyID uuid.UUID, slug string) (registrydomain.Attribute, error) {
	for _, item := range s.items {
		if item.StrategyID == strategyID && item.Slug == slug {
			return item, nil
		}
	}
	return registrydomain.Attribute{}, apperr.NotFoundError[registrydomain.Attribute]()
}

func (s *grpcGraphAttributeRepoStub) Update(_ context.Context, item *registrydomain.Attribute) (registrydomain.Attribute, error) {
	s.items[item.ID] = *item
	return *item, nil
}

func (s *grpcGraphAttributeRepoStub) Delete(_ context.Context, id uuid.UUID) error {
	delete(s.items, id)
	return nil
}

type grpcGraphAttributeListerStub struct {
	build func() []attribute.AttributeView
}

func (s *grpcGraphAttributeListerStub) List(context.Context, attribute.ListInput) (attribute.ListOutput, error) {
	attrs := s.build()
	return attribute.ListOutput{Attributes: attrs, Total: int64(len(attrs))}, nil
}

type grpcGraphAttributeValueRepoStub struct {
	attributes map[uuid.UUID]attributevalue.AttributeRef
	values     map[uuid.UUID]registrydomain.AttributeValue
}

func (s *grpcGraphAttributeValueRepoStub) Create(_ context.Context, item *registrydomain.AttributeValue) (registrydomain.AttributeValue, error) {
	item.ID = uuid.New()
	created := *item
	s.values[created.ID] = created
	return created, nil
}

func (s *grpcGraphAttributeValueRepoStub) GetAttributeByID(_ context.Context, id uuid.UUID) (attributevalue.AttributeRef, error) {
	item, ok := s.attributes[id]
	if !ok {
		return attributevalue.AttributeRef{}, apperr.NotFoundError[registrydomain.Attribute]()
	}
	return item, nil
}

func (s *grpcGraphAttributeValueRepoStub) GetAttributeBySlug(_ context.Context, strategyID uuid.UUID, slug string) (attributevalue.AttributeRef, error) {
	for _, item := range s.attributes {
		if item.StrategyID == strategyID && item.Slug == slug {
			return item, nil
		}
	}
	return attributevalue.AttributeRef{}, apperr.NotFoundError[registrydomain.Attribute]()
}

func (s *grpcGraphAttributeValueRepoStub) GetByID(_ context.Context, id uuid.UUID) (registrydomain.AttributeValue, error) {
	item, ok := s.values[id]
	if !ok {
		return registrydomain.AttributeValue{}, apperr.NotFoundError[registrydomain.AttributeValue]()
	}
	return item, nil
}

func (s *grpcGraphAttributeValueRepoStub) GetBySlug(_ context.Context, attributeID uuid.UUID, slug string) (registrydomain.AttributeValue, error) {
	for _, item := range s.values {
		if item.AttributeID == attributeID && item.Slug == slug {
			return item, nil
		}
	}
	return registrydomain.AttributeValue{}, apperr.NotFoundError[registrydomain.AttributeValue]()
}

func (s *grpcGraphAttributeValueRepoStub) Update(_ context.Context, item *registrydomain.AttributeValue) (registrydomain.AttributeValue, error) {
	s.values[item.ID] = *item
	return *item, nil
}

func (s *grpcGraphAttributeValueRepoStub) ReplaceRelations(_ context.Context, input attributevalue.ReplaceAttributeValueRelationsInput) error {
	current := s.values[input.FromValueID]
	next := make([]registrydomain.AttributeValueRelation, len(input.Relations))
	for i := range input.Relations {
		next[i] = registrydomain.AttributeValueRelation{
			FromAttributeID: input.FromAttributeID,
			FromValueID:     input.FromValueID,
			ToAttributeID:   input.Relations[i].ToAttributeID,
			ToValueID:       input.Relations[i].ToValueID,
		}
	}
	current.Relations = next
	s.values[input.FromValueID] = current
	return nil
}

func (s *grpcGraphAttributeValueRepoStub) Delete(_ context.Context, id uuid.UUID) error {
	delete(s.values, id)
	return nil
}

type grpcGraphClock struct {
	now time.Time
}

func (c grpcGraphClock) Now() time.Time { return c.now }

type grpcGraphTx struct{}

func (grpcGraphTx) WithinTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestHandler_BatchActionStrategyGraph_InvalidRefs(t *testing.T) {
	t.Parallel()

	handler := NewHandler(HandlerDependencies{})

	req := &strategyregistryv1.BatchActionStrategyGraphRequest{
		StrategyRef: uuid.New().String(),
		Actions: []*strategyregistryv1.StrategyGraphAction{
			{
				Action: &strategyregistryv1.StrategyGraphAction_DeleteValue{
					DeleteValue: &strategyregistryv1.GraphActionDeleteValue{
						AttributeRef: &strategyregistryv1.GraphAttributeRef{Id: stringPtr("bad")},
						ValueRef:     &strategyregistryv1.GraphValueRef{Id: stringPtr(uuid.New().String())},
					},
				},
			},
		},
	}

	_, err := handler.BatchActionStrategyGraph(context.Background(), req)
	require.Error(t, err)
	st := status.Convert(err)
	require.Equal(t, codes.InvalidArgument, st.Code())
	require.Equal(t, "invalid attributeRef", st.Message())
}

func TestHandler_BatchActionStrategyGraph_MapsValidRequest(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	strategyID := uuid.New()
	attributeID := uuid.New()
	valueID := uuid.New()

	attributeRepo := &grpcGraphAttributeRepoStub{
		items: map[uuid.UUID]registrydomain.Attribute{
			attributeID: {
				ID:          attributeID,
				StrategyID:  strategyID,
				Slug:        "timeframe",
				Name:        "Timeframe",
				Description: "tf",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}
	valueRepo := &grpcGraphAttributeValueRepoStub{
		attributes: map[uuid.UUID]attributevalue.AttributeRef{
			attributeID: {ID: attributeID, StrategyID: strategyID, Slug: "timeframe"},
		},
		values: map[uuid.UUID]registrydomain.AttributeValue{
			valueID: {
				ID:          valueID,
				StrategyID:  strategyID,
				AttributeID: attributeID,
				Slug:        "tf_15m",
				Value:       "15m",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}
	lister := &grpcGraphAttributeListerStub{
		build: func() []attribute.AttributeView {
			return []attribute.AttributeView{
				{
					ID:          attributeID,
					StrategyID:  strategyID,
					Slug:        "timeframe",
					Name:        "Timeframe",
					Description: "tf",
					Values: []attribute.AttributeValueView{
						{
							ID:          valueID,
							AttributeID: attributeID,
							Slug:        "tf_15m",
							Value:       "15m",
						},
					},
				},
			}
		},
	}

	handler := NewHandler(HandlerDependencies{
		StrategyGraph: strategygraph.NewService(
			&grpcGraphStrategyRepoStub{
				item: registrydomain.Strategy{
					ID:          strategyID,
					Slug:        "vb",
					Name:        "VB",
					Description: "graph",
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			attributeRepo,
			lister,
			valueRepo,
			grpcGraphTx{},
			grpcGraphClock{now: now},
		),
	})

	resp, err := handler.BatchActionStrategyGraph(context.Background(), &strategyregistryv1.BatchActionStrategyGraphRequest{
		StrategyRef: strategyID.String(),
		Actions: []*strategyregistryv1.StrategyGraphAction{
			{
				Action: &strategyregistryv1.StrategyGraphAction_UpdateAttribute{
					UpdateAttribute: &strategyregistryv1.GraphActionUpdateAttribute{
						AttributeRef: &strategyregistryv1.GraphAttributeRef{Slug: stringPtr("timeframe")},
						Name:         stringPtr("Execution Timeframe"),
					},
				},
			},
			{
				Action: &strategyregistryv1.StrategyGraphAction_CreateValue{
					CreateValue: &strategyregistryv1.GraphActionCreateValue{
						AttributeRef: &strategyregistryv1.GraphAttributeRef{Slug: stringPtr("timeframe")},
						Slug:         "tf_60m",
						Value:        "60m",
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "strategy graph batch action applied", resp.GetMessage())
	require.Equal(t, strategyID.String(), resp.GetData().GetId())
	require.Len(t, resp.GetData().GetParameters(), 1)
}
