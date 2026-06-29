package validation

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	emailPattern    = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	usernamePattern = regexp.MustCompile(`^[A-Za-z0-9]+(?:[._-][A-Za-z0-9]+)*$`)
)

const (
	UsernameMinLength = 3
	UsernameMaxLength = 32
)

// Shared rules/aliases.
func registerCommon(v *validator.Validate) {
	_ = v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return IsUsername(fl.Field().String())
	})
	_ = v.RegisterValidation("login_identifier", func(fl validator.FieldLevel) bool {
		return IsLoginIdentifier(fl.Field().String())
	})

	// Do not assign to the blank identifier: the method returns nothing.
	v.RegisterAlias("uuid4", "uuid")

	_ = v.RegisterValidation("password_no_space", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		return strings.IndexFunc(s, unicode.IsSpace) == -1
	})
}

func IsEmail(value string) bool {
	return emailPattern.MatchString(strings.TrimSpace(value))
}

func IsUsername(value string) bool {
	normalized := strings.TrimSpace(value)
	return len(normalized) >= UsernameMinLength &&
		len(normalized) <= UsernameMaxLength &&
		usernamePattern.MatchString(normalized)
}

func IsLoginIdentifier(value string) bool {
	normalized := strings.TrimSpace(value)
	if strings.Contains(normalized, "@") {
		return IsEmail(normalized)
	}
	return IsUsername(normalized)
}
