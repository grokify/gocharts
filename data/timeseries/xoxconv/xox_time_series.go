package xoxconv

/*
import (
	"time"

	"github.com/grokify/mogo/time/timeutil"

	"github.com/grokify/gocharts/v2/data/timeseries"
)

const (
	XoXClassNow = "Now"
	XoXClassMoM = "MoM"
	XoXClassQoQ = "QoQ"
	XoXClassYoY = "YoY"
)

// TimeSeriesXoX converts a `timeseries.TimeSeries` XoX statistics.
func TimeSeriesXoX(ts timeseries.TimeSeries) (XoXInfoMulti, error) {
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
*/
