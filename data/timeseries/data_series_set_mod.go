package timeseries

import (
	"time"

	"github.com/grokify/simplego/time/timeutil"
)

func (set *DataSeriesSet) ToMonth(inflate bool) DataSeriesSet {
	newDss := DataSeriesSet{
		Name:     set.Name,
		Series:   map[string]DataSeries{},
		Times:    set.Times,
		Interval: timeutil.Month,
		Order:    set.Order}
	for name, ds := range set.Series {
		newDss.Series[name] = ds.ToMonth(inflate)
	}
	newDss.Times = newDss.TimeSlice(true)
	return newDss
}

func (set *DataSeriesSet) ToMonthCumulative(popLast, inflate bool) (DataSeriesSet, error) {
	newDss := DataSeriesSet{
		Name:     set.Name,
		Series:   map[string]DataSeries{},
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

func (set *DataSeriesSet) PopLast() {
	times := set.TimeSlice(true)
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

func (set *DataSeriesSet) ToNewSeriesNames(seriesNames, seriesSetNames map[string]string) DataSeriesSet {
	newDss := DataSeriesSet{
		Name:     set.Name,
		Series:   map[string]DataSeries{},
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
