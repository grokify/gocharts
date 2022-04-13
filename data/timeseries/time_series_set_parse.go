package timeseries

import (
	"fmt"
	"strconv"
	"time"

	"github.com/grokify/gocharts/v2/data/table"
)

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
	tss := NewTimeSeriesSet("")
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
			return tss, fmt.Errorf("cannot parse time [%s] in row [%d]", row[timeColIdx], i)
		}
		countString := row[countColIdx]
		if isFloat {
			countFloat, err := strconv.ParseFloat(countString, 64)
			if err != nil {
				return tss, fmt.Errorf("cannot parse count as float64 [%s] in row [%d]", row[timeColIdx], i)
			}
			tss.AddFloat64(seriesName, dt, countFloat)
		} else {
			countInt, err := strconv.Atoi(countString)
			if err != nil {
				return tss, fmt.Errorf("cannot parse count as int [%s] in row [%d]", row[timeColIdx], i)
			}
			tss.AddInt64(seriesName, dt, int64(countInt))
		}
	}
	return tss, nil
}

func ParseTimeFuncMonthYear(s string) (time.Time, error) {
	return time.Parse("January 2006", s)
}

func ParseTimeFuncRFC3339(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func ParseTimeFuncYearDotMonth(s string) (time.Time, error) {
	return time.Parse("2006.01", s)
}
