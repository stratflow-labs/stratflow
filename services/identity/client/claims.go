package client

import "github.com/google/uuid"

// Claims is a public, service-agnostic representation of authenticated user data.
type Claims struct {
	UserID uuid.UUID
	Role   string
}
