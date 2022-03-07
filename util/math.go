package util

import (
	"errors"
	"math"
)

var ErrOverflow = errors.New("integer overflow")

func Int32(i int) (int32, error) {
	if i > math.MaxInt32 || i < math.MinInt32 {
		return 0, ErrOverflow
	}
	return int32(i), nil
}

func Int16(i int) (int16, error) {
	if i > math.MaxInt16 || i < math.MinInt16 {
		return 0, ErrOverflow
	}
	return int16(i), nil
}

func Int8(i int) (int8, error) {
	if i > math.MaxInt8 || i < math.MinInt8 {
		return 0, ErrOverflow
	}
	return int8(i), nil
}
