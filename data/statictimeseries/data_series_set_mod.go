package statictimeseries

import (
	"time"

	"github.com/grokify/simplego/time/timeutil"
)

func (set *DataSeriesSet) ToMonth() DataSeriesSet {
	newDss := DataSeriesSet{
		Name:     set.Name,
		Series:   map[string]DataSeries{},
		Times:    set.Times,
		Interval: timeutil.Month,
		Order:    set.Order}
	for name, ds := range set.Series {
		newDss.Series[name] = ds.ToMonth()
	}
	newDss.Times = newDss.GetTimeSlice(true)
	return newDss
}

func (set *DataSeriesSet) ToMonthCumulative(popLast bool) (DataSeriesSet, error) {
	newDss := DataSeriesSet{
		Name:     set.Name,
		Series:   map[string]DataSeries{},
		Times:    set.Times,
		Interval: timeutil.Month,
		Order:    set.Order}
	for name, ds := range set.Series {
		newDs, err := ds.ToMonthCumulative(newDss.Times...)
		if err != nil {
			return newDss, err
		}
		newDss.Series[name] = newDs
	}
	if popLast {
		newDss.PopLast()
	}
	newDss.Times = newDss.GetTimeSlice(true)
	return newDss, nil
}

func (set *DataSeriesSet) PopLast() {
	times := set.GetTimeSlice(true)
	if len(times) == 0 {
		return
	}
	last := times[len(times)-1]
	set.DeleteItemByTime(last)
}

func (set *DataSeriesSet) DeleteItemByTime(dt time.Time) {
	for id, ds := range set.Series {
		ds.DeleteByTime(dt)
		set.Series[id] = ds
	}
}
