package timeseries

func (set *TimeSeriesSet) MapOfMapFloat64() map[string]map[string]float64 {
	m := map[string]map[string]float64{}
	for tsName, ts := range set.Series {
		m[tsName] = ts.MapFloat64()
	}
	return m
}

func (set *TimeSeriesSet) MapOfMapInt64() map[string]map[string]int64 {
	m := map[string]map[string]int64{}
	for tsName, ts := range set.Series {
		m[tsName] = ts.MapInt64()
	}
	return m
}
