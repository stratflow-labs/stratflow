package strategy

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	strategydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
)

type CreateInput struct {
	Slug        string
	Name        string
	Description string
}

func (s *WriteService) Create(ctx context.Context, input CreateInput) (StrategyView, error) {
	now := s.clock.Now()
	id := uuid.New()
	item, err := strategydomain.NewStrategy(
		id,
		input.Slug,
		input.Name,
		input.Description,
		now,
	)
	if err != nil {
		var validationErr strategydomain.ValidationError
		if errors.As(err, &validationErr) {
			return StrategyView{}, apperr.Validation[strategydomain.Strategy]("create", validationErr.Fields)
		}
		return StrategyView{}, err
	}

	created, err := s.strategyRepo.Create(ctx, item)
	if err != nil {
		return StrategyView{}, fmt.Errorf("create strategy: %w", err)
	}

	return strategyViewFromDomain(created), nil
}
