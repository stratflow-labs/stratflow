package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/stratflow-labs/stratflow/internal/foundation/pagination"
	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"
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
	Users    []identitydomain.User
	Total    int64
	Page     int
	PageSize int
}

func (s *Service) List(ctx context.Context, input ListInput) (ListOutput, error) {
	page, err := pagination.Resolve(input.Page, input.PageSize, pagination.Config{
		DefaultPage:     defaultPage,
		DefaultPageSize: defaultPageSize,
		MaxPageSize:     maxPageSize,
	})
	if errors.Is(err, pagination.ErrOffsetOverflow) {
		return ListOutput{}, identitydomain.ErrUserPageTooLarge
	}
	if err != nil {
		return ListOutput{}, err
	}

	filter := ListFilter{
		Search:   input.Search,
		Page:     page.Page,
		PageSize: page.PageSize,
		Sort:     input.Sort,
	}

	users, total, err := s.userRepo.List(ctx, filter)
	if err != nil {
		return ListOutput{}, fmt.Errorf("list users: %w", err)
	}

	return ListOutput{Users: users, Total: total, Page: page.Page, PageSize: page.PageSize}, nil
}
