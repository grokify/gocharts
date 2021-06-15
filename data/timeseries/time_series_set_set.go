package timeseries

import (
	"sort"

	"github.com/grokify/simplego/time/timeutil"
)

type TimeSeriesSet2 struct {
	Name     string
	SetsMap  map[string]TimeSeriesSet
	Interval timeutil.Interval
}

func NewTimeSeriesSet2(name string) TimeSeriesSet2 {
	return TimeSeriesSet2{
		Name:    name,
		SetsMap: map[string]TimeSeriesSet{}}
}

func (tss2 *TimeSeriesSet2) AddItems(items ...TimeItem) {
	for _, item := range items {
		tss2.AddItem(item)
	}
}

func (tss2 *TimeSeriesSet2) AddItem(item TimeItem) {
	dss, ok := tss2.SetsMap[item.SeriesSetName]
	if !ok {
		dss = NewTimeSeriesSet(item.SeriesSetName)
		dss.Interval = tss2.Interval
		if len(item.SeriesName) > 0 {
			dss.Name = item.SeriesName
		}
	}
	dss.AddItem(item)
	tss2.SetsMap[item.SeriesSetName] = dss
}

func (tss2 *TimeSeriesSet2) SetNamesSorted() []string {
	names := []string{}
	for name := range tss2.SetsMap {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
