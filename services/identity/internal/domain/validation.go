package domain

import (
	"strings"
	"unicode"

	foundationvalidation "github.com/stratflow-labs/stratflow/internal/foundation/validator"
)

type RegistrationInput struct {
	Login    string
	Name     string
	LastName string
	Password string
	Email    *string
	Gender   *int
}

func NormalizeLoginInput(login, password string) (string, string, error) {
	identity := strings.TrimSpace(login)
	secret := strings.TrimSpace(password)
	if identity == "" || secret == "" {
		return "", "", ErrInvalidCredentials
	}

	if strings.Contains(identity, "@") {
		if !foundationvalidation.IsEmail(identity) {
			return "", "", ErrInvalidCredentials
		}
		return strings.ToLower(identity), secret, nil
	}

	username, err := NormalizeUsername(identity)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}
	return username, secret, nil
}

func NormalizeUsername(login string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(login))
	if normalized == "" {
		return "", ErrLoginRequired
	}
	if !foundationvalidation.IsUsername(normalized) {
		return "", ErrLoginInvalid
	}
	return normalized, nil
}

func NormalizeAccessToken(token string) (string, error) {
	normalized := strings.TrimSpace(token)
	if normalized == "" {
		return "", ErrAccessTokenInvalid
	}
	return normalized, nil
}

func NormalizeRegistrationInput(input RegistrationInput) (RegistrationInput, error) {
	normalizedLogin, err := NormalizeUsername(input.Login)
	if err != nil {
		return RegistrationInput{}, err
	}
	input.Login = normalizedLogin
	input.Name = strings.TrimSpace(input.Name)
	input.LastName = strings.TrimSpace(input.LastName)
	input.Password = strings.TrimSpace(input.Password)
	input.Email = normalizeOptionalString(input.Email)

	switch {
	case input.Name == "":
		return RegistrationInput{}, ErrNameRequired
	case input.LastName == "":
		return RegistrationInput{}, ErrLastNameRequired
	case input.Password == "":
		return RegistrationInput{}, ErrPasswordRequired
	case len(input.Password) < 5:
		return RegistrationInput{}, ErrPasswordTooShort
	case len(input.Password) > 32:
		return RegistrationInput{}, ErrPasswordTooLong
	case strings.IndexFunc(input.Password, unicode.IsSpace) >= 0:
		return RegistrationInput{}, ErrPasswordWhitespace
	case input.Email != nil && !foundationvalidation.IsEmail(*input.Email):
		return RegistrationInput{}, ErrEmailInvalid
	case input.Gender != nil && (*input.Gender < 0 || *input.Gender > 2):
		return RegistrationInput{}, ErrGenderInvalid
	}

	return input, nil
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
