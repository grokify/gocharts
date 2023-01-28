package timeseries

import "github.com/grokify/mogo/time/timeutil"

func (ts *TimeSeries) TimeUpdateIntervalStart() error {
	switch ts.Interval {
	case timeutil.Year:
		for rfc3339, ti := range ts.ItemMap {
			tm := timeutil.NewTimeMore(ti.Time, 0)
			if !tm.IsYearStart() {
				delete(ts.ItemMap, rfc3339)
				ti.Time = tm.YearStart()
				ts.AddItems(ti)
			}
		}
		return nil
	case timeutil.Month:
		for rfc3339, ti := range ts.ItemMap {
			tm := timeutil.NewTimeMore(ti.Time, 0)
			if !tm.IsMonthStart() {
				delete(ts.ItemMap, rfc3339)
				ti.Time = tm.MonthStart()
				ts.AddItems(ti)
			}
		}
		return nil
	}
	return ErrIntervalNotSupported
}
