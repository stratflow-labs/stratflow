package user

import (
	"context"
	"fmt"

	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"

	"github.com/google/uuid"
)

type CreateInput struct {
	Login    string
	Name     string
	LastName string
	Email    string
	Role     string
	Gender   *int
	Password string
}

func (s *Service) Create(ctx context.Context, input CreateInput) (identitydomain.User, error) {
	login, err := identitydomain.NormalizeUsername(input.Login)
	if err != nil {
		return identitydomain.User{}, fmt.Errorf("validation failed: %w", err)
	}
	input.Login = login

	if err := s.validateInput(input); err != nil {
		return identitydomain.User{}, fmt.Errorf("validation failed: %w", err)
	}

	item := s.buildUser(input)

	hash, err := s.hasher.Hash(ctx, identitydomain.SanitizeString(input.Password))
	if err != nil {
		return identitydomain.User{}, fmt.Errorf("hash password: %w", err)
	}
	item.PasswordHash = hash

	now := s.clock.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	created, err := s.userRepo.Create(ctx, &item)
	if err != nil {
		return identitydomain.User{}, fmt.Errorf("create user: %w", err)
	}

	return created, nil
}

func (s *Service) validateInput(input CreateInput) error {
	if err := s.validator.ValidateUserData(input.Name, input.Email, input.Role); err != nil {
		return err
	}

	if err := s.validator.ValidatePassword(input.Password); err != nil {
		return err
	}

	return nil
}

func (s *Service) buildUser(input CreateInput) identitydomain.User {
	return identitydomain.User{
		ID:       uuid.New(),
		Login:    input.Login,
		Name:     identitydomain.SanitizeString(input.Name),
		LastName: identitydomain.SanitizeString(input.LastName),
		Email:    identitydomain.SanitizeString(input.Email),
		Role:     identitydomain.SanitizeString(input.Role),
		Gender:   identitydomain.GenderOrZero(input.Gender),
	}
}
