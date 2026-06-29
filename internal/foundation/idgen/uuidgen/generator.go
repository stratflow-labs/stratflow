package uuidgen

import (
	"github.com/stratflow-labs/stratflow/internal/foundation/idgen"

	"github.com/google/uuid"
)

// Generator produces UUID v4 identifiers.
type Generator struct{}

var _ idgen.IDGenerator = Generator{}

func New() Generator {
	return Generator{}
}

func (Generator) NewID() uuid.UUID {
	return uuid.New()
}
