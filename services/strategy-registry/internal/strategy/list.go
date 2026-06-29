package strategy

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	"github.com/stratflow-labs/stratflow/internal/foundation/pagination"
	strategydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

type ListInput struct {
	Search   string
	Page     int
	PageSize int
	Sort     string
}

type ListOutput struct {
	Strategies []StrategyView
	Total      int64
	Page       int
	PageSize   int
}

func (s *ReadService) List(ctx context.Context, input ListInput) (ListOutput, error) {
	page, err := pagination.Resolve(input.Page, input.PageSize, pagination.Config{
		DefaultPage:     defaultPage,
		DefaultPageSize: defaultPageSize,
		MaxPageSize:     maxPageSize,
	})
	if errors.Is(err, pagination.ErrOffsetOverflow) {
		return ListOutput{}, apperr.PageTooLargeError[strategydomain.Strategy]()
	}
	if err != nil {
		return ListOutput{}, err
	}

	sort, err := strategydomain.ParseStrategySort(input.Sort)
	if err != nil {
		return ListOutput{}, err
	}

	filter := ListFilter{
		Search:   strings.TrimSpace(input.Search),
		Page:     page.Page,
		PageSize: page.PageSize,
		Sort:     sort,
	}

	strategies, total, err := s.strategyRepo.List(ctx, filter)
	if err != nil {
		return ListOutput{}, fmt.Errorf("list strategies: %w", err)
	}

	return ListOutput{
		Strategies: strategyViewsFromDomain(strategies),
		Total:      total,
		Page:       page.Page,
		PageSize:   page.PageSize,
	}, nil
}
