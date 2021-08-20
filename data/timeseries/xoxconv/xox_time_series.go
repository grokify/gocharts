package xoxconv

import (
	"time"

	"github.com/grokify/gocharts/data/timeseries"
	"github.com/grokify/simplego/time/timeutil"
)

// TimeSeriesXoX converts a `timeseries.TimeSeries`
// XoX statistics.
func TimeSeriesXoX(ts timeseries.TimeSeries) (XOXFloat64, error) {
	tsMonth := ts.ToMonth(true)

	xox := XOXFloat64{}
	tiNow, err := tsMonth.Last()
	if err != nil {
		return xox, err
	}
	xox.Now = Compare{
		XoX:   "NOW",
		Time:  tiNow.Time,
		Value: tiNow.Float64()}
	mago, err := tsMonth.Get(timeutil.TimeDt6SubNMonths(tiNow.Time, 1))
	if err == nil {
		xox.Month = Compare{
			XoX:   "MOM",
			Time:  mago.Time,
			Value: mago.Float64()}
		if mago.Float64() != 0 {
			xox.Month.Change = (tiNow.Float64() - mago.Float64()) / mago.Float64()
		}
	}
	qago, err := tsMonth.Get(timeutil.TimeDt6SubNMonths(tiNow.Time, 3))
	if err == nil {
		xox.Quarter = Compare{
			XoX:   "QOQ",
			Time:  qago.Time,
			Value: qago.Float64()}
		if qago.Float64() != 0 {
			xox.Quarter.Change = (tiNow.Float64() - qago.Float64()) / qago.Float64()
		}
	}
	yago, err := tsMonth.Get(timeutil.TimeDt6SubNMonths(tiNow.Time, 12))
	if err == nil {
		xox.Year = Compare{
			XoX:   "YOY",
			Time:  yago.Time,
			Value: yago.Float64()}
		if yago.Float64() != 0 {
			xox.Year.Change = (tiNow.Float64() - yago.Float64()) / yago.Float64()
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
