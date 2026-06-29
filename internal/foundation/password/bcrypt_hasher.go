package password

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher provides password hashing/verification for both auth and registration flows.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates hasher with the provided cost or bcrypt.DefaultCost when zero/negative.
func NewBcryptHasher(cost int) BcryptHasher {
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	return BcryptHasher{cost: cost}
}

// Hash returns bcrypt hash of a plaintext password.
func (h BcryptHasher) Hash(_ context.Context, plain string) (string, error) {
	out, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash: %w", err)
	}
	return string(out), nil
}

// Compare reports whether plaintext matches the hashed password.
// It returns false without error when hashes mismatch.
func (h BcryptHasher) Compare(_ context.Context, plain, hashed string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	if errors.Is(err, bcrypt.ErrHashTooShort) {
		return false, nil
	}
	var invalidPrefix bcrypt.InvalidHashPrefixError
	if errors.As(err, &invalidPrefix) {
		return false, nil
	}
	return false, fmt.Errorf("bcrypt compare: %w", err)
}
