package timeseries

import (
	"time"

	"github.com/grokify/simplego/time/timeutil"
)

func (set *TimeSeriesSet) ToMonth(inflate bool) TimeSeriesSet {
	newDss := TimeSeriesSet{
		Name:     set.Name,
		Series:   map[string]TimeSeries{},
		Times:    set.Times,
		Interval: timeutil.Month,
		Order:    set.Order}
	for name, ds := range set.Series {
		newDss.Series[name] = ds.ToMonth(inflate)
	}
	newDss.Times = newDss.TimeSlice(true)
	return newDss
}

func (set *TimeSeriesSet) ToMonthCumulative(popLast, inflate bool) (TimeSeriesSet, error) {
	newDss := TimeSeriesSet{
		Name:     set.Name,
		Series:   map[string]TimeSeries{},
		Times:    set.Times,
		Interval: timeutil.Month,
		Order:    set.Order}
	for name, ds := range set.Series {
		newDs, err := ds.ToMonthCumulative(inflate, newDss.Times...)
		if err != nil {
			return newDss, err
		}
		newDss.Series[name] = newDs
	}
	if popLast {
		newDss.PopLast()
	}
	newDss.Times = newDss.TimeSlice(true)
	return newDss, nil
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
	newDss := TimeSeriesSet{
		Name:     set.Name,
		Series:   map[string]TimeSeries{},
		Times:    set.Times,
		IsFloat:  set.IsFloat,
		Interval: timeutil.Month,
		Order:    []string{}}
	for _, ds := range set.Series {
		for _, item := range ds.ItemMap {
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
			newDss.AddItem(item)
		}
	}
	newDss.Inflate()
	return newDss
}
