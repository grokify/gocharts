package statictimeseries

import (
	"sort"

	"github.com/grokify/simplego/time/timeutil"
)

type DataSeriesSet2 struct {
	Name     string
	SetsMap  map[string]DataSeriesSet
	Interval timeutil.Interval
}

func NewDataSeriesSet2(name string) DataSeriesSet2 {
	return DataSeriesSet2{
		Name:    name,
		SetsMap: map[string]DataSeriesSet{}}
}

func (dss2 *DataSeriesSet2) AddItems(items ...DataItem) {
	for _, item := range items {
		dss2.AddItem(item)
	}
}

func (dss2 *DataSeriesSet2) AddItem(item DataItem) {
	dss, ok := dss2.SetsMap[item.SeriesSetName]
	if !ok {
		dss = NewDataSeriesSet(item.SeriesSetName)
		dss.Interval = dss2.Interval
		if len(item.SeriesName) > 0 {
			dss.Name = item.SeriesName
		}
	}
	dss.AddItem(item)
	dss2.SetsMap[item.SeriesSetName] = dss
}

func (dss2 *DataSeriesSet2) SetNamesSorted() []string {
	names := []string{}
	for name := range dss2.SetsMap {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
