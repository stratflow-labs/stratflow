package domain

import (
	"errors"
	"strings"
	"unicode"

	foundationvalidation "github.com/stratflow-labs/stratflow/internal/foundation/validator"
)

var (
	ErrPasswordEmpty      = errors.New("password is empty")
	ErrPasswordLength     = errors.New("password must be between 5 and 32 characters")
	ErrPasswordWhitespace = errors.New("password must not contain spaces")
	ErrNameEmpty          = errors.New("name is required")
	ErrEmailEmpty         = errors.New("email is required")
	ErrRoleEmpty          = errors.New("role is required")
)

// Validator validates user-related business rules.
type Validator struct{}

// NewValidator creates a new user domain validator.
func NewValidator() *Validator {
	return &Validator{}
}

// ValidatePassword checks if password meets domain requirements.
func (v *Validator) ValidatePassword(password string) error {
	trimmed := strings.TrimSpace(password)

	if trimmed == "" {
		return ErrPasswordEmpty
	}

	length := len(trimmed)
	if length < 5 || length > 32 {
		return ErrPasswordLength
	}

	if strings.IndexFunc(trimmed, unicode.IsSpace) >= 0 {
		return ErrPasswordWhitespace
	}

	return nil
}

// ValidateUserData checks basic user data requirements.
func (v *Validator) ValidateUserData(name, email, role string) error {
	if strings.TrimSpace(name) == "" {
		return ErrNameEmpty
	}

	if strings.TrimSpace(email) == "" {
		return ErrEmailEmpty
	}
	if !foundationvalidation.IsEmail(email) {
		return ErrEmailInvalid
	}

	if strings.TrimSpace(role) == "" {
		return ErrRoleEmpty
	}

	return nil
}

// SanitizeString trims whitespace from a string.
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// GenderOrZero returns the value or zero if nil.
func GenderOrZero(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}
