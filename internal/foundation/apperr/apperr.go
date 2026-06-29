package apperr

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Kind string

const (
	KindNotFound         Kind = "not_found"
	KindAlreadyExists    Kind = "already_exists"
	KindInvalidArgument  Kind = "invalid_argument"
	KindPermissionDenied Kind = "permission_denied"
	KindUnauthenticated  Kind = "unauthenticated"
	KindInternal         Kind = "internal"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrAlreadyExists  = errors.New("already exists")
	ErrUpdateEmpty    = errors.New("at least one field must be provided")
	ErrSortInvalid    = errors.New("sort is invalid")
	ErrPageTooLarge   = errors.New("page size is too large")
	ErrCloneEmpty     = errors.New("clone items are empty")
	ErrBatchEmpty     = errors.New("batch actions are empty")
	ErrBatchDuplicate = errors.New("duplicate items are not allowed")
	ErrSelfReference  = errors.New("self reference is not allowed")
)

type FieldViolation struct {
	Field   string
	Code    string
	Message string
}

type Error struct {
	Kind    Kind
	Code    string
	Entity  string
	Message string
	Fields  []FieldViolation
	Cause   error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	if t.Kind != "" && e.Kind != t.Kind {
		return false
	}
	if t.Code != "" && e.Code != t.Code {
		return false
	}
	return true
}

func New(kind Kind, code, message string) *Error {
	return &Error{
		Kind:    kind,
		Code:    code,
		Message: message,
	}
}

func Wrap(kind Kind, code, message string, cause error) *Error {
	return &Error{
		Kind:    kind,
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func Invalid[T any](action, errorType, message string, cause ...error) *Error {
	return newEntityActionError[T](KindInvalidArgument, action, errorType, message, firstCause(cause))
}

func NotFound[T any](action string, ref any, cause ...error) *Error {
	entity := EntityName[T]()
	return &Error{
		Kind:    KindNotFound,
		Code:    errorCode(entity, action, "notFound"),
		Entity:  entity,
		Message: fmt.Sprintf("%s %v not found", entity, ref),
		Cause:   firstCause(cause),
	}
}

func AlreadyExists[T any](action string, ref any, cause ...error) *Error {
	entity := EntityName[T]()
	return &Error{
		Kind:    KindAlreadyExists,
		Code:    errorCode(entity, action, "alreadyExists"),
		Entity:  entity,
		Message: fmt.Sprintf("%s %v already exists", entity, ref),
		Cause:   firstCause(cause),
	}
}

func Validation[T any](action string, fields []FieldViolation, cause ...error) *Error {
	entity := EntityName[T]()
	return &Error{
		Kind:    KindInvalidArgument,
		Code:    errorCode(entity, action, "validationFailed"),
		Entity:  entity,
		Message: entity + " validation failed",
		Fields:  cloneFields(fields),
		Cause:   firstCause(cause),
	}
}

func ValidationWithCode(code, message string, fields []FieldViolation, cause ...error) *Error {
	return &Error{
		Kind:    KindInvalidArgument,
		Code:    code,
		Message: message,
		Fields:  cloneFields(fields),
		Cause:   firstCause(cause),
	}
}

func NotFoundError[T any]() *Error {
	return newEntityError[T](KindNotFound, "notFound", "%s not found", ErrNotFound)
}

func AlreadyExistsError[T any]() *Error {
	return newEntityError[T](KindAlreadyExists, "alreadyExists", "%s already exists", ErrAlreadyExists)
}

func UpdateEmptyError[T any]() *Error {
	return newEntityError[T](KindInvalidArgument, "updateEmpty", "%s: at least one field must be provided", ErrUpdateEmpty)
}

func SortInvalidError[T any]() *Error {
	return newEntityError[T](KindInvalidArgument, "sortInvalid", "%s: sort is invalid", ErrSortInvalid)
}

func PageTooLargeError[T any]() *Error {
	return newEntityError[T](KindInvalidArgument, "pageTooLarge", "%s: page size is too large", ErrPageTooLarge)
}

func CloneEmptyError[T any]() *Error {
	return newEntityActionError[T](KindInvalidArgument, "clone", "itemsEmpty", "clone items are required", ErrCloneEmpty)
}

func BatchEmptyError[T any](action string) *Error {
	return newEntityActionError[T](KindInvalidArgument, action, "empty", EntityName[T]()+" batch actions are required", ErrBatchEmpty)
}

func DuplicateError[T any](action, errorType, message string) *Error {
	return newEntityActionError[T](KindInvalidArgument, action, errorType, message, ErrBatchDuplicate)
}

func DuplicateFieldInRequestError[T any](action, field string) *Error {
	return DuplicateError[T](
		action,
		"duplicate"+upperFirst(field),
		field+" values must be unique in request",
	)
}

func SelfReferenceError[T any](action, errorType, message string) *Error {
	return newEntityActionError[T](KindInvalidArgument, action, errorType, message, ErrSelfReference)
}

func Required(field string) FieldViolation {
	return FieldViolation{
		Field:   field,
		Code:    "required",
		Message: field + " is required",
	}
}

func RefRequired(field, subject string) FieldViolation {
	return FieldViolation{
		Field:   field,
		Code:    "required",
		Message: subject + " id or slug is required",
	}
}

func IsKind(err error, kind Kind) bool {
	var appErr *Error
	return errors.As(err, &appErr) && appErr.Kind == kind
}

func IsNotFound(err error) bool {
	return IsKind(err, KindNotFound)
}

func IsAlreadyExists(err error) bool {
	return IsKind(err, KindAlreadyExists)
}

func IsInvalidArgument(err error) bool {
	return IsKind(err, KindInvalidArgument)
}

func KindOf(err error) Kind {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Kind
	}
	return ""
}

func CodeOf(err error) string {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return ""
}

func Code(err error) string {
	return CodeOf(err)
}

func MessageOf(err error) string {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Message
	}
	return ""
}

func Message(err error) string {
	return MessageOf(err)
}

func FieldsOf(err error) []FieldViolation {
	var appErr *Error
	if !errors.As(err, &appErr) {
		return nil
	}
	return cloneFields(appErr.Fields)
}

func Fields(err error) []FieldViolation {
	return FieldsOf(err)
}

func EntityName[T any]() string {
	t := reflect.TypeFor[T]()
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	name := t.Name()
	if name == "" {
		return "unknown"
	}

	return lowerCamel(name)
}

func newEntityError[T any](kind Kind, errorType, message string, cause error) *Error {
	entity := EntityName[T]()
	return &Error{
		Kind:    kind,
		Code:    errorCode(entity, "", errorType),
		Entity:  entity,
		Message: fmt.Sprintf(message, entity),
		Cause:   cause,
	}
}

func newEntityActionError[T any](kind Kind, action, errorType, message string, cause error) *Error {
	entity := EntityName[T]()
	return &Error{
		Kind:    kind,
		Code:    errorCode(entity, action, errorType),
		Entity:  entity,
		Message: message,
		Cause:   cause,
	}
}

func errorCode(entity, action, errorType string) string {
	parts := make([]string, 0, 3)
	for _, part := range []string{entity, action, errorType} {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, ".")
}

func firstCause(causes []error) error {
	if len(causes) == 0 {
		return nil
	}
	return causes[0]
}

func cloneFields(fields []FieldViolation) []FieldViolation {
	if len(fields) == 0 {
		return nil
	}

	cloned := make([]FieldViolation, len(fields))
	copy(cloned, fields)
	return cloned
}

func lowerCamel(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
