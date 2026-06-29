package user

import (
	"context"
	"fmt"

	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"

	"github.com/google/uuid"
)

type UpdateInput struct {
	ID       uuid.UUID
	Name     *string
	LastName *string
	Email    *string
	Gender   *int
}

func (s *Service) Update(ctx context.Context, input UpdateInput) (identitydomain.User, error) {
	if input.ID == uuid.Nil {
		return identitydomain.User{}, identitydomain.ErrUserNotFound
	}

	var updated identitydomain.User
	err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		current, err := s.userRepo.GetByID(txCtx, input.ID)
		if err != nil {
			return fmt.Errorf("load user: %w", err)
		}

		current.ApplyProfileUpdate(identitydomain.ProfileUpdate{
			Name:     input.Name,
			LastName: input.LastName,
			Email:    input.Email,
			Gender:   input.Gender,
		})
		current.UpdatedAt = s.clock.Now()

		updated, err = s.userRepo.Update(txCtx, &current)
		if err != nil {
			return fmt.Errorf("update user: %w", err)
		}

		return nil
	})
	if err != nil {
		return identitydomain.User{}, fmt.Errorf("update user transaction: %w", err)
	}

	return updated, nil
}
