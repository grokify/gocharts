package timeseries

// LinearRegressionYearProjection only runs when interval=Year.
func (set *TimeSeriesSet) LinearRegressionYearProjection(years uint, constantYOY bool) error {
	for name, ts := range set.Series {
		err := ts.LinearRegressionYearProjection(years, constantYOY)
		if err != nil {
			return err
		}
		set.Series[name] = ts
	}
	return nil
}
