package timeseries

import (
	"errors"

	"github.com/grokify/mogo/time/timeutil"
	"gonum.org/v1/gonum/stat"
)

// LinearRegressionYear only runs when interval=Year.
func (ts *TimeSeries) LinearRegressionYear() (alpha, beta float64, err error) {
	if ts.Interval != timeutil.Year {
		err = errors.New("timeseries = year is required")
	}

	xs, ys := []float64{}, []float64{}

	for _, ti := range ts.ItemMap {
		xs = append(xs, float64(ti.Time.Year()))
		ys = append(ys, ti.Float64())
	}
	alpha, beta = stat.LinearRegression(xs, ys, nil, false)
	return
}

// LinearRegressionYearProjection only runs when interval=Year for now.
// Use this to build tables.
func (ts *TimeSeries) LinearRegressionYearProjection(years uint, constantYOY bool) error {
	if years == 0 {
		return nil
	}
	alpha, beta, err := ts.LinearRegressionYear()
	if err != nil {
		return err
	}
	_, maxYear := ts.MinMaxTimes()
	maxYear = timeutil.YearStart(maxYear)
	yoyProjectionFirstYear := 0.0
	for i := 0; i < int(years); i++ {
		dtThis := maxYear.AddDate(i+1, 0, 0)
		tiPrev, err := ts.Get(dtThis.AddDate(-1, 0, 0))
		if err != nil {
			// times := ts.TimeSlice(true)
			// fmtutil.PrintJSON(timeslice.TimeSlice(times).Format(time.RFC3339))
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
