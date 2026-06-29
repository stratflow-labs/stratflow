package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID
	Login           string
	Name            string
	LastName        string
	Email           string
	PasswordHash    string
	Role            string
	ImageUrl        string
	Gender          int
	IsEmailVerified bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ProfileUpdate struct {
	Name     *string
	LastName *string
	Email    *string
	Gender   *int
}

// ApplyProfileUpdate applies domain-level user profile mutation rules.
func (u *User) ApplyProfileUpdate(update ProfileUpdate) {
	if u == nil {
		return
	}

	if update.Name != nil {
		u.Name = SanitizeString(*update.Name)
	}
	if update.LastName != nil {
		u.LastName = SanitizeString(*update.LastName)
	}
	if update.Email != nil {
		email := SanitizeString(*update.Email)
		if email != "" && email != u.Email {
			u.Email = email
			u.IsEmailVerified = false
		}
	}
	if update.Gender != nil {
		u.Gender = *update.Gender
	}
}
