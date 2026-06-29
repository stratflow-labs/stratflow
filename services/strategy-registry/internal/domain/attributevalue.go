package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
)

type AttributeValue struct {
	ID          uuid.UUID
	StrategyID  uuid.UUID
	AttributeID uuid.UUID
	Slug        string
	Value       string
	Relations   []AttributeValueRelation
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdateAttributeValue struct {
	Slug  *string
	Value *string
}

func NewAttributeValue(id, strategyID, attributeID uuid.UUID, slug, value string, now time.Time) (*AttributeValue, error) {
	av := &AttributeValue{
		ID:          id,
		StrategyID:  strategyID,
		AttributeID: attributeID,
		Slug:        SanitizeString(slug),
		Value:       SanitizeString(value),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := av.Validate(); err != nil {
		return nil, err
	}

	return av, nil
}

func (av *AttributeValue) Update(update UpdateAttributeValue) error {
	if update.IsZero() {
		return apperr.UpdateEmptyError[AttributeValue]()
	}

	next := *av

	if update.Slug != nil {
		next.Slug = SanitizeString(*update.Slug)
	}
	if update.Value != nil {
		next.Value = SanitizeString(*update.Value)
	}

	if err := next.Validate(); err != nil {
		return err
	}

	*av = next
	return nil
}

func (u UpdateAttributeValue) IsZero() bool {
	return u.Slug == nil && u.Value == nil
}

func (av *AttributeValue) Validate() error {
	fields := make([]apperr.FieldViolation, 0, 6)

	requiredUUID(&fields, "id", av.ID)
	requiredUUID(&fields, "strategyId", av.StrategyID)
	requiredUUID(&fields, "attributeId", av.AttributeID)
	fields = append(fields, ValidateAttributeValue(av.Slug, av.Value)...)

	if len(fields) > 0 {
		return ValidationError{Fields: fields}
	}

	return nil
}
