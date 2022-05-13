package timeseries

import (
	"errors"
)

var (
	ErrIntervalNotSupported = errors.New("interval not supported")
	ErrNoTimeItem           = errors.New("no time item")
)
