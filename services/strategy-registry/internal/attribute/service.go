package attribute

import (
	"context"

	tx "github.com/stratflow-labs/stratflow/internal/foundation/tx"
)

type ReadService struct {
	attributeRepo      AttributeRepository
	attributeValueRepo AttributeValueLookup
}

func NewReadService(
	attributeRepo AttributeRepository,
	attributeValueRepo AttributeValueLookup,
) *ReadService {
	return &ReadService{
		attributeRepo:      attributeRepo,
		attributeValueRepo: attributeValueRepo,
	}
}

type WriteService struct {
	attributeRepo AttributeRepository
	txManager     tx.Manager
	clock         Clock
}

func NewWriteService(
	attributeRepo AttributeRepository,
	txManager tx.Manager,
	clock Clock,
) *WriteService {
	return &WriteService{
		attributeRepo: attributeRepo,
		txManager:     txManager,
		clock:         clock,
	}
}

type Service struct {
	read  *ReadService
	write *WriteService
}

func NewService(read *ReadService, write *WriteService) *Service {
	return &Service{
		read:  read,
		write: write,
	}
}

func (s *Service) Get(ctx context.Context, input GetInput) (AttributeView, error) {
	return s.read.Get(ctx, input)
}

func (s *Service) List(ctx context.Context, input ListInput) (ListOutput, error) {
	return s.read.List(ctx, input)
}

func (s *Service) Create(ctx context.Context, input CreateInput) (AttributeView, error) {
	return s.write.Create(ctx, input)
}

func (s *Service) Update(ctx context.Context, input UpdateInput) (AttributeView, error) {
	return s.write.Update(ctx, input)
}

func (s *Service) Delete(ctx context.Context, input DeleteInput) error {
	return s.write.Delete(ctx, input)
}
