// statictimeseriesdata provides tools for adding and formatting
// static time series data for reporting purposes.
package statictimeseries

import (
	"fmt"
	"sort"
	"time"

	"github.com/grokify/gotilla/sort/sortutil"
	"github.com/grokify/gotilla/time/month"
	"github.com/grokify/gotilla/time/timeutil"
	"github.com/grokify/gotilla/type/maputil"
	"github.com/pkg/errors"
)

type DataSeries struct {
	SeriesName string
	ItemMap    map[string]DataItem
	IsFloat    bool
	Interval   timeutil.Interval // Informational
}

func NewDataSeries() DataSeries {
	return DataSeries{ItemMap: map[string]DataItem{}}
}

type DataItem struct {
	SeriesName string
	Time       time.Time
	IsFloat    bool
	Value      int64
	ValueFloat float64
}

// AddItem adds data item. It will sum values when
// existing time unit is encountered.
func (series *DataSeries) AddItem(item DataItem) {
	if series.ItemMap == nil {
		series.ItemMap = map[string]DataItem{}
	}
	if len(item.SeriesName) == 0 {
		item.SeriesName = series.SeriesName
	}
	item.Time = item.Time.UTC()
	rfc := item.Time.Format(time.RFC3339)
	if _, ok := series.ItemMap[rfc]; !ok {
		series.ItemMap[rfc] = item
	} else {
		existingItem := series.ItemMap[rfc]
		existingItem.Value += item.Value
		existingItem.ValueFloat += item.ValueFloat
		series.ItemMap[rfc] = existingItem
	}
}

func (series *DataSeries) SetSeriesName(seriesName string) {
	series.SeriesName = seriesName
	for k, v := range series.ItemMap {
		v.SeriesName = seriesName
		series.ItemMap[k] = v
	}
}

func (series *DataSeries) Keys() []string {
	keys := []string{}
	for key := range series.ItemMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (series *DataSeries) ItemsSorted() []DataItem {
	keys := series.Keys()
	items := []DataItem{}
	for _, key := range keys {
		item, ok := series.ItemMap[key]
		if !ok {
			panic(fmt.Sprintf("KEY_NOT_FOUND [%s]", key))
		}
		items = append(items, item)
	}
	return items
}

func (series *DataSeries) Last() (DataItem, error) {
	items := series.ItemsSorted()
	if len(items) == 0 {
		return DataItem{}, errors.New("E_NO_ITEMS")
	}
	return items[len(items)-1], nil
}

func (series *DataSeries) Pop() (DataItem, error) {
	items := series.ItemsSorted()
	if len(items) == 0 {
		return DataItem{}, errors.New("E_NO_ERROR")
	}
	last := items[len(items)-1]
	rfc := last.Time.Format(time.RFC3339)
	delete(series.ItemMap, rfc)
	return last, nil
}

func (series *DataSeries) minMaxValuesInt64Only() (int64, int64) {
	int64s := []int64{}
	for _, item := range series.ItemMap {
		int64s = append(int64s, item.Value)
	}
	if len(int64s) == 0 {
		return 0, 0
	}
	sort.Sort(sortutil.Int64Slice(int64s))
	return int64s[0], int64s[len(int64s)-1]
}

func (series *DataSeries) minMaxValuesFloat64Only() (float64, float64) {
	float64s := []float64{}
	for _, item := range series.ItemMap {
		float64s = append(float64s, item.ValueFloat)
	}
	if len(float64s) == 0 {
		return 0, 0
	}
	float64slice := sort.Float64Slice(float64s)
	sort.Sort(float64slice)

	return float64slice[0], float64slice[len(float64slice)-1]
}

func (series *DataSeries) MinMaxValues() (int64, int64) {
	if series.IsFloat {
		min, max := series.minMaxValuesFloat64Only()
		return int64(min), int64(max)
	}
	return series.minMaxValuesInt64Only()
}

func (series *DataSeries) MinMaxValuesFloat64() (float64, float64) {
	if series.IsFloat {
		return series.minMaxValuesFloat64Only()
	}
	min, max := series.minMaxValuesInt64Only()
	return float64(min), float64(max)
}

func (series *DataSeries) MinValue() int64 {
	min, _ := series.MinMaxValues()
	return min
}

func (series *DataSeries) MaxValue() int64 {
	_, max := series.MinMaxValues()
	return max
}

func (series *DataSeries) ToMonth() DataSeries {
	newDataSeries := DataSeries{
		SeriesName: series.SeriesName,
		ItemMap:    map[string]DataItem{},
		IsFloat:    series.IsFloat,
		Interval:   timeutil.Month}
	for _, item := range series.ItemMap {
		newDataSeries.AddItem(DataItem{
			SeriesName: item.SeriesName,
			Time:       month.MonthBegin(item.Time, 0),
			IsFloat:    item.IsFloat,
			Value:      item.Value,
			ValueFloat: item.ValueFloat})
	}
	return newDataSeries
}

func AggregateSeries(s1 DataSeries) DataSeries {
	aggregate := NewDataSeries()
	sortedItems := s1.SortedItems()
	sum := int64(0)
	for _, atomicItem := range sortedItems {
		aggregateItem := DataItem{
			SeriesName: atomicItem.SeriesName,
			Time:       atomicItem.Time,
			Value:      atomicItem.Value + sum,
		}
		sum = aggregateItem.Value
		aggregate.AddItem(aggregateItem)
	}
	return aggregate
}

// SortedItems returns sorted DataItems. This currently uses
// a simple string sort on RFC3339 times. For dates that are not
// handled properly this way, this can be enhanced to use more
// proper comparison
func (series *DataSeries) SortedItems() []DataItem {
	itemsSorted := []DataItem{}
	timesSorted := maputil.StringKeysSorted(series.ItemMap)
	for _, rfc3339 := range timesSorted {
		itemsSorted = append(itemsSorted, series.ItemMap[rfc3339])
	}
	return itemsSorted
}

func DataSeriesTimeSeries(series *DataSeries, interval timeutil.Interval) []time.Time {
	return timeutil.TimeSeriesSlice(interval, DataSeriesItemTimes(series))
}

func DataSeriesItemTimes(series *DataSeries) []time.Time {
	times := []time.Time{}
	for _, item := range series.ItemMap {
		times = append(times, item.Time)
	}
	return times
}

func DataSeriesMinMaxTimes(series *DataSeries) (time.Time, time.Time) {
	return timeutil.SliceMinMax(DataSeriesItemTimes(series))
}

func (series *DataSeries) MinMaxTimes() (time.Time, time.Time) {
	return DataSeriesMinMaxTimes(series)
}

func DataSeriesDivide(numer, denom DataSeries) (DataSeries, error) {
	denomKeys := denom.Keys()
	ds := NewDataSeries()
	ds.IsFloat = true
	if numer.Interval == denom.Interval {
		ds.Interval = denom.Interval
	}
	ds.SeriesName = numer.SeriesName + " / " + denom.SeriesName
	for _, dKey := range denomKeys {
		nItem, nOk := numer.ItemMap[dKey]
		dItem, dOk := denom.ItemMap[dKey]
		if !nOk && !dOk {
			continue
		} else if !dOk || dItem.Value == 0 {
			return ds, fmt.Errorf("E_DENOM_MISSING_OR_ZERO TIME [%s] NUMERATOR [%v]",
				dKey, nItem.Value)
		}
		ds.AddItem(DataItem{
			Time:       dItem.Time,
			ValueFloat: float64(nItem.Value) / float64(dItem.Value),
			IsFloat:    true,
		})
	}
	return ds, nil
}
