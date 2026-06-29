package testkit

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	e2ecommon "github.com/stratflow-labs/stratflow/services/identity/tests/e2e/testkit"
	strategyregistryv1 "github.com/stratflow-labs/stratflow/services/strategy-registry/gen/go/proto/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StrategyGraph struct {
	StrategyID    string
	StrategySlug  string
	AttributeID   string
	AttributeSlug string
	ValueID       string
	ValueSlug     string
}

func CreateStrategyGraph(t *testing.T, cfg *axiom.Config, token string) StrategyGraph {
	t.Helper()

	client := ClientFromFixture(cfg).StrategyRegistry()
	suffix := e2ecommon.UniqueSuffix("strategy")
	strategySlug := "e2e-strategy-" + suffix
	attributeSlug := "e2e-param-" + suffix
	valueSlug := "e2e-value-" + suffix

	strategyResp, err := client.CreateStrategy(Context(cfg, token), &strategyregistryv1.CreateStrategyRequest{
		Slug:        strategySlug,
		Name:        "E2E Strategy",
		Description: "created by strategy-registry grpc e2e",
	})
	require.NoError(t, err)
	require.NotNil(t, strategyResp.GetData())
	strategyID := strategyResp.GetData().GetId()

	t.Cleanup(func() {
		_, cleanupErr := client.DeleteStrategy(Context(cfg, token), &strategyregistryv1.DeleteStrategyRequest{StrategyRef: strategyID})
		if cleanupErr != nil && status.Code(cleanupErr) != codes.NotFound {
			t.Fatalf("cleanup strategy %s: %v", strategyID, cleanupErr)
		}
	})

	attributeResp, err := client.CreateAttribute(Context(cfg, token), &strategyregistryv1.CreateAttributeRequest{
		StrategyRef: strategyID,
		Slug:        attributeSlug,
		Name:        "E2E Attribute",
		Description: "created by strategy-registry grpc e2e",
	})
	require.NoError(t, err)
	require.NotNil(t, attributeResp.GetData())
	attributeID := attributeResp.GetData().GetId()

	valueResp, err := client.CreateAttributeValue(Context(cfg, token), &strategyregistryv1.CreateAttributeValueRequest{
		StrategyRef:  strategyID,
		AttributeRef: attributeID,
		Slug:         valueSlug,
		Value:        "42",
		Relations:    []*strategyregistryv1.CreateAttributeValueRelationInput{},
	})
	require.NoError(t, err)
	require.NotNil(t, valueResp.GetData())

	return StrategyGraph{
		StrategyID:    strategyID,
		StrategySlug:  strategySlug,
		AttributeID:   attributeID,
		AttributeSlug: attributeSlug,
		ValueID:       valueResp.GetData().GetId(),
		ValueSlug:     valueSlug,
	}
}
