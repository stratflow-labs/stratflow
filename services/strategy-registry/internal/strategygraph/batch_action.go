package strategygraph

import (
	"context"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	"github.com/stratflow-labs/stratflow/internal/foundation/clock"
	tx "github.com/stratflow-labs/stratflow/internal/foundation/tx"
	attribute "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attribute"
	attributevalue "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attributevalue"
	registrydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

type Service struct {
	strategyRepo       StrategyRepository
	attributeRepo      AttributeRepository
	attributeLister    AttributeLister
	attributeValueRepo AttributeValueRepository
	txManager          tx.Manager
	clock              clock.Clock
}

func NewService(
	strategyRepo StrategyRepository,
	attributeRepo AttributeRepository,
	attributeLister AttributeLister,
	attributeValueRepo AttributeValueRepository,
	txManager tx.Manager,
	clock clock.Clock,
) *Service {
	return &Service{
		strategyRepo:       strategyRepo,
		attributeRepo:      attributeRepo,
		attributeLister:    attributeLister,
		attributeValueRepo: attributeValueRepo,
		txManager:          txManager,
		clock:              clock,
	}
}

func (s *Service) BatchAction(ctx context.Context, input BatchActionInput) (BatchActionOutput, error) {
	if input.StrategyID == uuid.Nil {
		return BatchActionOutput{}, apperr.NotFoundError[registrydomain.Strategy]()
	}
	if len(input.Actions) == 0 {
		return BatchActionOutput{}, apperr.BatchEmptyError[registrydomain.Strategy]("batchAction")
	}

	if _, err := s.loadGraph(ctx, input.StrategyID); err != nil {
		return BatchActionOutput{}, err
	}

	err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		for i := range input.Actions {
			if err := s.applyAction(txCtx, input.StrategyID, input.Actions[i]); err != nil {
				return fmt.Errorf("apply action[%d]: %w", i, err)
			}
		}
		return nil
	})
	if err != nil {
		return BatchActionOutput{}, fmt.Errorf("batch action strategy graph transaction: %w", err)
	}

	return s.loadGraph(ctx, input.StrategyID)
}

func (s *Service) applyAction(ctx context.Context, strategyID uuid.UUID, action Action) error {
	switch {
	case action.CreateAttribute != nil:
		return s.createAttribute(ctx, strategyID, *action.CreateAttribute)
	case action.UpdateAttribute != nil:
		return s.updateAttribute(ctx, strategyID, *action.UpdateAttribute)
	case action.DeleteAttribute != nil:
		return s.deleteAttribute(ctx, strategyID, *action.DeleteAttribute)
	case action.CreateValue != nil:
		return s.createValue(ctx, strategyID, *action.CreateValue)
	case action.UpdateValue != nil:
		return s.updateValue(ctx, strategyID, *action.UpdateValue)
	case action.DeleteValue != nil:
		return s.deleteValue(ctx, strategyID, *action.DeleteValue)
	case action.ReplaceRelations != nil:
		return s.replaceRelations(ctx, strategyID, *action.ReplaceRelations)
	default:
		return apperr.BatchEmptyError[registrydomain.Strategy]("batchAction")
	}
}

func (s *Service) createAttribute(ctx context.Context, strategyID uuid.UUID, input CreateInput) error {
	now := s.clock.Now()
	item := registrydomain.Attribute{
		StrategyID:  strategyID,
		Slug:        registrydomain.SanitizeString(input.Slug),
		Name:        registrydomain.SanitizeString(input.Name),
		Description: registrydomain.SanitizeString(input.Description),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if fields := registrydomain.ValidateAttribute(item.Slug, item.Name, item.Description); len(fields) > 0 {
		return apperr.Validation[registrydomain.Attribute]("batchCreate", fields)
	}
	_, err := s.attributeRepo.Create(ctx, &item)
	return err
}

func (s *Service) updateAttribute(ctx context.Context, strategyID uuid.UUID, input UpdateInput) error {
	current, err := s.resolveAttribute(ctx, strategyID, input.AttributeRef)
	if err != nil {
		return err
	}
	if current.StrategyID != strategyID {
		return apperr.NotFoundError[registrydomain.Attribute]()
	}
	if input.Slug == nil && input.Name == nil && input.Description == nil {
		return apperr.UpdateEmptyError[registrydomain.Attribute]()
	}

	if input.Slug != nil {
		current.Slug = registrydomain.SanitizeString(*input.Slug)
	}
	if input.Name != nil {
		current.Name = registrydomain.SanitizeString(*input.Name)
	}
	if input.Description != nil {
		current.Description = registrydomain.SanitizeString(*input.Description)
	}
	if fields := registrydomain.ValidateAttribute(current.Slug, current.Name, current.Description); len(fields) > 0 {
		return apperr.Validation[registrydomain.Attribute]("batchUpdate", fields)
	}
	current.UpdatedAt = s.clock.Now()
	_, err = s.attributeRepo.Update(ctx, &current)
	return err
}

func (s *Service) deleteAttribute(ctx context.Context, strategyID uuid.UUID, input DeleteInput) error {
	current, err := s.resolveAttribute(ctx, strategyID, input.AttributeRef)
	if err != nil {
		return err
	}
	if current.StrategyID != strategyID {
		return apperr.NotFoundError[registrydomain.Attribute]()
	}
	return s.attributeRepo.Delete(ctx, current.ID)
}

func (s *Service) createValue(ctx context.Context, strategyID uuid.UUID, input CreateValueInput) error {
	attrRef, err := s.resolveAttributeRef(ctx, strategyID, input.AttributeRef)
	if err != nil {
		return err
	}

	now := s.clock.Now()
	item := registrydomain.AttributeValue{
		StrategyID:  strategyID,
		AttributeID: attrRef.ID,
		Slug:        registrydomain.SanitizeString(input.Slug),
		Value:       registrydomain.SanitizeString(input.Value),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if fields := registrydomain.ValidateAttributeValue(item.Slug, item.Value); len(fields) > 0 {
		return apperr.Validation[registrydomain.AttributeValue]("batchCreate", fields)
	}
	_, err = s.attributeValueRepo.Create(ctx, &item)
	return err
}

func (s *Service) updateValue(ctx context.Context, strategyID uuid.UUID, input UpdateValueInput) error {
	attrRef, err := s.resolveAttributeRef(ctx, strategyID, input.AttributeRef)
	if err != nil {
		return err
	}
	current, err := s.resolveValue(ctx, strategyID, attrRef.ID, input.ValueRef)
	if err != nil {
		return err
	}
	if input.Slug == nil && input.Value == nil {
		return apperr.UpdateEmptyError[registrydomain.AttributeValue]()
	}

	if input.Slug != nil {
		current.Slug = registrydomain.SanitizeString(*input.Slug)
	}
	if input.Value != nil {
		current.Value = registrydomain.SanitizeString(*input.Value)
	}
	if fields := registrydomain.ValidateAttributeValue(current.Slug, current.Value); len(fields) > 0 {
		return apperr.Validation[registrydomain.AttributeValue]("batchUpdate", fields)
	}
	current.UpdatedAt = s.clock.Now()
	_, err = s.attributeValueRepo.Update(ctx, &current)
	return err
}

func (s *Service) deleteValue(ctx context.Context, strategyID uuid.UUID, input DeleteValueInput) error {
	attrRef, err := s.resolveAttributeRef(ctx, strategyID, input.AttributeRef)
	if err != nil {
		return err
	}
	current, err := s.resolveValue(ctx, strategyID, attrRef.ID, input.ValueRef)
	if err != nil {
		return err
	}
	return s.attributeValueRepo.Delete(ctx, current.ID)
}

func (s *Service) replaceRelations(ctx context.Context, strategyID uuid.UUID, input ReplaceRelationsInput) error {
	attrRef, err := s.resolveAttributeRef(ctx, strategyID, input.AttributeRef)
	if err != nil {
		return err
	}
	current, err := s.resolveValue(ctx, strategyID, attrRef.ID, input.ValueRef)
	if err != nil {
		return err
	}

	relations := make([]attributevalue.AttributeValueRelationInput, 0, len(input.Relations))
	seen := make(map[string]struct{}, len(input.Relations))
	for i := range input.Relations {
		targetAttribute, err := s.resolveAttributeRef(ctx, strategyID, input.Relations[i].AttributeRef)
		if err != nil {
			return fmt.Errorf("relation[%d] attribute: %w", i, err)
		}
		targetValue, err := s.resolveValue(ctx, strategyID, targetAttribute.ID, input.Relations[i].ValueRef)
		if err != nil {
			return fmt.Errorf("relation[%d] value: %w", i, err)
		}
		if targetAttribute.ID == current.AttributeID && targetValue.ID == current.ID {
			return apperr.SelfReferenceError[registrydomain.AttributeValue](
				"relation",
				"selfReference",
				"relation cannot reference itself",
			)
		}

		key := targetAttribute.ID.String() + "|" + targetValue.ID.String()
		if _, exists := seen[key]; exists {
			return apperr.DuplicateError[registrydomain.AttributeValue](
				"relation",
				"duplicate",
				"duplicate relations are not allowed",
			)
		}
		seen[key] = struct{}{}
		relations = append(relations, attributevalue.AttributeValueRelationInput{
			ToAttributeID: targetAttribute.ID,
			ToValueID:     targetValue.ID,
		})
	}

	return s.attributeValueRepo.ReplaceRelations(ctx, attributevalue.ReplaceAttributeValueRelationsInput{
		StrategyID:      strategyID,
		FromAttributeID: current.AttributeID,
		FromValueID:     current.ID,
		Relations:       relations,
	})
}

func (s *Service) resolveAttribute(ctx context.Context, strategyID uuid.UUID, ref EntityRef) (registrydomain.Attribute, error) {
	id, hasID, slug, err := normalizedEntityRef(ref)
	if err != nil {
		return registrydomain.Attribute{}, err
	}
	if hasID {
		item, err := s.attributeRepo.GetByID(ctx, id)
		if err != nil {
			return registrydomain.Attribute{}, err
		}
		if item.StrategyID != strategyID {
			return registrydomain.Attribute{}, apperr.NotFoundError[registrydomain.Attribute]()
		}
		if err := ensureOptionalSlugMatch(item.Slug, slug); err != nil {
			return registrydomain.Attribute{}, err
		}
		return item, nil
	}

	return s.attributeRepo.GetBySlug(ctx, strategyID, slug)
}

func (s *Service) resolveAttributeRef(ctx context.Context, strategyID uuid.UUID, ref EntityRef) (attributevalue.AttributeRef, error) {
	id, hasID, slug, err := normalizedEntityRef(ref)
	if err != nil {
		return attributevalue.AttributeRef{}, err
	}
	if hasID {
		item, err := s.attributeValueRepo.GetAttributeByID(ctx, id)
		if err != nil {
			return attributevalue.AttributeRef{}, err
		}
		if item.StrategyID != strategyID {
			return attributevalue.AttributeRef{}, apperr.NotFoundError[registrydomain.Attribute]()
		}
		if err := ensureOptionalSlugMatch(item.Slug, slug); err != nil {
			return attributevalue.AttributeRef{}, err
		}
		return item, nil
	}

	return s.attributeValueRepo.GetAttributeBySlug(ctx, strategyID, slug)
}

func (s *Service) resolveValue(ctx context.Context, strategyID, attributeID uuid.UUID, ref EntityRef) (registrydomain.AttributeValue, error) {
	id, hasID, slug, err := normalizedEntityRef(ref)
	if err != nil {
		return registrydomain.AttributeValue{}, err
	}
	if hasID {
		item, err := s.attributeValueRepo.GetByID(ctx, id)
		if err != nil {
			return registrydomain.AttributeValue{}, err
		}
		if item.StrategyID != strategyID || item.AttributeID != attributeID {
			return registrydomain.AttributeValue{}, apperr.Invalid[registrydomain.AttributeValue](
				"relation",
				"combinationNotFound",
				"relation id/slug combination not found",
			)
		}
		if err := ensureOptionalSlugMatch(item.Slug, slug); err != nil {
			return registrydomain.AttributeValue{}, err
		}
		return item, nil
	}

	item, err := s.attributeValueRepo.GetBySlug(ctx, attributeID, slug)
	if err != nil {
		return registrydomain.AttributeValue{}, err
	}
	if item.StrategyID != strategyID {
		return registrydomain.AttributeValue{}, apperr.Invalid[registrydomain.AttributeValue](
			"relation",
			"combinationNotFound",
			"relation id/slug combination not found",
		)
	}
	return item, nil
}

func normalizedEntityRef(ref EntityRef) (uuid.UUID, bool, string, error) {
	id, hasID, err := normalizeEntityRef(ref.ID, ref.Slug)
	if err != nil {
		return uuid.Nil, false, "", err
	}
	return id, hasID, sanitizeOptionalSlug(ref.Slug), nil
}

func normalizeEntityRef(idPtr *uuid.UUID, slugPtr *string) (uuid.UUID, bool, error) {
	if idPtr != nil && *idPtr != uuid.Nil {
		return *idPtr, true, nil
	}
	if slugPtr == nil || registrydomain.SanitizeString(*slugPtr) == "" {
		return uuid.Nil, false, apperr.Invalid[registrydomain.AttributeValue](
			"relation",
			"combinationNotFound",
			"relation id/slug combination not found",
		)
	}
	return uuid.Nil, false, nil
}

func sanitizeOptionalSlug(slugPtr *string) string {
	if slugPtr == nil {
		return ""
	}
	return registrydomain.SanitizeString(*slugPtr)
}

func ensureOptionalSlugMatch(actualSlug, expectedSlug string) error {
	if expectedSlug != "" && actualSlug != expectedSlug {
		return apperr.Invalid[registrydomain.AttributeValue](
			"relation",
			"combinationNotFound",
			"relation id/slug combination not found",
		)
	}
	return nil
}

func (s *Service) loadGraph(ctx context.Context, strategyID uuid.UUID) (BatchActionOutput, error) {
	item, err := s.strategyRepo.GetByID(ctx, strategyID)
	if err != nil {
		return BatchActionOutput{}, fmt.Errorf("load strategy graph strategy: %w", err)
	}

	attributes, err := s.loadGraphAttributes(ctx, strategyID)
	if err != nil {
		return BatchActionOutput{}, fmt.Errorf("load strategy graph attributes: %w", err)
	}

	return BatchActionOutput{
		Strategy:   strategyViewFromDomain(item),
		Attributes: attributes,
	}, nil
}

func (s *Service) loadGraphAttributes(ctx context.Context, strategyID uuid.UUID) ([]attribute.AttributeView, error) {
	var (
		page       = 1
		total      int64
		attributes []attribute.AttributeView
	)

	for {
		out, err := s.attributeLister.List(ctx, attribute.ListInput{
			StrategyID: strategyID,
			Page:       page,
			PageSize:   listAllPageSize,
			Sort:       string(attribute.AttributeSortCreatedAtDesc),
		})
		if err != nil {
			return nil, err
		}
		if page == 1 {
			total = out.Total
		}

		attributes = append(attributes, out.Attributes...)
		if len(out.Attributes) == 0 || int64(len(attributes)) >= total {
			return attributes, nil
		}

		page++
	}
}
