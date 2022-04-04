package timeseries

import "strings"

func (ts *TimeSeries) TimeSeriesMonthYOY() TimeSeries {
	return ts.TimeSeriesMonthXOX(-1, 0, 0, "YoY")
}

func (ts *TimeSeries) TimeSeriesMonthQOQ() TimeSeries {
	return ts.TimeSeriesMonthXOX(0, -3, 0, "QoQ")
}

func (ts *TimeSeries) TimeSeriesMonthMOM() TimeSeries {
	return ts.TimeSeriesMonthXOX(0, -1, 0, "MoM")
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
	times := tsm.TimeSlice(true)
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
