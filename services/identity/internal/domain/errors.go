package domain

import "errors"

var (
	ErrInvalidCredentials  = errors.New("auth: invalid credentials")
	ErrUserNotFound        = errors.New("auth: user not found")
	ErrPasswordMismatch    = errors.New("auth: password mismatch")
	ErrAccessTokenInvalid  = errors.New("auth: access token invalid")
	ErrLoginRequired       = errors.New("auth: login is required")
	ErrNameRequired        = errors.New("auth: name is required")
	ErrLastNameRequired    = errors.New("auth: last name is required")
	ErrPasswordRequired    = errors.New("auth: password is required")
	ErrPasswordTooShort    = errors.New("auth: password is too short")
	ErrPasswordTooLong     = errors.New("auth: password is too long")
	ErrLoginAlreadyUsed    = errors.New("auth: login already used")
	ErrLoginInvalid        = errors.New("auth: login format is invalid")
	ErrEmailAlreadyUsed    = errors.New("auth: email already used")
	ErrEmailInvalid        = errors.New("auth: email is invalid")
	ErrGenderInvalid       = errors.New("auth: gender is invalid")
	ErrUserPageTooLarge    = errors.New("auth: page is too large")
	ErrAccessTokenNotFound = errors.New("access token not found")
)

func IsLoginValidation(err error) bool {
	return errors.Is(err, ErrInvalidCredentials) ||
		errors.Is(err, ErrUserNotFound) ||
		errors.Is(err, ErrPasswordMismatch)
}

func IsRegistrationValidation(err error) bool {
	return errors.Is(err, ErrLoginRequired) ||
		errors.Is(err, ErrNameRequired) ||
		errors.Is(err, ErrLastNameRequired) ||
		errors.Is(err, ErrPasswordRequired) ||
		errors.Is(err, ErrPasswordTooShort) ||
		errors.Is(err, ErrPasswordTooLong) ||
		errors.Is(err, ErrPasswordWhitespace) ||
		errors.Is(err, ErrLoginAlreadyUsed) ||
		errors.Is(err, ErrLoginInvalid) ||
		errors.Is(err, ErrEmailAlreadyUsed) ||
		errors.Is(err, ErrEmailInvalid) ||
		errors.Is(err, ErrGenderInvalid)
}
