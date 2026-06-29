package clock

import "time"

// Clock defines a minimal contract for retrieving the current time.
type Clock interface {
	Now() time.Time
}
