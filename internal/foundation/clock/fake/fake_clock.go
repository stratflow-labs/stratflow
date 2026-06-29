package fakeclock

import (
	"time"

	"github.com/stratflow-labs/stratflow/internal/foundation/clock"
)

// Clock is a controllable clock for tests.
type Clock struct {
	Fixed time.Time
}

var _ clock.Clock = (*Clock)(nil)

func New(fixed time.Time) Clock {
	return Clock{Fixed: fixed}
}

func (c Clock) Now() time.Time {
	return c.Fixed
}
