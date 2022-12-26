package timeseries

import (
	"os"
	"sort"

	"github.com/grokify/mogo/os/osutil"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/stringsutil"
)

type TimeSeriesSets struct {
	Name      string
	SetsMap   map[string]TimeSeriesSet
	Order     []string
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

func (sets *TimeSeriesSets) SetNames() []string {
	names := []string{}
	for name := range sets.SetsMap {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (sets *TimeSeriesSets) SeriesNames() []string {
	seriesNames := []string{}
	for _, set := range sets.SetsMap {
		for seriesName := range set.Series {
			seriesNames = append(seriesNames, seriesName)
		}
	}
	return stringsutil.SliceCondenseSpace(seriesNames, true, true)
}

// WriteJSON writes the data to a JSON file. To write a minimized JSON
// file use an empty string for `prefix` and `indent`.
func (sets *TimeSeriesSets) WriteJSON(filename string, perm os.FileMode, prefix, indent string) error {
	return osutil.WriteFileJSON(filename, sets, perm, prefix, indent)
}
