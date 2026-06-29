package attributevalue

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	attributeValuedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
)

type UpdateInput struct {
	ID          uuid.UUID
	StrategyID  uuid.UUID
	AttributeID uuid.UUID
	Slug        *string
	Value       *string
	Relations   *[]AttributeValueRelationInput
}

func (in UpdateInput) HasChanges() bool {
	return in.Slug != nil || in.Value != nil || in.Relations != nil
}

func (s *Service) Update(ctx context.Context, input UpdateInput) (AttributeValueView, error) {
	if input.ID == uuid.Nil {
		return AttributeValueView{}, apperr.NotFoundError[attributeValuedomain.AttributeValue]()
	}

	if input.StrategyID == uuid.Nil {
		return AttributeValueView{}, apperr.NotFoundError[attributeValuedomain.Strategy]()
	}
	if input.AttributeID == uuid.Nil {
		return AttributeValueView{}, apperr.NotFoundError[attributeValuedomain.Attribute]()
	}

	if !input.HasChanges() {
		return AttributeValueView{}, apperr.UpdateEmptyError[attributeValuedomain.AttributeValue]()
	}

	var updated attributeValuedomain.AttributeValue
	err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		current, err := s.attributeValueRepo.GetByID(txCtx, input.ID)
		if err != nil {
			return fmt.Errorf("load attributeValue: %w", err)
		}
		if current.StrategyID != input.StrategyID || current.AttributeID != input.AttributeID {
			return apperr.NotFoundError[attributeValuedomain.AttributeValue]()
		}

		if input.Slug != nil {
			current.Slug = attributeValuedomain.SanitizeString(*input.Slug)
		}
		if input.Value != nil {
			current.Value = attributeValuedomain.SanitizeString(*input.Value)
		}

		if fields := attributeValuedomain.ValidateAttributeValue(current.Slug, current.Value); len(fields) > 0 {
			return apperr.Validation[attributeValuedomain.AttributeValue]("update", fields)
		}

		var relations []attributeValuedomain.AttributeValueRelation
		if input.Relations != nil {
			relations = make([]attributeValuedomain.AttributeValueRelation, 0, len(*input.Relations))
			for i := range *input.Relations {
				relation, err := attributeValuedomain.NewAttributeValueRelation(
					current.AttributeID,
					current.ID,
					(*input.Relations)[i].ToAttributeID,
					(*input.Relations)[i].ToValueID,
				)
				if err != nil {
					return fmt.Errorf("relation[%d]: %w", i, err)
				}
				relations = append(relations, relation)
			}

			relations = attributeValuedomain.DedupeAttributeValueRelations(relations)
			if err := attributeValuedomain.ValidateAttributeValueRelations(relations); err != nil {
				return fmt.Errorf("relations validation: %w", err)
			}
		}

		current.UpdatedAt = s.clock.Now()

		updated, err = s.attributeValueRepo.Update(txCtx, &current)
		if err != nil {
			return fmt.Errorf("update attributeValue: %w", err)
		}

		if input.Relations != nil {
			if err := s.attributeValueRepo.ReplaceRelations(txCtx, ReplaceAttributeValueRelationsInput{
				StrategyID:      input.StrategyID,
				FromAttributeID: current.AttributeID,
				FromValueID:     current.ID,
				Relations:       relationInputsFromDomain(relations),
			}); err != nil {
				return fmt.Errorf("replace attributeValue relations: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return AttributeValueView{}, fmt.Errorf("update attributeValue transaction: %w", err)
	}

	view := attributeValueViewFromDomain(updated)
	relationsByFromValueID, err := s.attributeValueRepo.ListRelationsByFromValueIDsForAttributeValues(
		ctx,
		input.StrategyID,
		[]uuid.UUID{updated.ID},
	)
	if err != nil {
		return AttributeValueView{}, fmt.Errorf("load updated attributeValue relations: %w", err)
	}
	if relations, ok := relationsByFromValueID[updated.ID]; ok {
		view.Relations = relations
	}

	return view, nil
}
