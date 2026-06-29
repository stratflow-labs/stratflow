package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	"github.com/stratflow-labs/stratflow/internal/foundation/sort"
)

type Attribute struct {
	ID          uuid.UUID
	StrategyID  uuid.UUID
	Slug        string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdateAttribute struct {
	Slug        *string
	Name        *string
	Description *string
}

func NewAttribute(id, strategyID uuid.UUID, slug, name, description string, now time.Time) (*Attribute, error) {
	a := &Attribute{
		ID:          id,
		StrategyID:  strategyID,
		Slug:        SanitizeString(slug),
		Name:        SanitizeString(name),
		Description: SanitizeString(description),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := a.Validate(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Attribute) Update(update UpdateAttribute) error {
	if update.IsZero() {
		return apperr.UpdateEmptyError[Attribute]()
	}

	next := *a

	if update.Slug != nil {
		next.Slug = SanitizeString(*update.Slug)
	}
	if update.Name != nil {
		next.Name = SanitizeString(*update.Name)
	}
	if update.Description != nil {
		next.Description = SanitizeString(*update.Description)
	}

	if err := next.Validate(); err != nil {
		return err
	}

	*a = next
	return nil
}

func (u UpdateAttribute) IsZero() bool {
	return u.Slug == nil &&
		u.Name == nil &&
		u.Description == nil
}

func (a *Attribute) Validate() error {
	fields := make([]apperr.FieldViolation, 0, 6)

	requiredUUID(&fields, "id", a.ID)
	requiredUUID(&fields, "strategyId", a.StrategyID)
	fields = append(fields, ValidateAttribute(a.Slug, a.Name, a.Description)...)

	if len(fields) > 0 {
		return ValidationError{Fields: fields}
	}

	return nil
}

type AttributeSort string

const (
	AttributeSortCreatedAtDesc AttributeSort = "created_at_desc"
	AttributeSortCreatedAtAsc  AttributeSort = "created_at_asc"
)

var validAttributeSorts = []AttributeSort{
	AttributeSortCreatedAtDesc,
	AttributeSortCreatedAtAsc,
}

func ParseAttributeSort(s string) (AttributeSort, error) {
	return sort.ParseSort(
		s,
		validAttributeSorts,
		AttributeSortCreatedAtDesc,
		func() error { return apperr.SortInvalidError[Attribute]() },
	)
}
