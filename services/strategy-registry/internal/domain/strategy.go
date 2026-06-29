package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	"github.com/stratflow-labs/stratflow/internal/foundation/sort"
)

type Strategy struct {
	ID          uuid.UUID
	Slug        string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdateStrategy struct {
	Slug        *string
	Name        *string
	Description *string
}

func NewStrategy(id uuid.UUID, slug, name, description string, now time.Time) (*Strategy, error) {
	s := &Strategy{
		ID:          id,
		Slug:        SanitizeString(slug),
		Name:        SanitizeString(name),
		Description: SanitizeString(description),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.Validate(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Strategy) Update(update UpdateStrategy) error {
	if update.IsZero() {
		return apperr.UpdateEmptyError[Strategy]()
	}

	next := *s

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

	*s = next
	return nil
}

func (u UpdateStrategy) IsZero() bool {
	return u.Slug == nil &&
		u.Name == nil &&
		u.Description == nil
}

func (s *Strategy) Validate() error {
	fields := make([]apperr.FieldViolation, 0, 5)

	requiredUUID(&fields, "id", s.ID)
	fields = append(fields, ValidateStrategy(s.Slug, s.Name, s.Description)...)

	if len(fields) > 0 {
		return ValidationError{Fields: fields}
	}

	return nil
}

type StrategySort string

const (
	StrategySortCreatedAtDesc StrategySort = "created_at_desc"
	StrategySortCreatedAtAsc  StrategySort = "created_at_asc"
)

var validStrategySorts = []StrategySort{
	StrategySortCreatedAtDesc,
	StrategySortCreatedAtAsc,
}

// ParseStrategySort парсит строку сортировки в StrategySort
func ParseStrategySort(s string) (StrategySort, error) {
	return sort.ParseSort(
		s,
		validStrategySorts,
		StrategySortCreatedAtDesc,
		func() error { return apperr.SortInvalidError[Strategy]() },
	)
}
