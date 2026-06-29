package user

import (
	"context"

	"github.com/stratflow-labs/stratflow/internal/foundation/clock"
	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *identitydomain.User) (identitydomain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (identitydomain.User, error)
	List(ctx context.Context, filter ListFilter) ([]identitydomain.User, int64, error)
	Update(ctx context.Context, user *identitydomain.User) (identitydomain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type PasswordHasher interface {
	Hash(ctx context.Context, plain string) (string, error)
}

type UserValidator interface {
	ValidatePassword(password string) error
	ValidateUserData(name, email, role string) error
}

type Clock interface {
	clock.Clock
}

type ListFilter struct {
	Search   string
	Page     int
	PageSize int
	Sort     string
}
