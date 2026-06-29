package systemclock

import (
	"time"

	"github.com/stratflow-labs/stratflow/internal/foundation/clock"
)

// Clock returns wall-clock time in UTC.
type Clock struct{}

var _ clock.Clock = Clock{}

func New() Clock {
	return Clock{}
}

func (Clock) Now() time.Time {
	return time.Now().UTC()
}
