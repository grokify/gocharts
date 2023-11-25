package timeseries

import (
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

// ToYear aggregates time values into months. `inflate` is used to add months with `0` values.
func (set *TimeSeriesSet) ToYear(inflate, popLast bool) (TimeSeriesSet, error) {
	newTSS := TimeSeriesSet{
		Name:     set.Name,
		Series:   map[string]TimeSeries{},
		Interval: timeutil.IntervalYear,
		Order:    set.Order}
	for name, ts := range set.Series {
		newTSS.Series[name] = ts.ToYear()
	}
	if popLast {
		newTSS.PopLast()
	}
	newTSS.Inflate()
	newTSS.Times = newTSS.TimeSlice(true)
	return newTSS, nil
}

// ToMonth aggregates time values into months. `inflate` is used to add months with `0` values.
func (set *TimeSeriesSet) ToMonth(cumulative, inflate, popLast bool, monthsFilter []time.Month) (TimeSeriesSet, error) {
	if cumulative {
		return set.toMonthCumulative(inflate, popLast)
	}
	newTSS := TimeSeriesSet{
		Name:     set.Name,
		Series:   map[string]TimeSeries{},
		Times:    set.Times,
		Interval: timeutil.IntervalMonth,
		Order:    set.Order}
	for name, ts := range set.Series {
		newTSS.Series[name] = ts.ToMonth(inflate, monthsFilter...)
	}
	if popLast {
		newTSS.PopLast()
	}
	newTSS.Inflate()
	newTSS.Times = newTSS.TimeSlice(true)
	return newTSS, nil
}

func (set *TimeSeriesSet) toMonthCumulative(inflate, popLast bool) (TimeSeriesSet, error) {
	newTSS := TimeSeriesSet{
		Name:     set.Name,
		Series:   map[string]TimeSeries{},
		Times:    set.Times,
		Interval: timeutil.IntervalMonth,
		Order:    set.Order}
	for seriesName, ts := range set.Series {
		newTS, err := ts.ToMonthCumulative(inflate, newTSS.Times...)
		if err != nil {
			return newTSS, err
		}
		newTSS.Series[seriesName] = newTS
	}
	if popLast {
		newTSS.PopLast()
	}
	newTSS.Inflate()
	newTSS.Times = newTSS.TimeSlice(true)
	return newTSS, nil
}

func (set *TimeSeriesSet) PopLast() {
	times := set.TimeSlice(true)
	if len(times) == 0 {
		return
	}
	last := times[len(times)-1]
	set.DeleteTime(last)
}

func (set *TimeSeriesSet) DeleteTime(dt time.Time) {
	for id, ds := range set.Series {
		ds.DeleteTime(dt)
		set.Series[id] = ds
	}
}

func (set *TimeSeriesSet) ToNewSeriesNames(seriesNames, seriesSetNames map[string]string) TimeSeriesSet {
	newTSS := TimeSeriesSet{
		Name:     set.Name,
		Series:   map[string]TimeSeries{},
		Times:    set.Times,
		IsFloat:  set.IsFloat,
		Interval: timeutil.IntervalMonth,
		Order:    []string{}}
	for _, ts := range set.Series {
		for _, item := range ts.ItemMap {
			if len(seriesNames) > 0 {
				if newSeriesName, ok := seriesNames[item.SeriesName]; ok {
					item.SeriesName = newSeriesName
				}
			}
			if len(seriesSetNames) > 0 {
				if newSeriesSetName, ok := seriesSetNames[item.SeriesName]; ok {
					item.SeriesSetName = newSeriesSetName
				}
			}
			newTSS.AddItems(item)
		}
	}
	newTSS.Inflate()
	return newTSS
}
