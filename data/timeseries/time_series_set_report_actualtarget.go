package timeseries

import (
	"errors"
	"strconv"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/time/year"
)

func (set *TimeSeriesSet) TableActualTarget() (*table.Table, error) {
	var times timeutil.Times
	set.Inflate()
	tbl := table.NewTable("")
	tbl.Columns = append(tbl.Columns, "Series")
	tbl.FormatMap = map[int]string{0: table.FormatString, -1: table.FormatFloat}
	if set.Interval == timeutil.IntervalYear {
		times = year.TimesYearStarts(set.Times...)
		for _, dt := range times {
			tbl.Columns = append(tbl.Columns, strconv.Itoa(dt.Year()))
		}
	} else if set.Interval == timeutil.IntervalMonth {
		times = month.TimesMonthStarts(set.Times...)
		for _, dt := range times {
			tbl.Columns = append(tbl.Columns, dt.Format("Jan 2006"))
		}
	} else {
		return nil, ErrIntervalNotSupported
	}
	for _, pair := range set.ActualTargetPairs {
		actualTS, ok := set.Series[pair.ActualSeriesName]
		if !ok {
			return nil, errors.New("actual timeseries not found")
		}
		targetTS, ok := set.Series[pair.TargetSeriesName]
		if !ok {
			return nil, errors.New("target timeseries not found")
		}
		rowActual := []string{actualTS.SeriesName}
		rowTarget := []string{targetTS.SeriesName}
		rowDiff := []string{actualTS.SeriesName + " vs. " + targetTS.SeriesName}
		for _, dt := range times {
			itemActual, errActual := actualTS.Get(dt)
			if errActual != nil {
				rowActual = append(rowActual, "0")
			} else {
				rowActual = append(rowActual, strconvutil.Ftoa(itemActual.Float64(), -1))
			}
			itemTarget, errTarget := targetTS.Get(dt)
			if errTarget != nil {
				rowTarget = append(rowTarget, "0")
			} else {
				rowTarget = append(rowTarget, strconvutil.Ftoa(itemTarget.Float64(), -1))
			}
			if errActual != nil || errTarget != nil {
				rowDiff = append(rowDiff, "0")
			} else {
				rowDiff = append(rowDiff, strconvutil.Ftoa((itemActual.Float64()-itemTarget.Float64())/itemTarget.Float64(), -1))
			}
		}
		tbl.Rows = append(tbl.Rows, rowActual, rowTarget, rowDiff)
	}
	return &tbl, nil
}
