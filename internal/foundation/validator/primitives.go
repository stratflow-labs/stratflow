package validation

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:[-_][a-z0-9]+)*$`)

type SlugConfig struct {
	MinLen int
	MaxLen int
}

func RequiredString(fields *[]apperr.FieldViolation, field, value string) {
	if value == "" {
		*fields = append(*fields, apperr.Required(field))
	}
}

func RequiredUUID(fields *[]apperr.FieldViolation, field string, value uuid.UUID) {
	if value == uuid.Nil {
		*fields = append(*fields, apperr.Required(field))
	}
}

func MaxString(fields *[]apperr.FieldViolation, field, value string, max int) {
	if len(value) > max {
		*fields = append(*fields, apperr.FieldViolation{
			Field:   field,
			Code:    "maxLength",
			Message: field + " must be at most " + strconv.Itoa(max) + " characters",
		})
	}
}

func Slug(fields *[]apperr.FieldViolation, field, value string, cfg SlugConfig) {
	if value == "" {
		*fields = append(*fields, apperr.Required(field))
		return
	}

	if len(value) < cfg.MinLen {
		*fields = append(*fields, apperr.FieldViolation{
			Field:   field,
			Code:    "minLength",
			Message: field + " must be at least " + strconv.Itoa(cfg.MinLen) + " characters",
		})
	}

	if len(value) > cfg.MaxLen {
		*fields = append(*fields, apperr.FieldViolation{
			Field:   field,
			Code:    "maxLength",
			Message: field + " must be at most " + strconv.Itoa(cfg.MaxLen) + " characters",
		})
	}

	if !slugPattern.MatchString(value) {
		*fields = append(*fields, apperr.FieldViolation{
			Field:   field,
			Code:    "invalidFormat",
			Message: field + " must contain only lowercase letters, numbers, hyphens, and underscores",
		})
	}
}

func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}
