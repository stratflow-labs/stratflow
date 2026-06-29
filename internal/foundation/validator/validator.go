package validation

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	once sync.Once
	v    *validator.Validate
)

// V is a singleton *validator.Validate with shared settings.
func V() *validator.Validate {
	once.Do(func() {
		val := validator.New()

		// Error field names come from json tags; otherwise use lowerCamelCase.
		val.RegisterTagNameFunc(func(f reflect.StructField) string {
			if j := f.Tag.Get("json"); j != "" && j != "-" {
				return strings.Split(j, ",")[0]
			}
			if f.Name == "" {
				return ""
			}
			return strings.ToLower(f.Name[:1]) + f.Name[1:]
		})

		registerCommon(val)
		v = val
	})
	return v
}

// Register is an extension point for features (call it from sync.Once in the feature).
func Register(fn func(*validator.Validate) error) error { return fn(V()) }
