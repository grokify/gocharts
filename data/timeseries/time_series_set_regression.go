package timeseries

// LinearRegressionYearProjection only runs when interval=Year.
func (set *TimeSeriesSet) LinearRegressionYearProjection(years uint32, constantYOY bool) error {
	for name, ts := range set.Series {
		if err := ts.LinearRegressionYearProjection(years, constantYOY); err != nil {
			return err
		} else {
			set.Series[name] = ts
		}
	}
	return nil
}
