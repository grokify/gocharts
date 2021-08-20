package xoxconv

import (
	"time"

	"github.com/grokify/gocharts/data/timeseries"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/time/timeutil"
)

// TimeSeriesXoX converts a `timeseries.TimeSeries`
// XoX statistics.
func TimeSeriesXoX(ts timeseries.TimeSeries) (XOXFloat64, error) {
	tsMonth := ts.ToMonth(true)
	fmtutil.PrintJSON(tsMonth)

	xox := XOXFloat64{}
	tiNow, err := tsMonth.Last()
	if err != nil {
		return xox, err
	}
	xox.Now = Compare{
		XoX:   "NOW",
		Time:  tiNow.Time,
		Value: tiNow.ValueFloat64()}
	mago, err := tsMonth.Get(timeutil.TimeDt6SubNMonths(tiNow.Time, 1))
	if err == nil {
		xox.Month = Compare{
			XoX:   "MOM",
			Time:  mago.Time,
			Value: mago.ValueFloat64()}
		if mago.ValueFloat64() != 0 {
			xox.Month.Change = (tiNow.ValueFloat64() - mago.ValueFloat64()) / mago.ValueFloat64()
		}
	}
	qago, err := tsMonth.Get(timeutil.TimeDt6SubNMonths(tiNow.Time, 3))
	if err == nil {
		xox.Quarter = Compare{
			XoX:   "QOQ",
			Time:  qago.Time,
			Value: qago.ValueFloat64()}
		if qago.ValueFloat64() != 0 {
			xox.Quarter.Change = (tiNow.ValueFloat64() - qago.ValueFloat64()) / qago.ValueFloat64()
		}
	}
	yago, err := tsMonth.Get(timeutil.TimeDt6SubNMonths(tiNow.Time, 12))
	if err == nil {
		xox.Year = Compare{
			XoX:   "YOY",
			Time:  yago.Time,
			Value: yago.ValueFloat64()}
		if yago.ValueFloat64() != 0 {
			xox.Year.Change = (tiNow.ValueFloat64() - yago.ValueFloat64()) / yago.ValueFloat64()
		}
	}
	return xox, nil
}

type XOXFloat64 struct {
	Now     Compare
	Month   Compare
	Quarter Compare
	Year    Compare
}

type Compare struct {
	Time   time.Time
	Value  float64
	Change float64
	XoX    string
}
