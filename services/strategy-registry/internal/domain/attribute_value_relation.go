package domain

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	"github.com/stratflow-labs/stratflow/internal/foundation/sort"
)

type AttributeValueRelation struct {
	FromAttributeID uuid.UUID
	FromValueID     uuid.UUID
	ToAttributeID   uuid.UUID
	ToValueID       uuid.UUID
}

func NewAttributeValueRelation(
	fromAttributeID uuid.UUID,
	fromValueID uuid.UUID,
	toAttributeID uuid.UUID,
	toValueID uuid.UUID,
) (AttributeValueRelation, error) {
	relation := AttributeValueRelation{
		FromAttributeID: fromAttributeID,
		FromValueID:     fromValueID,
		ToAttributeID:   toAttributeID,
		ToValueID:       toValueID,
	}

	if err := relation.Validate(); err != nil {
		return AttributeValueRelation{}, err
	}

	return relation, nil
}

func (r AttributeValueRelation) Validate() error {
	var fields []apperr.FieldViolation

	validateAttributeValueRelationFields(&fields, "", r)

	if len(fields) > 0 {
		return ValidationError{Fields: fields}
	}

	return nil
}

func (r AttributeValueRelation) IsSelfReference() bool {
	return r.FromAttributeID == r.ToAttributeID &&
		r.FromValueID == r.ToValueID
}

func ValidateAttributeValueRelations(relations []AttributeValueRelation) error {
	var fields []apperr.FieldViolation

	seen := make(map[AttributeValueRelation]struct{}, len(relations))

	for i, relation := range relations {
		fieldPrefix := "relations[" + strconv.Itoa(i) + "]"

		validateAttributeValueRelationFields(&fields, fieldPrefix, relation)

		if _, ok := seen[relation]; ok {
			fields = append(fields, apperr.FieldViolation{
				Field:   fieldPrefix,
				Code:    "duplicate",
				Message: "duplicate relation is not allowed",
			})
			continue
		}

		seen[relation] = struct{}{}
	}

	if len(fields) > 0 {
		return ValidationError{Fields: fields}
	}

	return nil
}

func DedupeAttributeValueRelations(relations []AttributeValueRelation) []AttributeValueRelation {
	if len(relations) == 0 {
		return []AttributeValueRelation{}
	}

	deduped := make([]AttributeValueRelation, 0, len(relations))
	seen := make(map[AttributeValueRelation]struct{}, len(relations))
	for i := range relations {
		if _, ok := seen[relations[i]]; ok {
			continue
		}
		seen[relations[i]] = struct{}{}
		deduped = append(deduped, relations[i])
	}

	return deduped
}

func validateAttributeValueRelationFields(
	fields *[]apperr.FieldViolation,
	prefix string,
	relation AttributeValueRelation,
) {
	field := func(name string) string {
		if prefix == "" {
			return name
		}
		return prefix + "." + name
	}

	requiredUUID(fields, field("fromAttributeId"), relation.FromAttributeID)
	requiredUUID(fields, field("fromValueId"), relation.FromValueID)
	requiredUUID(fields, field("toAttributeId"), relation.ToAttributeID)
	requiredUUID(fields, field("toValueId"), relation.ToValueID)

	if relation.IsSelfReference() {
		*fields = append(*fields, apperr.FieldViolation{
			Field:   prefix,
			Code:    "selfReference",
			Message: "relation cannot reference itself",
		})
	}
}

type AttributeValueSort string

const (
	AttributeValueSortCreatedAtDesc AttributeValueSort = "created_at_desc"
	AttributeValueSortCreatedAtAsc  AttributeValueSort = "created_at_asc"
)

var validAttributeValueSorts = []AttributeValueSort{
	AttributeValueSortCreatedAtDesc,
	AttributeValueSortCreatedAtAsc,
}

func ParseAttributeValueSort(s string) (AttributeValueSort, error) {
	return sort.ParseSort(
		s,
		validAttributeValueSorts,
		AttributeValueSortCreatedAtDesc,
		func() error { return apperr.SortInvalidError[AttributeValue]() },
	)
}
