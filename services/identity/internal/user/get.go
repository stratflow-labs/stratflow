package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"
)

func (s *Service) Get(ctx context.Context, id uuid.UUID) (identitydomain.User, error) {
	if id == uuid.Nil {
		return identitydomain.User{}, identitydomain.ErrUserNotFound
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return identitydomain.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}
