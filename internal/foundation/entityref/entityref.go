package entityref

import (
	"strings"

	"github.com/google/uuid"
)

// EntityRef represents a reference to an entity by ID or slug.
// Used for flexible lookups in repositories and API handlers.
type EntityRef struct {
	ID   uuid.UUID
	Slug string
}

// RefByID creates an EntityRef with only ID set.
func RefByID(id uuid.UUID) EntityRef {
	return EntityRef{ID: id}
}

// RefBySlug creates an EntityRef with only slug set (normalized).
func RefBySlug(slug string) EntityRef {
	return EntityRef{Slug: SanitizeString(slug)}
}

// IsZero returns true if both ID and Slug are empty.
func (r EntityRef) IsZero() bool {
	return r.ID == uuid.Nil && strings.TrimSpace(r.Slug) == ""
}

// HasID returns true if ID is set.
func (r EntityRef) HasID() bool {
	return r.ID != uuid.Nil
}

// NormalizedSlug returns the slug with trimmed whitespace.
func (r EntityRef) NormalizedSlug() string {
	return SanitizeString(r.Slug)
}

// SanitizeString trims whitespace from a string.
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}
