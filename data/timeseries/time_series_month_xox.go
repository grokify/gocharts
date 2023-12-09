package timeseries

import (
	"errors"
	"strings"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

const (
	XoXClassNow = "Now"
	XoXClassMoM = "MoM"
	XoXClassQoQ = "QoQ"
	XoXClassYoY = "YoY"
)

func (ts *TimeSeries) TimeSeriesMonthYOY() TimeSeries {
	return ts.TimeSeriesMonthXOX(-1, 0, 0, XoXClassYoY)
}

func (ts *TimeSeries) TimeSeriesMonthQOQ() TimeSeries {
	return ts.TimeSeriesMonthXOX(0, -3, 0, XoXClassQoQ)
}

func (ts *TimeSeries) TimeSeriesMonthMOM() TimeSeries {
	return ts.TimeSeriesMonthXOX(0, -1, 0, XoXClassMoM)
}

func (ts *TimeSeries) TimeSeriesMonthXOX(years, months, days int, suffix string) TimeSeries {
	tsm := ts.ToMonth(true)
	tsXOX := NewTimeSeries(tsm.SeriesName)
	suffix = strings.TrimSpace(suffix)
	if len(suffix) > 0 {
		if len(tsm.SeriesName) > 0 {
			tsXOX.SeriesName += " " + suffix
		} else {
			tsXOX.SeriesName = suffix
		}
	}
	tsXOX.IsFloat = true
	tsXOX.Interval = tsm.Interval
	times := tsm.Times(true)
	times.SortReverse()

	for _, dt := range times {
		dtThis := dt
		dtPast := dt.AddDate(years, months, days)
		tiThis, err := tsm.Get(dtThis)
		if err != nil {
			panic("cannot find this time")
		}
		tiPast, err := ts.Get(dtPast)
		if err != nil {
			continue
		}
		tsXOX.AddFloat64(dtThis, (tiThis.Float64()-tiPast.Float64())/tiPast.Float64())
	}

	return tsXOX
}

func (ts *TimeSeries) TimeSeriesYearYOY(suffix string) (TimeSeries, error) {
	if ts.Interval != timeutil.IntervalYear {
		return TimeSeries{}, errors.New("interval year is required")
	}
	// tsm := ts.ToMonth(true)
	tsYOY := NewTimeSeries(ts.SeriesName)
	suffix = strings.TrimSpace(suffix)
	if len(suffix) > 0 {
		if len(ts.SeriesName) > 0 {
			tsYOY.SeriesName += " " + suffix
		} else {
			tsYOY.SeriesName = suffix
		}
	}
	tsYOY.IsFloat = true
	tsYOY.Interval = ts.Interval
	times := ts.Times(true)
	times.SortReverse()

	for _, dt := range times {
		dtThis := dt
		dtPast := dt.AddDate(-1, 0, 0)
		tiThis, err := ts.Get(dtThis)
		if err != nil {
			panic("cannot find this time")
		}
		tiPast, err := ts.Get(dtPast)
		if err != nil {
			continue
		}
		tsYOY.AddFloat64(dtThis, (tiThis.Float64()-tiPast.Float64())/tiPast.Float64())
	}

	return tsYOY, nil
}

func (ts *TimeSeries) XoXInfoMulti() (XoXInfoMulti, error) {
	tsMonth := ts.ToMonth(true)

	xox := XoXInfoMulti{}
	tiNow, err := tsMonth.Last()
	if err != nil {
		return xox, err
	}
	xox.Now = XoXInfo{
		XoX:   XoXClassNow,
		Time:  tiNow.Time,
		Value: tiNow.Float64()}
	mago, err := tsMonth.Get(timeutil.TimeDT6AddNMonths(tiNow.Time, -1))
	if err == nil {
		xox.Month = XoXInfo{
			XoX:   XoXClassMoM,
			Time:  mago.Time,
			Value: mago.Float64()}
		if mago.Float64() != 0 {
			xox.Month.Change = (tiNow.Float64() - mago.Float64()) / mago.Float64()
		}
	}
	qago, err := tsMonth.Get(timeutil.TimeDT6AddNMonths(tiNow.Time, -3))
	if err == nil {
		xox.Quarter = XoXInfo{
			XoX:   XoXClassQoQ,
			Time:  qago.Time,
			Value: qago.Float64()}
		if qago.Float64() != 0 {
			xox.Quarter.Change = (tiNow.Float64() - qago.Float64()) / qago.Float64()
		}
	}
	yago, err := tsMonth.Get(timeutil.TimeDT6AddNMonths(tiNow.Time, -12))
	if err == nil {
		xox.Year = XoXInfo{
			XoX:   XoXClassYoY,
			Time:  yago.Time,
			Value: yago.Float64()}
		if yago.Float64() != 0 {
			xox.Year.Change = (tiNow.Float64() - yago.Float64()) / yago.Float64()
		}
	}
	return xox, nil
}

type XoXInfoMulti struct {
	Now     XoXInfo
	Month   XoXInfo
	Quarter XoXInfo
	Year    XoXInfo
}

type XoXInfo struct {
	Time   time.Time
	Value  float64
	Change float64
	XoX    string
}
