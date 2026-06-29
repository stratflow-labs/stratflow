package strategygraph

import (
	"context"
	"testing"
	"time"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	attribute "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attribute"
	attributevalue "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attributevalue"
	registrydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type batchActionClock struct {
	now time.Time
}

func (c batchActionClock) Now() time.Time { return c.now }

type batchActionStrategyRepoStub struct {
	item registrydomain.Strategy
}

func (s *batchActionStrategyRepoStub) GetByID(_ context.Context, id uuid.UUID) (registrydomain.Strategy, error) {
	if id != s.item.ID {
		return registrydomain.Strategy{}, apperr.NotFoundError[registrydomain.Strategy]()
	}
	return s.item, nil
}

type batchActionAttributeListerStub struct {
	build func() []attribute.AttributeView
}

func (s *batchActionAttributeListerStub) List(context.Context, attribute.ListInput) (attribute.ListOutput, error) {
	return attribute.ListOutput{
		Attributes: s.build(),
		Total:      int64(len(s.build())),
	}, nil
}

type batchActionAttributeRepoStub struct {
	items map[uuid.UUID]registrydomain.Attribute
}

func (s *batchActionAttributeRepoStub) Create(_ context.Context, item *registrydomain.Attribute) (registrydomain.Attribute, error) {
	item.ID = uuid.New()
	created := *item
	s.items[created.ID] = created
	return created, nil
}

func (s *batchActionAttributeRepoStub) GetByID(_ context.Context, id uuid.UUID) (registrydomain.Attribute, error) {
	item, ok := s.items[id]
	if !ok {
		return registrydomain.Attribute{}, apperr.NotFoundError[registrydomain.Attribute]()
	}
	return item, nil
}

func (s *batchActionAttributeRepoStub) GetBySlug(_ context.Context, strategyID uuid.UUID, slug string) (registrydomain.Attribute, error) {
	for _, item := range s.items {
		if item.StrategyID == strategyID && item.Slug == slug {
			return item, nil
		}
	}
	return registrydomain.Attribute{}, apperr.NotFoundError[registrydomain.Attribute]()
}

func (s *batchActionAttributeRepoStub) Update(_ context.Context, item *registrydomain.Attribute) (registrydomain.Attribute, error) {
	if _, ok := s.items[item.ID]; !ok {
		return registrydomain.Attribute{}, apperr.NotFoundError[registrydomain.Attribute]()
	}
	updated := *item
	s.items[updated.ID] = updated
	return updated, nil
}

func (s *batchActionAttributeRepoStub) Delete(_ context.Context, id uuid.UUID) error {
	delete(s.items, id)
	return nil
}

type batchActionAttributeValueRepoStub struct {
	attributes map[uuid.UUID]attributevalue.AttributeRef
	values     map[uuid.UUID]registrydomain.AttributeValue
}

func (s *batchActionAttributeValueRepoStub) Create(_ context.Context, item *registrydomain.AttributeValue) (registrydomain.AttributeValue, error) {
	item.ID = uuid.New()
	created := *item
	s.values[created.ID] = created
	return created, nil
}

func (s *batchActionAttributeValueRepoStub) GetAttributeByID(_ context.Context, id uuid.UUID) (attributevalue.AttributeRef, error) {
	item, ok := s.attributes[id]
	if !ok {
		return attributevalue.AttributeRef{}, apperr.NotFoundError[registrydomain.Attribute]()
	}
	return item, nil
}

func (s *batchActionAttributeValueRepoStub) GetAttributeBySlug(_ context.Context, strategyID uuid.UUID, slug string) (attributevalue.AttributeRef, error) {
	for _, item := range s.attributes {
		if item.StrategyID == strategyID && item.Slug == slug {
			return item, nil
		}
	}
	return attributevalue.AttributeRef{}, apperr.NotFoundError[registrydomain.Attribute]()
}

func (s *batchActionAttributeValueRepoStub) GetByID(_ context.Context, id uuid.UUID) (registrydomain.AttributeValue, error) {
	item, ok := s.values[id]
	if !ok {
		return registrydomain.AttributeValue{}, apperr.NotFoundError[registrydomain.AttributeValue]()
	}
	return item, nil
}

func (s *batchActionAttributeValueRepoStub) GetBySlug(_ context.Context, attributeID uuid.UUID, slug string) (registrydomain.AttributeValue, error) {
	for _, item := range s.values {
		if item.AttributeID == attributeID && item.Slug == slug {
			return item, nil
		}
	}
	return registrydomain.AttributeValue{}, apperr.NotFoundError[registrydomain.AttributeValue]()
}

func (s *batchActionAttributeValueRepoStub) Update(_ context.Context, item *registrydomain.AttributeValue) (registrydomain.AttributeValue, error) {
	if _, ok := s.values[item.ID]; !ok {
		return registrydomain.AttributeValue{}, apperr.NotFoundError[registrydomain.AttributeValue]()
	}
	updated := *item
	s.values[updated.ID] = updated
	return updated, nil
}

func (s *batchActionAttributeValueRepoStub) ReplaceRelations(_ context.Context, input attributevalue.ReplaceAttributeValueRelationsInput) error {
	value, ok := s.values[input.FromValueID]
	if !ok {
		return apperr.NotFoundError[registrydomain.AttributeValue]()
	}
	next := make([]registrydomain.AttributeValueRelation, len(input.Relations))
	for i := range input.Relations {
		next[i] = registrydomain.AttributeValueRelation{
			FromAttributeID: input.FromAttributeID,
			FromValueID:     input.FromValueID,
			ToAttributeID:   input.Relations[i].ToAttributeID,
			ToValueID:       input.Relations[i].ToValueID,
		}
	}
	value.Relations = next
	s.values[input.FromValueID] = value
	return nil
}

func (s *batchActionAttributeValueRepoStub) Delete(_ context.Context, id uuid.UUID) error {
	delete(s.values, id)
	return nil
}

type batchActionTxManagerStub struct{}

func (batchActionTxManagerStub) WithinTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestBatchActionStrategyGraphUseCase_Execute_AppliesOrderedActionsWithSlugRefs(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	strategyID := uuid.New()
	riskAttributeID := uuid.New()
	timeframeAttributeID := uuid.New()
	riskValueID := uuid.New()
	tf15ValueID := uuid.New()
	tf30ValueID := uuid.New()

	attributeRepo := &batchActionAttributeRepoStub{
		items: map[uuid.UUID]registrydomain.Attribute{
			riskAttributeID: {
				ID:          riskAttributeID,
				StrategyID:  strategyID,
				Slug:        "risk_level",
				Name:        "Risk Level",
				Description: "risk",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			timeframeAttributeID: {
				ID:          timeframeAttributeID,
				StrategyID:  strategyID,
				Slug:        "timeframe",
				Name:        "Timeframe",
				Description: "tf",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}
	attributeValueRepo := &batchActionAttributeValueRepoStub{
		attributes: map[uuid.UUID]attributevalue.AttributeRef{
			riskAttributeID:      {ID: riskAttributeID, StrategyID: strategyID, Slug: "risk_level"},
			timeframeAttributeID: {ID: timeframeAttributeID, StrategyID: strategyID, Slug: "timeframe"},
		},
		values: map[uuid.UUID]registrydomain.AttributeValue{
			riskValueID: {
				ID:          riskValueID,
				StrategyID:  strategyID,
				AttributeID: riskAttributeID,
				Slug:        "very_high",
				Value:       "Very High",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			tf15ValueID: {
				ID:          tf15ValueID,
				StrategyID:  strategyID,
				AttributeID: timeframeAttributeID,
				Slug:        "tf_15m",
				Value:       "15m",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			tf30ValueID: {
				ID:          tf30ValueID,
				StrategyID:  strategyID,
				AttributeID: timeframeAttributeID,
				Slug:        "tf_30m",
				Value:       "30m",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}

	lister := &batchActionAttributeListerStub{
		build: func() []attribute.AttributeView {
			out := make([]attribute.AttributeView, 0, len(attributeRepo.items))
			for _, attrItem := range attributeRepo.items {
				view := attribute.AttributeView{
					ID:          attrItem.ID,
					StrategyID:  attrItem.StrategyID,
					Slug:        attrItem.Slug,
					Name:        attrItem.Name,
					Description: attrItem.Description,
					Values:      []attribute.AttributeValueView{},
				}
				for _, valueItem := range attributeValueRepo.values {
					if valueItem.AttributeID != attrItem.ID {
						continue
					}
					relations := make([]attribute.AttributeValueRelationView, len(valueItem.Relations))
					for i := range valueItem.Relations {
						targetAttr := attributeValueRepo.attributes[valueItem.Relations[i].ToAttributeID]
						targetValue := attributeValueRepo.values[valueItem.Relations[i].ToValueID]
						relations[i] = attribute.AttributeValueRelationView{
							FromAttributeID: valueItem.Relations[i].FromAttributeID,
							FromValueID:     valueItem.Relations[i].FromValueID,
							ToAttributeID:   valueItem.Relations[i].ToAttributeID,
							ToValueID:       valueItem.Relations[i].ToValueID,
							ToAttributeSlug: targetAttr.Slug,
							ToValueSlug:     targetValue.Slug,
						}
					}
					view.Values = append(view.Values, attribute.AttributeValueView{
						ID:          valueItem.ID,
						AttributeID: valueItem.AttributeID,
						Slug:        valueItem.Slug,
						Value:       valueItem.Value,
						Relations:   relations,
					})
				}
				out = append(out, view)
			}
			return out
		},
	}

	svc := NewService(
		&batchActionStrategyRepoStub{
			item: registrydomain.Strategy{
				ID:          strategyID,
				Slug:        "volatility-breakout",
				Name:        "Volatility Breakout",
				Description: "strategy graph",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		attributeRepo,
		lister,
		attributeValueRepo,
		batchActionTxManagerStub{},
		batchActionClock{now: now},
	)

	out, err := svc.BatchAction(context.Background(), BatchActionStrategyGraphInput{
		StrategyID: strategyID,
		Actions: []StrategyGraphAction{
			{
				CreateValue: &BatchCreateValueInput{
					AttributeRef: StrategyGraphEntityRef{Slug: new("timeframe")},
					Slug:         "tf_60m",
					Value:        "60m",
				},
			},
			{
				UpdateAttribute: &BatchUpdateInput{
					AttributeRef: StrategyGraphEntityRef{ID: &riskAttributeID},
					Name:         new("Risk Range"),
					Description:  new("risk range"),
				},
			},
			{
				DeleteValue: &BatchDeleteValueInput{
					AttributeRef: StrategyGraphEntityRef{Slug: new("timeframe")},
					ValueRef:     StrategyGraphEntityRef{Slug: new("tf_30m")},
				},
			},
			{
				ReplaceRelations: &BatchReplaceRelationsInput{
					AttributeRef: StrategyGraphEntityRef{Slug: new("risk_level")},
					ValueRef:     StrategyGraphEntityRef{Slug: new("very_high")},
					Relations: []StrategyGraphRelationTargetInput{
						{
							AttributeRef: StrategyGraphEntityRef{Slug: new("timeframe")},
							ValueRef:     StrategyGraphEntityRef{Slug: new("tf_15m")},
						},
						{
							AttributeRef: StrategyGraphEntityRef{Slug: new("timeframe")},
							ValueRef:     StrategyGraphEntityRef{Slug: new("tf_60m")},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, strategyID, out.Strategy.ID)

	riskAttr, err := attributeRepo.GetByID(context.Background(), riskAttributeID)
	require.NoError(t, err)
	require.Equal(t, "Risk Range", riskAttr.Name)

	_, err = attributeValueRepo.GetBySlug(context.Background(), timeframeAttributeID, "tf_30m")
	require.ErrorIs(t, err, apperr.NotFoundError[registrydomain.AttributeValue]())
	tf60, err := attributeValueRepo.GetBySlug(context.Background(), timeframeAttributeID, "tf_60m")
	require.NoError(t, err)

	riskValue, err := attributeValueRepo.GetByID(context.Background(), riskValueID)
	require.NoError(t, err)
	require.Len(t, riskValue.Relations, 2)
	require.Equal(t, tf15ValueID, riskValue.Relations[0].ToValueID)
	require.Equal(t, tf60.ID, riskValue.Relations[1].ToValueID)
}

func TestBatchActionStrategyGraphUseCase_Execute_RejectsDuplicateRelations(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	strategyID := uuid.New()
	attributeID := uuid.New()
	valueID := uuid.New()
	targetAttributeID := uuid.New()
	targetValueID := uuid.New()

	svc := NewService(
		&batchActionStrategyRepoStub{item: registrydomain.Strategy{ID: strategyID}},
		&batchActionAttributeRepoStub{
			items: map[uuid.UUID]registrydomain.Attribute{
				attributeID:       {ID: attributeID, StrategyID: strategyID, Slug: "source"},
				targetAttributeID: {ID: targetAttributeID, StrategyID: strategyID, Slug: "target"},
			},
		},
		&batchActionAttributeListerStub{build: func() []attribute.AttributeView { return nil }},
		&batchActionAttributeValueRepoStub{
			attributes: map[uuid.UUID]attributevalue.AttributeRef{
				attributeID:       {ID: attributeID, StrategyID: strategyID, Slug: "source"},
				targetAttributeID: {ID: targetAttributeID, StrategyID: strategyID, Slug: "target"},
			},
			values: map[uuid.UUID]registrydomain.AttributeValue{
				valueID:       {ID: valueID, StrategyID: strategyID, AttributeID: attributeID, Slug: "a"},
				targetValueID: {ID: targetValueID, StrategyID: strategyID, AttributeID: targetAttributeID, Slug: "b"},
			},
		},
		batchActionTxManagerStub{},
		batchActionClock{now: now},
	)

	_, err := svc.BatchAction(context.Background(), BatchActionStrategyGraphInput{
		StrategyID: strategyID,
		Actions: []StrategyGraphAction{
			{
				ReplaceRelations: &BatchReplaceRelationsInput{
					AttributeRef: StrategyGraphEntityRef{ID: &attributeID},
					ValueRef:     StrategyGraphEntityRef{ID: &valueID},
					Relations: []StrategyGraphRelationTargetInput{
						{
							AttributeRef: StrategyGraphEntityRef{ID: &targetAttributeID},
							ValueRef:     StrategyGraphEntityRef{ID: &targetValueID},
						},
						{
							AttributeRef: StrategyGraphEntityRef{Slug: new("target")},
							ValueRef:     StrategyGraphEntityRef{Slug: new("b")},
						},
					},
				},
			},
		},
	})
	require.ErrorIs(t, err, apperr.DuplicateError[registrydomain.AttributeValue]("relation", "duplicate", "duplicate relations are not allowed"))
}
