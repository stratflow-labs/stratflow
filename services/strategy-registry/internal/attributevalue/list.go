package attributevalue

import (
	"context"
	"errors"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	"github.com/stratflow-labs/stratflow/internal/foundation/pagination"
	attributeValuedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

type ListInput struct {
	StrategyID  uuid.UUID
	AttributeID uuid.UUID
	Search      string
	Page        int
	PageSize    int
	Sort        string
}

type ListOutput struct {
	AttributeValues []AttributeValueView
	Total           int64
	Page            int
	PageSize        int
}

func (s *Service) List(ctx context.Context, input ListInput) (ListOutput, error) {
	if input.StrategyID == uuid.Nil {
		return ListOutput{}, apperr.NotFoundError[attributeValuedomain.Strategy]()
	}
	if input.AttributeID == uuid.Nil {
		return ListOutput{}, apperr.NotFoundError[attributeValuedomain.Attribute]()
	}

	page, err := pagination.Resolve(input.Page, input.PageSize, pagination.Config{
		DefaultPage:     defaultPage,
		DefaultPageSize: defaultPageSize,
		MaxPageSize:     maxPageSize,
	})
	if errors.Is(err, pagination.ErrOffsetOverflow) {
		return ListOutput{}, apperr.PageTooLargeError[attributeValuedomain.AttributeValue]()
	}
	if err != nil {
		return ListOutput{}, err
	}

	sort, err := attributeValuedomain.ParseAttributeValueSort(input.Sort)
	if err != nil {
		return ListOutput{}, err
	}

	filter := ListFilter{
		StrategyID:  input.StrategyID,
		AttributeID: input.AttributeID,
		Search:      input.Search,
		Page:        page.Page,
		PageSize:    page.PageSize,
		Sort:        sort,
	}

	attributeValues, total, err := s.attributeValueRepo.List(ctx, filter)
	if err != nil {
		return ListOutput{}, fmt.Errorf("list attributeValues: %w", err)
	}

	views := attributeValueViewsFromDomain(attributeValues)
	if len(views) == 0 {
		return ListOutput{AttributeValues: views, Total: total, Page: page.Page, PageSize: page.PageSize}, nil
	}

	valueIDs := make([]uuid.UUID, 0, len(views))
	for i := range views {
		valueIDs = append(valueIDs, views[i].ID)
	}

	relationsByFromValueID, err := s.attributeValueRepo.ListRelationsByFromValueIDsForAttributeValues(ctx, input.StrategyID, valueIDs)
	if err != nil {
		return ListOutput{}, fmt.Errorf("list attributeValue relations: %w", err)
	}

	for i := range views {
		relations, ok := relationsByFromValueID[views[i].ID]
		if !ok {
			continue
		}
		views[i].Relations = relations
	}

	return ListOutput{AttributeValues: views, Total: total, Page: page.Page, PageSize: page.PageSize}, nil
}
