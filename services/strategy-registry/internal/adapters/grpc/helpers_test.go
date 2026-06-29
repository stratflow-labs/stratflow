package strategygrpc

import (
	"errors"
	"testing"
	"time"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	attribute "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attribute"
	attributevalue "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attributevalue"
	"github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
	strategy "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/strategy"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestParseRefHelpers(t *testing.T) {
	t.Parallel()

	strategyID := uuid.New()
	attributeID := uuid.New()
	valueID := uuid.New()
	rawValueID := " " + valueID.String() + " "

	gotStrategyID, gotAttributeID, gotValueID, err := parseRefs(strategyID.String(), attributeID.String(), rawValueID)
	require.NoError(t, err)
	require.Equal(t, strategyID, gotStrategyID)
	require.Equal(t, attributeID, gotAttributeID)
	require.Equal(t, valueID, gotValueID)

	optional, ok, err := parseOptionalUUID(&rawValueID)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, optional)
	require.Equal(t, valueID, *optional)

	blank := "   "
	optional, ok, err = parseOptionalUUID(&blank)
	require.NoError(t, err)
	require.False(t, ok)
	require.Nil(t, optional)
}

func TestParseRefHelpers_ReturnAppErrors(t *testing.T) {
	t.Parallel()

	_, err := parseStrategyID("bad")
	require.ErrorIs(t, err, apperr.NotFoundError[domain.Strategy]())

	_, err = parseAttributeID(uuid.Nil.String())
	require.ErrorIs(t, err, apperr.NotFoundError[domain.Attribute]())

	_, err = parseValueID("bad")
	require.ErrorIs(t, err, apperr.NotFoundError[domain.AttributeValue]())

	raw := "bad"
	_, _, err = parseOptionalUUID(&raw)
	require.EqualError(t, err, "invalid UUID length: 3")

	_, _, err = parseStrategyAndAttributeIDs(uuid.New().String(), "bad")
	require.ErrorIs(t, err, apperr.NotFoundError[domain.Attribute]())
}

func TestMapError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		err     error
		code    codes.Code
		reason  string
		message string
	}{
		{
			name:    "wrapped strategy not found",
			err:     errors.Join(errors.New("wrapped"), apperr.NotFoundError[domain.Strategy]()),
			code:    codes.NotFound,
			reason:  "strategy.notFound",
			message: "strategy not found",
		},
		{
			name:    "relation duplicate",
			err:     apperr.DuplicateError[domain.AttributeValue]("relation", "duplicate", "duplicate relations are not allowed"),
			code:    codes.InvalidArgument,
			reason:  "attributeValue.relationDuplicate",
			message: "duplicate relations are not allowed",
		},
		{
			name:    "graph batch empty",
			err:     apperr.BatchEmptyError[domain.Strategy]("batchAction"),
			code:    codes.InvalidArgument,
			reason:  "strategyGraph.batchActionEmpty",
			message: "strategy batch actions are required",
		},
		{
			name:    "fallback internal",
			err:     errors.New("boom"),
			code:    codes.Internal,
			reason:  "strategyRegistry.internal",
			message: "Internal Server Error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			st := status.Convert(mapError(tc.err))
			require.Equal(t, tc.code, st.Code())
			require.Equal(t, tc.message, st.Message())

			require.Len(t, st.Details(), 1)
			info, ok := st.Details()[0].(*errdetails.ErrorInfo)
			require.True(t, ok)
			require.Equal(t, tc.reason, info.Reason)
		})
	}
}

func TestMappingHelpers(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 18, 13, 0, 0, 0, time.UTC)
	strategyID := uuid.New()
	attributeID := uuid.New()
	valueID := uuid.New()
	relatedAttributeID := uuid.New()
	relatedValueID := uuid.New()

	mappedStrategy := mapStrategyWithParameters(
		strategy.StrategyView{
			ID:          strategyID,
			Slug:        "trend-following",
			Name:        "Trend Following",
			Description: "core strategy",
			CreatedAt:   now,
			UpdatedAt:   now.Add(time.Minute),
		},
		[]attribute.AttributeView{
			{
				ID:          attributeID,
				StrategyID:  strategyID,
				Slug:        "timeframe",
				Name:        "Timeframe",
				Description: "candles",
				CreatedAt:   now,
				UpdatedAt:   now.Add(time.Minute),
				Values: []attribute.AttributeValueView{
					{
						ID:          valueID,
						AttributeID: attributeID,
						Slug:        "1h",
						Value:       "1H",
						CreatedAt:   now,
						UpdatedAt:   now.Add(2 * time.Minute),
						Relations: []attribute.AttributeValueRelationView{
							{
								FromAttributeID: attributeID,
								FromValueID:     valueID,
								ToAttributeID:   relatedAttributeID,
								ToValueID:       relatedValueID,
								ToAttributeSlug: " related-attr ",
								ToValueSlug:     " related-value ",
							},
						},
					},
				},
			},
		},
	)

	require.Equal(t, strategyID.String(), mappedStrategy.Id)
	require.Len(t, mappedStrategy.Parameters, 1)
	require.Len(t, mappedStrategy.Parameters[0].Values, 1)
	relation := mappedStrategy.Parameters[0].Values[0].Relations[0]
	require.Equal(t, "related-attr", relation.GetToAttributeSlug())
	require.Equal(t, "related-value", relation.GetToValueSlug())

	relationNoSlugs := mapRelation(
		attributeID.String(),
		valueID.String(),
		relatedAttributeID.String(),
		relatedValueID.String(),
		"   ",
		"",
	)
	require.Nil(t, relationNoSlugs.ToAttributeSlug)
	require.Nil(t, relationNoSlugs.ToValueSlug)
	require.Equal(t, "trimmed", *stringPtr("  trimmed  "))
	require.Nil(t, stringPtr(" \n "))

	mappedValue := mapAttributeValue(attributevalue.AttributeValueView{
		ID:          valueID,
		AttributeID: attributeID,
		Slug:        "1h",
		Value:       "1H",
		CreatedAt:   now,
		UpdatedAt:   now.Add(time.Minute),
	})
	require.Equal(t, attributeID.String(), mappedValue.AttributeId)
	require.Empty(t, mappedValue.Relations)
}

func TestInvalidArgument(t *testing.T) {
	t.Parallel()

	st := status.Convert(invalidArgument("request.invalid", "request is invalid"))
	require.Equal(t, codes.InvalidArgument, st.Code())
	require.Equal(t, "request is invalid", st.Message())
}
