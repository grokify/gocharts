package timeseries

import (
	"time"

	"github.com/grokify/simplego/time/timeutil"
)

func (set *TimeSeriesSet) ToMonth(inflate bool) TimeSeriesSet {
	newTss := TimeSeriesSet{
		Name:     set.Name,
		Series:   map[string]TimeSeries{},
		Times:    set.Times,
		Interval: timeutil.Month,
		Order:    set.Order}
	for name, ds := range set.Series {
		newTss.Series[name] = ds.ToMonth(inflate)
	}
	newTss.Times = newTss.TimeSlice(true)
	return newTss
}

func (set *TimeSeriesSet) ToMonthCumulative(popLast, inflate bool) (TimeSeriesSet, error) {
	newTss := TimeSeriesSet{
		Name:     set.Name,
		Series:   map[string]TimeSeries{},
		Times:    set.Times,
		Interval: timeutil.Month,
		Order:    set.Order}
	for name, ds := range set.Series {
		newDs, err := ds.ToMonthCumulative(inflate, newTss.Times...)
		if err != nil {
			return newTss, err
		}
		newTss.Series[name] = newDs
	}
	if popLast {
		newTss.PopLast()
	}
	newTss.Times = newTss.TimeSlice(true)
	return newTss, nil
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
	newTss := TimeSeriesSet{
		Name:     set.Name,
		Series:   map[string]TimeSeries{},
		Times:    set.Times,
		IsFloat:  set.IsFloat,
		Interval: timeutil.Month,
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
			newTss.AddItems(item)
		}
	}
	newTss.Inflate()
	return newTss
}
