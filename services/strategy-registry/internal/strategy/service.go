package strategy

import (
	"context"

	"github.com/google/uuid"
	tx "github.com/stratflow-labs/stratflow/internal/foundation/tx"
)

type ReadService struct {
	strategyRepo StrategyRepository
}

func NewReadService(strategyRepo StrategyRepository) *ReadService {
	return &ReadService{strategyRepo: strategyRepo}
}

type WriteService struct {
	strategyRepo      StrategyRepository
	strategyCloneRepo StrategyCloneRepository
	txManager         tx.Manager
	clock             Clock
}

func NewWriteService(
	strategyRepo StrategyRepository,
	strategyCloneRepo StrategyCloneRepository,
	txManager tx.Manager,
	clock Clock,
) *WriteService {
	return &WriteService{
		strategyRepo:      strategyRepo,
		strategyCloneRepo: strategyCloneRepo,
		txManager:         txManager,
		clock:             clock,
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

func (s *Service) Get(ctx context.Context, id uuid.UUID) (StrategyView, error) {
	return s.read.Get(ctx, id)
}

func (s *Service) List(ctx context.Context, input ListInput) (ListOutput, error) {
	return s.read.List(ctx, input)
}

func (s *Service) Create(ctx context.Context, input CreateInput) (StrategyView, error) {
	return s.write.Create(ctx, input)
}

func (s *Service) Update(ctx context.Context, input UpdateInput) (StrategyView, error) {
	return s.write.Update(ctx, input)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.write.Delete(ctx, id)
}

func (s *Service) Clone(ctx context.Context, input CloneInput) (CloneOutput, error) {
	return s.write.Clone(ctx, input)
}
