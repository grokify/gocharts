package timeseries

import (
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/time/timeutil"
	"gonum.org/v1/gonum/stat"
)

// LinearRegression returns the `alpha` and `beta` for the data series.
// It currently only supports `month` and `year` time intervals. When
// `month` is used, the time is converted to `MonthContinuous` from
// `github.com/grokify/mogo/time/month`.
func (ts *TimeSeries) LinearRegression() (alpha, beta float64, err error) {
	xs, ys := []float64{}, []float64{}
	switch ts.Interval {
	case timeutil.IntervalYear:
		for _, ti := range ts.ItemMap {
			xs = append(xs, float64(ti.Time.Year()))
			ys = append(ys, ti.Float64())
		}
		alpha, beta = stat.LinearRegression(xs, ys, nil, false)
	case timeutil.IntervalMonth:
		for _, ti := range ts.ItemMap {
			xs = append(xs, float64(month.TimeToMonthContinuous(ti.Time)))
			ys = append(ys, ti.Float64())
		}
		alpha, beta = stat.LinearRegression(xs, ys, nil, false)
	default:
		err = ErrIntervalNotSupported
	}
	return
}

// LinearRegressionYearProjection only runs when interval=Year for now. Use this to build tables.
func (ts *TimeSeries) LinearRegressionYearProjection(years uint, constantYOY bool) error {
	if ts.Interval != timeutil.IntervalYear {
		return ErrIntervalNotSupported
	}
	if years == 0 {
		return nil
	}
	alpha, beta, err := ts.LinearRegression()
	if err != nil {
		return err
	}
	_, maxYear := ts.MinMaxTimes()
	maxYear = timeutil.NewTimeMore(maxYear, 0).YearStart()
	yoyProjectionFirstYear := 0.0
	for i := 0; i < int(years); i++ {
		dtThis := maxYear.AddDate(i+1, 0, 0)
		tiPrev, err := ts.Get(dtThis.AddDate(-1, 0, 0))
		if err != nil {
			// times := ts.TimeSlice(true)
			// fmtutil.PrintJSON(timeutil.Times(times).Format(time.RFC3339))
			// fmt.Printf("TPREV [%s]\n", dtThis.AddDate(-1, 0, 0).Format(time.RFC3339))
			panic("prev date do not, but should, exist")
		}
		if i == 0 {
			projectionThis := alpha + beta*float64(dtThis.Year())
			yoyProjectionFirstYear = (projectionThis - tiPrev.Float64()) / tiPrev.Float64()
			ts.AddFloat64(dtThis, projectionThis)
		} else if constantYOY {
			ts.AddFloat64(dtThis, tiPrev.Float64()*(1+yoyProjectionFirstYear))
		} else {
			ts.AddFloat64(dtThis, alpha+beta*float64(dtThis.Year()))
		}
	}
	return nil
}
