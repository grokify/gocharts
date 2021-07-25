package timeseries

import (
	"sort"

	"github.com/grokify/simplego/time/timeutil"
)

type TimeSeriesSets struct {
	Name      string
	SetsMap   map[string]TimeSeriesSet
	KeyIsTime bool
	Interval  timeutil.Interval
}

func NewTimeSeriesSets(name string) TimeSeriesSets {
	return TimeSeriesSets{
		Name:    name,
		SetsMap: map[string]TimeSeriesSet{}}
}

func (sets *TimeSeriesSets) AddItems(items ...TimeItem) {
	for _, item := range items {
		set, ok := sets.SetsMap[item.SeriesSetName]
		if !ok {
			set = NewTimeSeriesSet(item.SeriesSetName)
			set.Interval = sets.Interval
			if len(item.SeriesName) > 0 {
				set.Name = item.SeriesName
			}
		}
		set.AddItems(item)
		sets.SetsMap[item.SeriesSetName] = set
	}
}

func (sets *TimeSeriesSets) SetNamesSorted() []string {
	names := []string{}
	for name := range sets.SetsMap {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
