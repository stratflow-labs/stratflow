package attributevalue

import (
	"context"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	attributeValuedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type CreateAttributeValueRelationInput struct {
	ToAttributeID   *uuid.UUID
	ToAttributeSlug *string
	ToValueID       *uuid.UUID
	ToValueSlug     *string
}

type CreateInput struct {
	StrategyID  uuid.UUID
	AttributeID uuid.UUID
	Slug        string
	Value       string
	Relations   []CreateAttributeValueRelationInput
}

func (s *Service) Create(ctx context.Context, input *CreateInput) (AttributeValueView, error) {
	if input == nil {
		return AttributeValueView{}, fmt.Errorf("validation failed: input is required")
	}
	if input.StrategyID == uuid.Nil {
		return AttributeValueView{}, apperr.NotFoundError[attributeValuedomain.Strategy]()
	}
	if input.AttributeID == uuid.Nil {
		return AttributeValueView{}, apperr.NotFoundError[attributeValuedomain.Attribute]()
	}

	now := s.clock.Now()
	id := uuid.New()
	attributeValue, err := attributeValuedomain.NewAttributeValue(
		id,
		input.StrategyID,
		input.AttributeID,
		input.Slug,
		input.Value,
		now,
	)
	if err != nil {
		return AttributeValueView{}, apperr.Invalid[attributeValuedomain.AttributeValue](
			"create",
			"validation",
			"attribute value validation failed",
			err,
		)
	}

	var created attributeValuedomain.AttributeValue
	if err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		var err error
		created, err = s.attributeValueRepo.Create(txCtx, attributeValue)
		if err != nil {
			return fmt.Errorf("create attributeValue: %w", err)
		}

		if len(input.Relations) == 0 {
			return nil
		}

		relations, err := s.resolveCreateRelations(txCtx, created, input.Relations)
		if err != nil {
			return fmt.Errorf("resolve relations: %w", err)
		}

		if err := s.attributeValueRepo.ReplaceRelations(txCtx, ReplaceAttributeValueRelationsInput{
			StrategyID:      created.StrategyID,
			FromAttributeID: created.AttributeID,
			FromValueID:     created.ID,
			Relations:       relations,
		}); err != nil {
			return fmt.Errorf("create attributeValue relations: %w", err)
		}

		return nil
	}); err != nil {
		return AttributeValueView{}, fmt.Errorf("create attributeValue transaction: %w", err)
	}

	view := attributeValueViewFromDomain(created)
	relationsByFromValueID, err := s.attributeValueRepo.ListRelationsByFromValueIDsForAttributeValues(ctx, created.StrategyID, []uuid.UUID{created.ID})
	if err != nil {
		return AttributeValueView{}, fmt.Errorf("load created attributeValue relations: %w", err)
	}
	if relations, ok := relationsByFromValueID[created.ID]; ok {
		view.Relations = relations
	}

	return view, nil
}

func (s *Service) resolveCreateRelations(
	ctx context.Context,
	created attributeValuedomain.AttributeValue,
	items []CreateAttributeValueRelationInput,
) ([]AttributeValueRelationInput, error) {
	resolved := make([]AttributeValueRelationInput, 0, len(items))
	seen := make(map[string]struct{}, len(items))

	for i := range items {
		param, err := s.resolveTargetAttribute(ctx, created.StrategyID, items[i])
		if err != nil {
			return nil, fmt.Errorf("relation[%d] target attribute: %w", i, err)
		}

		value, err := s.resolveTargetValue(ctx, created.StrategyID, param.ID, items[i])
		if err != nil {
			return nil, fmt.Errorf("relation[%d] target value: %w", i, err)
		}

		if param.ID == created.AttributeID && value.ID == created.ID {
			return nil, apperr.SelfReferenceError[attributeValuedomain.AttributeValue](
				"relation",
				"selfReference",
				"relation cannot reference itself",
			)
		}

		key := param.ID.String() + "|" + value.ID.String()
		if _, exists := seen[key]; exists {
			return nil, apperr.DuplicateError[attributeValuedomain.AttributeValue](
				"relation",
				"duplicate",
				"duplicate relations are not allowed",
			)
		}
		seen[key] = struct{}{}

		resolved = append(resolved, AttributeValueRelationInput{
			ToAttributeID: param.ID,
			ToValueID:     value.ID,
		})
	}

	return resolved, nil
}

func (s *Service) resolveTargetAttribute(
	ctx context.Context,
	strategyID uuid.UUID,
	input CreateAttributeValueRelationInput,
) (AttributeRef, error) {
	id, hasID := normalizeUUIDPtr(input.ToAttributeID)
	slug, hasSlug := normalizeStringPtr(input.ToAttributeSlug)

	if !hasID && !hasSlug {
		return AttributeRef{}, apperr.Validation[attributeValuedomain.AttributeValue](
			"relationAttributeRef",
			[]apperr.FieldViolation{apperr.RefRequired("toAttributeRef", "attribute")},
		)
	}

	var (
		byID   AttributeRef
		bySlug AttributeRef
		err    error
	)

	if hasID {
		byID, err = s.attributeValueRepo.GetAttributeByID(ctx, id)
		if err != nil {
			return AttributeRef{}, err
		}
		if byID.StrategyID != strategyID {
			return AttributeRef{}, apperr.NotFoundError[attributeValuedomain.Attribute]()
		}
	}

	if hasSlug {
		bySlug, err = s.attributeValueRepo.GetAttributeBySlug(ctx, strategyID, slug)
		if err != nil {
			return AttributeRef{}, err
		}
	}

	if hasID && hasSlug && byID.ID != bySlug.ID {
		return AttributeRef{}, apperr.Invalid[attributeValuedomain.AttributeValue](
			"relation",
			"combinationNotFound",
			"relation id/slug combination not found",
		)
	}
	if hasID {
		return byID, nil
	}
	return bySlug, nil
}

func (s *Service) resolveTargetValue(
	ctx context.Context,
	strategyID, attributeID uuid.UUID,
	input CreateAttributeValueRelationInput,
) (attributeValuedomain.AttributeValue, error) {
	id, hasID := normalizeUUIDPtr(input.ToValueID)
	slug, hasSlug := normalizeStringPtr(input.ToValueSlug)

	if !hasID && !hasSlug {
		return attributeValuedomain.AttributeValue{}, apperr.Validation[attributeValuedomain.AttributeValue](
			"relationValueRef",
			[]apperr.FieldViolation{apperr.RefRequired("toValueRef", "value")},
		)
	}

	var (
		byID   attributeValuedomain.AttributeValue
		bySlug attributeValuedomain.AttributeValue
		err    error
	)

	if hasID {
		byID, err = s.attributeValueRepo.GetByID(ctx, id)
		if err != nil {
			return attributeValuedomain.AttributeValue{}, err
		}
		if byID.StrategyID != strategyID || byID.AttributeID != attributeID {
			return attributeValuedomain.AttributeValue{}, apperr.Invalid[attributeValuedomain.AttributeValue](
				"relation",
				"combinationNotFound",
				"relation id/slug combination not found",
			)
		}
	}

	if hasSlug {
		bySlug, err = s.attributeValueRepo.GetBySlug(ctx, attributeID, slug)
		if err != nil {
			return attributeValuedomain.AttributeValue{}, err
		}
		if bySlug.StrategyID != strategyID {
			return attributeValuedomain.AttributeValue{}, apperr.Invalid[attributeValuedomain.AttributeValue](
				"relation",
				"combinationNotFound",
				"relation id/slug combination not found",
			)
		}
	}

	if hasID && hasSlug && byID.ID != bySlug.ID {
		return attributeValuedomain.AttributeValue{}, apperr.Invalid[attributeValuedomain.AttributeValue](
			"relation",
			"combinationNotFound",
			"relation id/slug combination not found",
		)
	}
	if hasID {
		return byID, nil
	}
	return bySlug, nil
}

func normalizeUUIDPtr(v *uuid.UUID) (uuid.UUID, bool) {
	if v == nil || *v == uuid.Nil {
		return uuid.Nil, false
	}
	return *v, true
}

func normalizeStringPtr(v *string) (string, bool) {
	if v == nil {
		return "", false
	}
	s := attributeValuedomain.SanitizeString(*v)
	if s == "" {
		return "", false
	}
	return s, true
}
