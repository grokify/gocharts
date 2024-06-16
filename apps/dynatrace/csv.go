package dynatrace

import (
	"strings"
	"time"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/time/timeutil"
)

func ReadCSVTimeDurations(filename string, toTimeUnit time.Duration, prec int, addSuffix bool) (*table.Table, error) {
	if t, err := table.ReadTableSimple(filename, ","); err != nil {
		return nil, err
	} else if err = UpdateTimes(t); err != nil {
		return nil, err
	} else {
		return t, UpdateDurations(t, toTimeUnit, prec, addSuffix)
	}
}

// UpdateTimes updates times to RFC-3339 times.
func UpdateTimes(t *table.Table) error {
	if t == nil {
		return table.ErrTableCannotBeNil
	}
	t.FormatMap[0] = table.FormatTime
	return t.FormatColumns(0, 0, func(s string) (string, error) {
		if dt, err := time.Parse(timeutil.SQLTimestampMinutes, s); err != nil {
			return s, err
		} else {
			return dt.Format(time.RFC3339), nil
		}
	}, true)
}

// UpdateDurations updates durations to canonical durations.
func UpdateDurations(t *table.Table, toTimeUnit time.Duration, prec int, addSuffix bool) error {
	if t == nil {
		return table.ErrTableCannotBeNil
	}
	if !addSuffix {
		if prec == 0 {
			t.FormatMap[-1] = table.FormatInt
		} else {
			t.FormatMap[-1] = table.FormatFloat
		}
	}
	return t.FormatColumns(1, -1, func(s string) (string, error) {
		s = strings.TrimSpace(s)
		if s == "" {
			return timeutil.DurationStringUnit(0, toTimeUnit, prec, addSuffix), nil
		} else if d, err := time.ParseDuration(s); err != nil {
			return s, err
		} else {
			return timeutil.DurationStringUnit(d, toTimeUnit, prec, addSuffix), nil
		}
	}, true)
}
