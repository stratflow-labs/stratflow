package client

import "errors"

var (
	// ErrUnauthorized signals token is missing or invalid.
	ErrUnauthorized = errors.New("unauthorized")
)
