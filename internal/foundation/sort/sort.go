// internal/foundation/sort/sort.go
package sort

import (
	"strings"
)

func ParseSort[T ~string](
	s string,
	validValues []T,
	defaultValue T,
	errorFactory func() error,
) (T, error) {
	normalized := T(strings.ToLower(strings.TrimSpace(s)))

	if normalized == "" {
		return defaultValue, nil
	}

	for _, v := range validValues {
		if normalized == v {
			return normalized, nil
		}
	}

	return defaultValue, errorFactory()
}
