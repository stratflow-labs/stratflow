package user

import tx "github.com/stratflow-labs/stratflow/internal/foundation/tx"

type Service struct {
	userRepo  UserRepository
	txManager tx.Manager
	hasher    PasswordHasher
	validator UserValidator
	clock     Clock
}

func NewService(
	userRepo UserRepository,
	txManager tx.Manager,
	hasher PasswordHasher,
	validator UserValidator,
	clock Clock,
) *Service {
	return &Service{
		userRepo:  userRepo,
		txManager: txManager,
		hasher:    hasher,
		validator: validator,
		clock:     clock,
	}
}
