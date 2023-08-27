package format

import (
	"strconv"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

func FormatStrings(val string, col uint) (any, error) {
	return val, nil
}

func FormatStringAndInts(val string, colIdx uint) (any, error) {
	if colIdx == 0 {
		return val, nil
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		return val, err
	}
	return num, nil
}

func FormatStringAndFloats(val string, colIdx uint) (any, error) {
	if colIdx == 0 {
		return val, nil
	}
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return val, err
	}
	return num, nil
}

func FormatTimeAndInts(val string, colIdx uint) (any, error) {
	if colIdx == 0 {
		dt, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return val, err
		} else {
			return dt, nil
		}
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		return val, err
	}
	return num, nil
}

func FormatMonthAndFloats(val string, colIdx uint) (any, error) {
	if colIdx == 0 {
		dt, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return val, err
		} else {
			return dt.Format(timeutil.ISO8601YM), nil
		}
	}
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return val, err
	}
	return num, nil
}

func FormatDateAndFloats(val string, colIdx uint) (any, error) {
	if colIdx == 0 {
		dt, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return val, err
		} else {
			return dt.Format(timeutil.DateMDY), nil
		}
	}
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return val, err
	}
	return num, nil
}

func FormatTimeAndFloats(val string, colIdx uint) (any, error) {
	if colIdx == 0 {
		dt, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return val, err
		} else {
			return dt, nil
		}
	}
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return val, err
	}
	return num, nil
}
