package timeseries

import "github.com/grokify/mogo/time/timeutil"

func (ts *TimeSeries) TimeUpdateIntervalStart() error {
	switch ts.Interval {
	case timeutil.Year:
		for rfc3339, ti := range ts.ItemMap {
			if !timeutil.IsYearStart(ti.Time) {
				delete(ts.ItemMap, rfc3339)
				ti.Time = timeutil.YearStart(ti.Time)
				ts.AddItems(ti)
			}
		}
		return nil
	}
	return ErrIntervalNotSupported
}
