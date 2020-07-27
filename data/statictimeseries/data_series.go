// statictimeseriesdata provides tools for adding and formatting
// static time series data for reporting purposes.
package statictimeseries

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gotilla/sort/sortutil"
	"github.com/grokify/gotilla/time/month"
	"github.com/grokify/gotilla/time/timeutil"
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

func (item *DataItem) ValueInt64() int64 {
	if item.IsFloat {
		return int64(item.ValueFloat)
	}
	return item.Value
}

func (item *DataItem) ValueFloat64() float64 {
	if item.IsFloat {
		return item.ValueFloat
	}
	return float64(item.Value)
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

// Keys returns a sorted listed of Item keys.
func (series *DataSeries) Keys() []string {
	// maputil.StringKeysSorted(series.ItemMap)
	keys := []string{}
	for key := range series.ItemMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// ItemsSorted returns sorted DataItems. This currently uses
// a simple string sort on RFC3339 times.
func (series *DataSeries) ItemsSorted() []DataItem {
	keys := series.Keys()
	itemsSorted := []DataItem{}
	for _, key := range keys {
		item, ok := series.ItemMap[key]
		if !ok {
			panic(fmt.Sprintf("KEY_NOT_FOUND [%s]", key))
		}
		itemsSorted = append(itemsSorted, item)
	}
	return itemsSorted
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

func (ds *DataSeries) GetTimeSlice(sortSlice bool) timeutil.TimeSlice {
	times := timeutil.TimeSlice{}
	for _, item := range ds.ItemMap {
		times = append(times, item.Time)
	}
	if sortSlice {
		sort.Sort(times)
	}
	return times
}

func (ds *DataSeries) DeleteByTime(dt time.Time) {
	delete(ds.ItemMap, dt.Format(time.RFC3339))
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

func (ds *DataSeries) ToMonthCumulative(timesInput ...time.Time) (DataSeries, error) {
	newDataSeries := DataSeries{
		SeriesName: ds.SeriesName,
		ItemMap:    map[string]DataItem{},
		IsFloat:    ds.IsFloat,
		Interval:   timeutil.Month}
	dsMonth := ds.ToMonth()
	var min time.Time
	var max time.Time
	var err error
	if len(timesInput) > 0 {
		min, max, err = timeutil.TimeSliceMinMax(timesInput)
		if err != nil {
			return newDataSeries, err
		}
	} else {
		min, max, err = timeutil.TimeSliceMinMax(dsMonth.GetTimeSlice(false))
		if err != nil {
			return newDataSeries, err
		}
	}
	times := timeutil.TimeSeriesSlice(timeutil.Month, []time.Time{min, max})
	cItems := []DataItem{}
	for _, t := range times {
		rfc := t.Format(time.RFC3339)
		if item, ok := dsMonth.ItemMap[rfc]; ok {
			if len(cItems) > 0 {
				prevCItem := cItems[len(cItems)-1]
				cItems = append(cItems, DataItem{
					SeriesName: newDataSeries.SeriesName,
					IsFloat:    newDataSeries.IsFloat,
					Time:       t,
					Value:      item.Value + prevCItem.Value,
					ValueFloat: item.ValueFloat + prevCItem.ValueFloat})
			} else {
				cItems = append(cItems, DataItem{
					SeriesName: newDataSeries.SeriesName,
					IsFloat:    newDataSeries.IsFloat,
					Time:       t,
					Value:      item.Value,
					ValueFloat: item.ValueFloat})
			}
		} else {
			if len(cItems) > 0 {
				prevCItem := cItems[len(cItems)-1]
				cItems = append(cItems, DataItem{
					SeriesName: newDataSeries.SeriesName,
					IsFloat:    newDataSeries.IsFloat,
					Time:       t,
					Value:      prevCItem.Value,
					ValueFloat: prevCItem.ValueFloat})
			} else {
				cItems = append(cItems, DataItem{
					SeriesName: newDataSeries.SeriesName,
					IsFloat:    newDataSeries.IsFloat,
					Time:       t,
					Value:      0,
					ValueFloat: 0})
			}
		}
	}
	for _, cItem := range cItems {
		newDataSeries.AddItem(cItem)
	}
	return newDataSeries, nil
}

func (series *DataSeries) ToQuarter() DataSeries {
	newDataSeries := DataSeries{
		SeriesName: series.SeriesName,
		ItemMap:    map[string]DataItem{},
		IsFloat:    series.IsFloat,
		Interval:   timeutil.Quarter}
	for _, item := range series.ItemMap {
		newDataSeries.AddItem(DataItem{
			SeriesName: item.SeriesName,
			Time:       timeutil.QuarterStart(item.Time),
			IsFloat:    item.IsFloat,
			Value:      item.Value,
			ValueFloat: item.ValueFloat})
	}
	return newDataSeries
}

func (ds *DataSeries) WriteXLSX(filename, sheetname, col1, col2 string) error {
	col1 = strings.TrimSpace(col1)
	col2 = strings.TrimSpace(col2)
	if len(col1) == 0 {
		col1 = "Date"
	}
	if len(col2) == 0 {
		col2 = "Value"
	}
	rows := [][]interface{}{{col1, col2}}
	items := ds.ItemsSorted()
	for _, item := range items {
		if ds.IsFloat {
			rows = append(rows, []interface{}{item.Time, item.ValueFloat})
		} else {
			rows = append(rows, []interface{}{item.Time, item.Value})
		}
	}
	return table.WriteXLSXInterface(filename, table.SheetData{
		SheetName: sheetname,
		Rows:      rows})
}

func AggregateSeries(s1 DataSeries) DataSeries {
	aggregate := NewDataSeries()
	sortedItems := s1.ItemsSorted()
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
