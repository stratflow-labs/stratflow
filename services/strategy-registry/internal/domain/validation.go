package domain

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	validation "github.com/stratflow-labs/stratflow/internal/foundation/validator"
)

const (
	slugMinLen        = 2
	slugMaxLen        = 64
	nameMaxLen        = 120
	valueMaxLen       = 255
	descriptionMaxLen = 1000
)

type ValidationError struct {
	Fields []apperr.FieldViolation
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed with %d field violation(s)", len(e.Fields))
}

func ValidateStrategy(slug, name, description string) []apperr.FieldViolation {
	return validateSlugNameDescription(slug, name, description)
}

func ValidateAttribute(slug, name, description string) []apperr.FieldViolation {
	return validateSlugNameDescription(slug, name, description)
}

func ValidateAttributeValue(slug, value string) []apperr.FieldViolation {
	fields := make([]apperr.FieldViolation, 0, 2)
	validation.Slug(&fields, "slug", slug, validation.SlugConfig{
		MinLen: slugMinLen,
		MaxLen: slugMaxLen,
	})
	validation.RequiredString(&fields, "value", value)
	validation.MaxString(&fields, "value", value, valueMaxLen)
	return fields
}

func validateSlugNameDescription(slug, name, description string) []apperr.FieldViolation {
	fields := make([]apperr.FieldViolation, 0, 3)
	validation.Slug(&fields, "slug", slug, validation.SlugConfig{
		MinLen: slugMinLen,
		MaxLen: slugMaxLen,
	})
	validation.RequiredString(&fields, "name", name)
	validation.MaxString(&fields, "name", name, nameMaxLen)
	validation.MaxString(&fields, "description", description, descriptionMaxLen)
	return fields
}

func requiredUUID(fields *[]apperr.FieldViolation, field string, value uuid.UUID) {
	validation.RequiredUUID(fields, field, value)
}

func SanitizeString(s string) string {
	return validation.SanitizeString(s)
}
