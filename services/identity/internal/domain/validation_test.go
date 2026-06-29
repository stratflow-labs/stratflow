package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeLoginInput(t *testing.T) {
	identity, password, err := NormalizeLoginInput("  user@example.com ", " secret ")
	require.NoError(t, err)
	require.Equal(t, "user@example.com", identity)
	require.Equal(t, "secret", password)

	_, _, err = NormalizeLoginInput("  ", "secret")
	require.ErrorIs(t, err, ErrInvalidCredentials)

	_, _, err = NormalizeLoginInput("user@example.com", " ")
	require.ErrorIs(t, err, ErrInvalidCredentials)

	identity, _, err = NormalizeLoginInput(" Trader-One ", "secret")
	require.NoError(t, err)
	require.Equal(t, "trader-one", identity)

	_, _, err = NormalizeLoginInput("invalid@login", "secret")
	require.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestNormalizeUsername(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "letters and digits", input: "Trader42", want: "trader42"},
		{name: "common separators", input: "trader.one-two_3", want: "trader.one-two_3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeUsername(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}

	for _, invalid := range []string{"ab", "user@example.com", "_user", "user_", "user..name", "user name"} {
		t.Run("invalid "+invalid, func(t *testing.T) {
			_, err := NormalizeUsername(invalid)
			require.ErrorIs(t, err, ErrLoginInvalid)
		})
	}
}

func TestNormalizeAccessToken(t *testing.T) {
	token, err := NormalizeAccessToken("  abc ")
	require.NoError(t, err)
	require.Equal(t, "abc", token)

	_, err = NormalizeAccessToken("  ")
	require.ErrorIs(t, err, ErrAccessTokenInvalid)
}

func TestNormalizeRegistrationInput(t *testing.T) {
	email := " user@example.com "
	gender := 1

	input, err := NormalizeRegistrationInput(RegistrationInput{
		Login:    " trader ",
		Name:     " John ",
		LastName: " Smith ",
		Password: " secret123 ",
		Email:    &email,
		Gender:   &gender,
	})
	require.NoError(t, err)
	require.Equal(t, "trader", input.Login)
	require.Equal(t, "John", input.Name)
	require.Equal(t, "Smith", input.LastName)
	require.Equal(t, "secret123", input.Password)
	require.NotNil(t, input.Email)
	require.Equal(t, "user@example.com", *input.Email)
	require.NotNil(t, input.Gender)
	require.Equal(t, 1, *input.Gender)
}

func TestNormalizeRegistrationInput_OptionalFields(t *testing.T) {
	invalidEmail := "invalid"
	invalidGender := 9

	_, err := NormalizeRegistrationInput(RegistrationInput{
		Login:    "trader",
		Name:     "John",
		LastName: "Smith",
		Password: "secret123",
		Email:    &invalidEmail,
	})
	require.ErrorIs(t, err, ErrEmailInvalid)

	_, err = NormalizeRegistrationInput(RegistrationInput{
		Login:    "trader",
		Name:     "John",
		LastName: "Smith",
		Password: "secret123",
		Gender:   &invalidGender,
	})
	require.ErrorIs(t, err, ErrGenderInvalid)
}

func TestDomainValidationClassifiers(t *testing.T) {
	require.True(t, IsLoginValidation(ErrInvalidCredentials))
	require.True(t, IsLoginValidation(ErrUserNotFound))
	require.True(t, IsLoginValidation(ErrPasswordMismatch))
	require.False(t, IsLoginValidation(ErrEmailInvalid))

	require.True(t, IsRegistrationValidation(ErrLoginRequired))
	require.True(t, IsRegistrationValidation(ErrLoginInvalid))
	require.True(t, IsRegistrationValidation(ErrEmailInvalid))
	require.False(t, IsRegistrationValidation(ErrInvalidCredentials))
}

func TestValidatorRejectsInvalidEmail(t *testing.T) {
	err := NewValidator().ValidateUserData("John", "not-an-email", "user")
	require.ErrorIs(t, err, ErrEmailInvalid)
}
