// statictimeseriesdata provides tools for adding and formatting
// static time series data for reporting purposes.
package statictimeseries

import (
	"sort"
)

type DataSeriesSetSimpleSet struct {
	Name    string
	SetsMap map[string]DataSeriesSetSimple
}

func NewDataSeriesSetSimpleSet(name string) DataSeriesSetSimpleSet {
	return DataSeriesSetSimpleSet{
		Name:    name,
		SetsMap: map[string]DataSeriesSetSimple{}}
}

func (ds3set *DataSeriesSetSimpleSet) AddItem(setName string, item DataItem) {
	ds3, ok := ds3set.SetsMap[setName]
	if !ok {
		ds3 = NewDataSeriesSetSimple()
	}
	ds3.AddItem(item)
	ds3set.SetsMap[setName] = ds3
}

func (ds3set *DataSeriesSetSimpleSet) SetNamesSorted() []string {
	names := []string{}
	for name := range ds3set.SetsMap {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
