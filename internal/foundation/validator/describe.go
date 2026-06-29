package validation

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	CodeRequired     = "validation.required"
	CodeInvalidValue = "validation.invalidValue"
	CodeInvalidType  = "validation.invalidType"
	CodeInvalidJSON  = "validation.invalidJson"
)

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// todo errors.AsType[T](err) not errors.As(err, &target). ALWAYS use errors.AsType when checking if error matches a specific type.
// Describe normalizes validator errors into a unified format.
func Describe(err error) []FieldError {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		out := make([]FieldError, 0, len(verrs))
		for _, ve := range verrs {
			field := toJSONField(ve.Field())
			tag, param := ve.Tag(), ve.Param()
			out = append(out, FieldError{
				Field:   field,
				Code:    codeByTag(tag),
				Message: humanMsg(field, tag, param),
			})
		}
		return out
	}
	return DescribeBindingError(err)
}

func codeByTag(tag string) string {
	switch tag {
	case "required":
		return CodeRequired
	case "uuid", "uuid4":
		return CodeInvalidType
	default:
		return CodeInvalidValue
	}
}

func humanMsg(field, tag, param string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "uuid", "uuid4":
		return field + " must be a valid UUID"
	case "email":
		return field + " must be a valid email address"
	case "username":
		return field + " must be 3-32 characters and contain only letters, digits, dots, underscores, or hyphens"
	case "login_identifier":
		return field + " must be a valid username or email address"
	case "min":
		return field + " must be at least " + param + " characters long"
	case "max":
		return field + " must be at most " + param + " characters long"
	case "len":
		return field + " must be exactly " + param + " characters long"
	default:
		return field + " validation failed (" + tag + ")"
	}
}

func DescribeBindingError(err error) []FieldError {
	if err == nil {
		return nil
	}
	var out []FieldError

	var vErrors validator.ValidationErrors
	//todo errors.AsType[T](err) not errors.As(err, &target). ALWAYS use errors.AsType when checking if error matches a specific type.
	if errors.As(err, &vErrors) {
		for _, ve := range vErrors {
			field := toJSONField(ve.Field())
			code := codeByTag(ve.Tag())
			message := humanMsg(field, ve.Tag(), ve.Param())
			out = append(out, FieldError{Field: field, Code: code, Message: message})
		}
		return out
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		field := toJSONField(typeErr.Field)
		message := "expected " + typeErr.Type.String() + " but got " + typeErr.Value
		out = append(out, FieldError{Field: field, Code: CodeInvalidType, Message: message})
		return out
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		message := "invalid JSON at offset " + strconv.FormatInt(syntaxErr.Offset, 10)
		out = append(out, FieldError{Field: "", Code: CodeInvalidJSON, Message: message})
		return out
	}

	if errors.Is(err, io.EOF) {
		out = append(out, FieldError{Field: "", Code: CodeInvalidJSON, Message: "request body is empty"})
		return out
	}

	out = append(out, FieldError{Field: "", Code: CodeInvalidValue, Message: err.Error()})
	return out
}

func toJSONField(field string) string {
	if field == "" {
		return ""
	}
	return strings.ToLower(field[:1]) + field[1:]
}
