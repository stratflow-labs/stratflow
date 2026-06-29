package attribute

import (
	"context"
	"errors"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	"github.com/stratflow-labs/stratflow/internal/foundation/pagination"
	attributevalue "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attributevalue"
	attributedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

type ListInput struct {
	StrategyID uuid.UUID
	Search     string
	Page       int
	PageSize   int
	Sort       string
}

type ListOutput struct {
	Attributes []AttributeView
	Total      int64
	Page       int
	PageSize   int
}

func (s *ReadService) List(ctx context.Context, input ListInput) (ListOutput, error) {
	if input.StrategyID == uuid.Nil {
		return ListOutput{}, apperr.NotFoundError[attributedomain.Strategy]()
	}

	filter, err := buildListFilter(input)
	if err != nil {
		return ListOutput{}, err
	}

	attributes, total, err := s.attributeRepo.List(ctx, filter)
	if err != nil {
		return ListOutput{}, fmt.Errorf("list attributes: %w", err)
	}

	views := attributeViewsFromDomain(attributes)
	if len(views) == 0 {
		return ListOutput{Attributes: views, Total: total, Page: filter.Page, PageSize: filter.PageSize}, nil
	}

	views, err = s.enrichViews(ctx, input.StrategyID, views)
	if err != nil {
		return ListOutput{}, err
	}

	return ListOutput{Attributes: views, Total: total, Page: filter.Page, PageSize: filter.PageSize}, nil
}

func buildListFilter(input ListInput) (ListFilter, error) {
	page, err := pagination.Resolve(input.Page, input.PageSize, pagination.Config{
		DefaultPage:     defaultPage,
		DefaultPageSize: defaultPageSize,
		MaxPageSize:     maxPageSize,
	})
	if errors.Is(err, pagination.ErrOffsetOverflow) {
		return ListFilter{}, apperr.PageTooLargeError[attributedomain.Attribute]()
	}
	if err != nil {
		return ListFilter{}, err
	}

	sort, err := attributedomain.ParseAttributeSort(input.Sort)
	if err != nil {
		return ListFilter{}, err
	}

	return ListFilter{
		StrategyID: input.StrategyID,
		Search:     input.Search,
		Page:       page.Page,
		PageSize:   page.PageSize,
		Sort:       sort,
	}, nil
}

func (s *ReadService) enrichViews(ctx context.Context, strategyID uuid.UUID, views []AttributeView) ([]AttributeView, error) {
	if s.attributeValueRepo == nil {
		return views, nil
	}

	attributeIDs := collectAttributeIDs(views)
	valuesByAttributeID, err := s.attributeValueRepo.ListValueGraphByAttributeIDs(ctx, strategyID, attributeIDs)
	if err != nil {
		return nil, fmt.Errorf("list attribute value graph: %w", err)
	}
	assignValuesToViews(views, valuesByAttributeID)

	return views, nil
}

func collectAttributeIDs(views []AttributeView) []uuid.UUID {
	attributeIDs := make([]uuid.UUID, 0, len(views))
	for i := range views {
		attributeIDs = append(attributeIDs, views[i].ID)
	}
	return attributeIDs
}

func assignValuesToViews(
	views []AttributeView,
	valuesByAttributeID map[uuid.UUID][]attributevalue.AttributeValueView,
) {
	for i := range views {
		values, ok := valuesByAttributeID[views[i].ID]
		if !ok {
			continue
		}
		views[i].Values = values
	}
}
