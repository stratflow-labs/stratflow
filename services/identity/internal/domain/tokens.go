package domain

import "github.com/google/uuid"

type TokenClaims struct {
	UserID uuid.UUID
	Role   string
}

type IssuedToken struct {
	Value string
}
