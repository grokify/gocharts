package timeseries

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/time/timeutil"
)

// ParseTableTimeSeriesSetMatrixRows creates a `TimeSeriesSet` from a `table.Table` where the
// various time series are in rows and the time intervals are columns `1:`. `funcStringToTime`
// can be used to convert column names to times and set to `timeutil.ParseTimeCanonicalFunc("Jan 2006")`.
// `funcStringToInt` can be set to `strconvutil.AtoiMoreFunc(",", ".")`.
func ParseTableTimeSeriesSetMatrixRows(tbl table.Table, interval timeutil.Interval, isFloat bool,
	funcStringToTime func(s string) (time.Time, error),
	funcStringToInt func(s string) (int, error),
	funcStringToFloat func(s string) (float64, error),
) (*TimeSeriesSet, error) {
	tss := NewTimeSeriesSet("")
	tss.IsFloat = isFloat
	if len(tbl.Columns) < 2 {
		return nil, errors.New("columns length cannot be less than 2")
	}
	colTimes, err := strconvutil.SliceAtotFunc(funcStringToTime, tbl.Columns[1:])
	if err != nil {
		return nil, err
	}
	for _, row := range tbl.Rows {
		if len(row) == 0 {
			continue
		} else if len(row) < 2 {
			return nil, errors.New("row length cannot be less than 2")
		} else if len(row) > len(tbl.Columns) {
			return nil, errors.New("row is longer than columns")
		}
		seriesName := row[0]
		for i := 1; i < len(row); i++ {
			dt := colTimes[i-1]
			countString := strings.TrimSpace(row[i])
			if countString == "" || countString == "0" {
				tss.AddInt64(seriesName, dt, 0)
			} else if isFloat {
				if floatVal, err := strconvutil.AtofFunc(funcStringToFloat, countString); err != nil {
					return nil, err
				} else {
					tss.AddFloat64(seriesName, dt, floatVal)
				}
			} else {
				if intVal, err := strconvutil.AtoiFunc(funcStringToInt, countString); err != nil {
					return nil, err
				} else {
					tss.AddInt64(seriesName, dt, int64(intVal))
				}
			}
		}
	}
	return &tss, nil
}

// ParseTableTimeSeriesSetMatrix create a `TimeSeriesSet` from a `table.Table` using the least
// amount of input to populate the data structure. The time must be in column 0 and the series
// names must be in the column headers.
func ParseTableTimeSeriesSetMatrix(tbl table.Table, isFloat bool, timeParseFunc func(s string) (time.Time, error)) (TimeSeriesSet, error) {
	if timeParseFunc == nil {
		timeParseFunc = ParseTimeFuncRFC3339
	}
	tss := NewTimeSeriesSet("")
	tss.IsFloat = isFloat
	for y, row := range tbl.Rows {
		if len(row) <= 1 {
			continue
		}
		dt, err := timeParseFunc(row[0])
		if err != nil {
			return tss, fmt.Errorf("cannot parse time [%s] in row [%d]", row[0], y)
		}
		for x := 1; x < len(row); x++ {
			if x >= len(tbl.Columns) {
				return tss, fmt.Errorf("no column header for column [%d] on row [%d]", x, y)
			}
			seriesName := tbl.Columns[x]
			countString := row[x]
			if isFloat {
				countFloat, err := strconv.ParseFloat(countString, 64)
				if err != nil {
					return tss, fmt.Errorf("cannot parse count as float64 [%s] in row [%d]", row[x], y)
				}
				tss.AddFloat64(seriesName, dt, countFloat)
			} else {
				countInt, err := strconv.Atoi(countString)
				if err != nil {
					return tss, fmt.Errorf("cannot parse count as int [%s] in row [%d]", row[x], y)
				}
				tss.AddInt64(seriesName, dt, int64(countInt))
			}
		}
	}
	tss.Times = tss.TimeSlice(true)
	return tss, nil
}

// ParseTableTimeSeriesSetFlat create a `TimeSeriesSet` from a `table.Table` using the least
// amount of input to populate the data structure. It does not set the following
// parameters which must be set manually: `Name`, `Interval`.
func ParseTableTimeSeriesSetFlat(tbl table.Table, timeColIdx, seriesNameColIdx, countColIdx uint, isFloat bool, timeParseFunc func(s string) (time.Time, error)) (TimeSeriesSet, error) {
	if timeParseFunc == nil {
		timeParseFunc = ParseTimeFuncRFC3339
	}
	tss := NewTimeSeriesSet(tbl.Name)
	tss.IsFloat = isFloat
	for i, row := range tbl.Rows {
		if int(seriesNameColIdx) >= len(row) {
			return tss, fmt.Errorf("colIdx [%d] not present in row [%d]", seriesNameColIdx, i)
		}
		if int(timeColIdx) >= len(row) {
			return tss, fmt.Errorf("colIdx [%d] not present in row [%d]", timeColIdx, i)
		}
		if int(countColIdx) >= len(row) {
			return tss, fmt.Errorf("colIdx [%d] not present in row [%d]", countColIdx, i)
		}
		seriesName := row[seriesNameColIdx]

		dt, err := timeParseFunc(row[timeColIdx])
		if err != nil {
			return tss, fmt.Errorf("cannot parse time [%s] in row [%d] err [%s]", row[timeColIdx], i, err.Error())
		}
		countString := row[countColIdx]
		if isFloat {
			countFloat, err := strconv.ParseFloat(countString, 64)
			if err != nil {
				return tss, fmt.Errorf("cannot parse count as float64 [%s] in row [%d] err [%s]", row[countColIdx], i, err.Error())
			}
			tss.AddFloat64(seriesName, dt, countFloat)
		} else {
			countInt, err := strconv.Atoi(countString)
			if err != nil {
				return tss, fmt.Errorf("cannot parse count as int [%s] in row [%d] err [%s]", row[countColIdx], i, err.Error())
			}
			tss.AddInt64(seriesName, dt, int64(countInt))
		}
	}
	return tss, nil
}

func ParseTimeFuncMonthYear(s string) (time.Time, error) {
	return timeutil.ParseTimeCanonical("Jan 2006", s)
}

func ParseTimeFuncRFC3339(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func ParseTimeFuncYearDotMonth(s string) (time.Time, error) {
	return time.Parse("2006.01", s)
}
