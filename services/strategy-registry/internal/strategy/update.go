package strategy

import (
	"context"
	"errors"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	strategydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
)

type UpdateInput struct {
	Ref         strategydomain.EntityRef
	Slug        *string
	Name        *string
	Description *string
}

func (in UpdateInput) HasChanges() bool {
	return in.Slug != nil || in.Name != nil || in.Description != nil
}

func (s *WriteService) Update(ctx context.Context, input UpdateInput) (StrategyView, error) {
	if input.Ref.IsZero() {
		return StrategyView{}, apperr.Validation[strategydomain.Strategy](
			"update",
			[]apperr.FieldViolation{apperr.RefRequired("strategyRef", "strategy")},
		)
	}
	if !input.HasChanges() {
		return StrategyView{}, apperr.UpdateEmptyError[strategydomain.Strategy]()
	}

	var updated strategydomain.Strategy

	err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		current, err := s.resolveStrategy(txCtx, input.Ref)
		if err != nil {
			return fmt.Errorf("fetch strategy: %w", err)
		}

		domainCmd := strategydomain.UpdateStrategy{
			Slug:        input.Slug,
			Name:        input.Name,
			Description: input.Description,
		}

		if err := current.Update(domainCmd); err != nil {
			var validationErr strategydomain.ValidationError
			if errors.As(err, &validationErr) {
				return apperr.Validation[strategydomain.Strategy]("update", validationErr.Fields)
			}
			return err
		}

		current.UpdatedAt = s.clock.Now()

		updated, err = s.strategyRepo.Update(txCtx, &current)
		if err != nil {
			return fmt.Errorf("save strategy: %w", err)
		}

		return nil
	})

	if err != nil {
		return StrategyView{}, err
	}

	return strategyViewFromDomain(updated), nil
}
