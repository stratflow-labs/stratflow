package idgen

import "github.com/google/uuid"

// IDGenerator abstracts identifier generation.
type IDGenerator interface {
	NewID() uuid.UUID
}
