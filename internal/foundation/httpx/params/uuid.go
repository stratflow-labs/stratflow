package params

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// UUID extracts a route parameter and parses it as UUID.
func UUID(r *http.Request, name string) (uuid.UUID, error) {
	value := strings.TrimSpace(r.PathValue(name))
	if value == "" {
		return uuid.Nil, fmt.Errorf("missing route param %q", name)
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse route param %q: %w", name, err)
	}

	return id, nil
}
