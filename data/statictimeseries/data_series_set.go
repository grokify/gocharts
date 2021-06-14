package statictimeseries

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/grokify/simplego/sort/sortutil"
	"github.com/grokify/simplego/time/month"
	"github.com/grokify/simplego/time/timeutil"
	"github.com/grokify/simplego/type/stringsutil"
	"github.com/pkg/errors"
)

type DataSeriesSet struct {
	Name     string
	Series   map[string]DataSeries
	Times    []time.Time
	Order    []string
	IsFloat  bool
	Interval timeutil.Interval
}

func NewDataSeriesSet(name string) DataSeriesSet {
	return DataSeriesSet{
		Name:   name,
		Series: map[string]DataSeries{},
		Times:  []time.Time{},
		Order:  []string{}}
}

func (set *DataSeriesSet) AddItems(items ...DataItem) {
	for _, item := range items {
		set.AddItem(item)
	}
}

func (set *DataSeriesSet) AddItem(item DataItem) {
	if _, ok := set.Series[item.SeriesName]; !ok {
		set.Series[item.SeriesName] =
			DataSeries{
				SeriesName: item.SeriesName,
				ItemMap:    map[string]DataItem{},
				IsFloat:    item.IsFloat,
				Interval:   set.Interval}
	}
	dataSeries := set.Series[item.SeriesName]
	dataSeries.AddItem(item)
	set.Series[item.SeriesName] = dataSeries
	set.Times = append(set.Times, item.Time)
}

func (set *DataSeriesSet) AddDataSeries(dataSeries ...DataSeries) error {
	for _, ds := range dataSeries {
		ds.SeriesName = strings.TrimSpace(ds.SeriesName)
		if len(ds.SeriesName) == 0 {
			return errors.New("E_DataSeriesSet.AddDataSeries_NO_DataSeries.SeriesName")
		}
		for _, item := range ds.ItemMap {
			if len(item.SeriesName) == 0 || item.SeriesName != ds.SeriesName {
				item.SeriesName = ds.SeriesName
			}
			set.AddItem(item)
		}
	}
	return nil
}

func (set *DataSeriesSet) Inflate() {
	set.Times = set.TimeSlice(true)
	if len(set.Order) > 0 {
		set.Order = stringsutil.SliceCondenseSpace(set.Order, true, false)
	} else {
		order := []string{}
		for name := range set.Series {
			order = append(order, name)
		}
		sort.Strings(order)
		set.Order = order
	}
}

func (set *DataSeriesSet) SeriesNames() []string {
	seriesNames := []string{}
	for seriesName := range set.Series {
		seriesNames = append(seriesNames, seriesName)
	}
	sort.Strings(seriesNames)
	return seriesNames
}

func (set *DataSeriesSet) GetSeriesByIndex(index int) (DataSeries, error) {
	if len(set.Order) == 0 && len(set.Series) > 0 {
		set.Inflate()
	}
	if index < len(set.Order) {
		name := set.Order[index]
		if ds, ok := set.Series[name]; ok {
			return ds, nil
		}
	}
	return DataSeries{}, fmt.Errorf("E_CANNOT_FIND_INDEX_[%d]_SET_COUNT_[%d]", index, len(set.Order))
}

func (set *DataSeriesSet) Item(seriesName, rfc3339 string) (DataItem, error) {
	di := DataItem{}
	dss, ok := set.Series[seriesName]
	if !ok {
		return di, fmt.Errorf("SeriesName not found [%s]", seriesName)
	}
	item, ok := dss.ItemMap[rfc3339]
	if !ok {
		return di, fmt.Errorf("SeriesName found [%s] Time not found [%s]", seriesName, rfc3339)
	}
	return item, nil
}

func (set *DataSeriesSet) TimeSlice(sortAsc bool) sortutil.TimeSlice {
	times := []time.Time{}
	for _, ds := range set.Series {
		for _, item := range ds.ItemMap {
			times = append(times, item.Time)
		}
	}
	times = timeutil.Sort(timeutil.Distinct(times))
	return month.TimeSeriesMonth(sortAsc, times...)
}

func (set *DataSeriesSet) TimeStrings() []string {
	times := []string{}
	for _, ds := range set.Series {
		for rfc3339 := range ds.ItemMap {
			times = append(times, rfc3339)
		}
	}
	return stringsutil.SliceCondenseSpace(times, true, true)
}

func (set *DataSeriesSet) MinMaxTimes() (time.Time, time.Time) {
	values := sortutil.TimeSlice{}
	for _, ds := range set.Series {
		min, max := ds.MinMaxTimes()
		values = append(values, min, max)
	}
	sort.Sort(values)
	return values[0], values[len(values)-1]
}

func (set *DataSeriesSet) MinMaxValues() (int64, int64) {
	values := sortutil.Int64Slice{}
	for _, ds := range set.Series {
		min, max := ds.MinMaxValues()
		values = append(values, min, max)
	}
	sort.Sort(values)
	return values[0], values[len(values)-1]
}

func (set *DataSeriesSet) MinMaxValuesFloat64() (float64, float64) {
	values := sort.Float64Slice{}
	for _, ds := range set.Series {
		min, max := ds.MinMaxValuesFloat64()
		values = append(values, min, max)
	}
	sort.Sort(values)
	return values[0], values[len(values)-1]
}

type RowInt64 struct {
	Name         string
	DisplayName  string
	HavePlusOne  bool
	ValuePlusOne int64
	Values       []int64
}

func (row *RowInt64) Flatten(conv func(v int64) string) []string {
	strs := []string{row.Name}
	for _, v := range row.Values {
		strs = append(strs, conv(v))
	}
	return strs
}

type RowFloat64 struct {
	Name   string
	Values []float64
}

func (row *RowFloat64) Flatten(conv func(v float64) string, preCount int, preVal string) []string {
	strs := []string{row.Name}
	for i := 0; i < preCount; i++ {
		strs = append(strs, preVal)
	}
	for _, v := range row.Values {
		strs = append(strs, conv(v))
	}
	return strs
}
