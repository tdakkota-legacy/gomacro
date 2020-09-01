package derive

import "errors"

var ErrCycleDetected = errors.New("cycle detected")
var ErrInvalidType = errors.New("got invalid type")
