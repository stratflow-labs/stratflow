package strategy

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	strategydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
)

type CloneInput struct {
	Items []CloneStrategyItemInput
}

type CloneStrategyItemInput struct {
	SourceStrategyID uuid.UUID
	Slug             string
}

type CloneOutput struct {
	Strategies []StrategyView
	Total      int64
}

func (s *WriteService) Clone(ctx context.Context, input CloneInput) (CloneOutput, error) {
	if len(input.Items) == 0 {
		return CloneOutput{}, apperr.CloneEmptyError[strategydomain.Strategy]()
	}

	specs := make([]CloneStrategySpec, 0, len(input.Items))
	seenSlugs := make(map[string]struct{}, len(input.Items))
	for i := range input.Items {
		if input.Items[i].SourceStrategyID == uuid.Nil {
			return CloneOutput{}, apperr.NotFoundError[strategydomain.Strategy]()
		}
		slug := strategydomain.SanitizeString(input.Items[i].Slug)
		if slug == "" {
			return CloneOutput{}, apperr.Validation[strategydomain.Strategy]("clone", []apperr.FieldViolation{apperr.Required("slug")})
		}
		if _, exists := seenSlugs[slug]; exists {
			return CloneOutput{}, apperr.DuplicateFieldInRequestError[strategydomain.Strategy]("clone", "slug")
		}
		seenSlugs[slug] = struct{}{}

		specs = append(specs, CloneStrategySpec{
			SourceStrategyID: input.Items[i].SourceStrategyID,
			Slug:             slug,
		})
	}

	var cloned []strategydomain.Strategy
	if err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		var err error
		cloned, err = s.strategyCloneRepo.CloneBatch(txCtx, specs)
		if err != nil {
			return fmt.Errorf("clone strategies: %w", err)
		}
		return nil
	}); err != nil {
		return CloneOutput{}, fmt.Errorf("clone strategies transaction: %w", err)
	}

	views := strategyViewsFromDomain(cloned)
	return CloneOutput{
		Strategies: views,
		Total:      int64(len(views)),
	}, nil
}
