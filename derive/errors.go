package derive

import "errors"

var (
	// ErrCycleDetected reports that code cannot be derived due to infinite cycle.
	ErrCycleDetected = errors.New("cycle detected")
	// ErrInvalidType reports that type check failed.
	ErrInvalidType = errors.New("got invalid type")
)
