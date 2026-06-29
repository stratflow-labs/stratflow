package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"
)

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return identitydomain.ErrUserNotFound
	}

	err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		if err := s.userRepo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete user: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
