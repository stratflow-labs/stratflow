package strategygrpc

import (
	"context"
	"net/http"
	"strings"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	"github.com/stratflow-labs/stratflow/internal/grpcserver"
	registrydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type optionalGraphRef struct {
	ID   *uuid.UUID
	Slug *string
}

func noContent(ctx context.Context) *emptypb.Empty {
	setHTTPStatus(ctx, http.StatusNoContent)
	return &emptypb.Empty{}
}

func setHTTPStatus(ctx context.Context, statusCode int) {
	grpcserver.SetHTTPStatus(ctx, statusCode)
}

func parseStrategyID(raw string) (uuid.UUID, error) {
	return parseUUID(raw, apperr.NotFoundError[registrydomain.Strategy]())
}

func parseAttributeID(raw string) (uuid.UUID, error) {
	return parseUUID(raw, apperr.NotFoundError[registrydomain.Attribute]())
}

func parseValueID(raw string) (uuid.UUID, error) {
	return parseUUID(raw, apperr.NotFoundError[registrydomain.AttributeValue]())
}

func parseUUID(raw string, err error) (uuid.UUID, error) {
	id, parseErr := uuid.Parse(strings.TrimSpace(raw))
	if parseErr != nil || id == uuid.Nil {
		return uuid.Nil, err
	}
	return id, nil
}

func parseOptionalUUID(raw *string) (_ *uuid.UUID, ok bool, err error) {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return nil, false, nil
	}

	id, err := uuid.Parse(strings.TrimSpace(*raw))
	if err != nil || id == uuid.Nil {
		return nil, false, err
	}
	return &id, true, nil
}

func parseOptionalGraphRef(idRaw, slugRaw *string, notFoundErr error) (optionalGraphRef, error) {
	id, ok, err := parseOptionalUUID(idRaw)
	if err != nil {
		return optionalGraphRef{}, notFoundErr
	}

	var slug *string
	if slugRaw != nil {
		normalized := strings.TrimSpace(*slugRaw)
		if normalized != "" {
			slug = &normalized
		}
	}

	if !ok && slug == nil {
		return optionalGraphRef{}, notFoundErr
	}

	return optionalGraphRef{
		ID:   id,
		Slug: slug,
	}, nil
}

func parseStrategyAndAttributeIDs(strategyRef, attributeRef string) (uuid.UUID, uuid.UUID, error) {
	strategyID, err := parseStrategyID(strategyRef)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	attributeID, err := parseAttributeID(attributeRef)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	return strategyID, attributeID, nil
}

func parseRefs(strategyRef, attributeRef, valueRef string) (uuid.UUID, uuid.UUID, uuid.UUID, error) {
	strategyID, attributeID, err := parseStrategyAndAttributeIDs(strategyRef, attributeRef)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, err
	}
	valueID, err := parseValueID(valueRef)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, err
	}
	return strategyID, attributeID, valueID, nil
}
