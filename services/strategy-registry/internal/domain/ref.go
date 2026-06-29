package domain

import "github.com/google/uuid"

type EntityRef struct {
	ID   uuid.UUID
	Slug string
}

func RefByID(id uuid.UUID) EntityRef {
	return EntityRef{ID: id}
}

func RefBySlug(slug string) EntityRef {
	return EntityRef{Slug: SanitizeString(slug)}
}

func (r EntityRef) IsZero() bool {
	return r.ID == uuid.Nil && SanitizeString(r.Slug) == ""
}

func (r EntityRef) HasID() bool {
	return r.ID != uuid.Nil
}

func (r EntityRef) NormalizedSlug() string {
	return SanitizeString(r.Slug)
}
