package timeseries

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/grokify/gocharts/data/point"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/simplego/sort/sortutil"
	"github.com/grokify/simplego/time/month"
	"github.com/grokify/simplego/time/timeutil"
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

// AddItem adds data item. It will sum values when
// existing time unit is encountered.
func (ds *DataSeries) AddItem(item DataItem) {
	if ds.ItemMap == nil {
		ds.ItemMap = map[string]DataItem{}
	}
	if len(item.SeriesName) == 0 {
		item.SeriesName = ds.SeriesName
	}
	item.Time = item.Time.UTC()
	rfc := item.Time.Format(time.RFC3339)
	if existingItem, ok := ds.ItemMap[rfc]; ok {
		existingItem.Value += item.Value
		existingItem.ValueFloat += item.ValueFloat
		ds.ItemMap[rfc] = existingItem
	} else {
		ds.ItemMap[rfc] = item
	}
}

func (ds *DataSeries) SetSeriesName(seriesName string) {
	ds.SeriesName = seriesName
	for k, v := range ds.ItemMap {
		v.SeriesName = seriesName
		ds.ItemMap[k] = v
	}
}

// Keys returns a sorted listed of Item keys.
func (ds *DataSeries) Keys() []string {
	keys := []string{}
	for key := range ds.ItemMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// ItemsSorted returns sorted DataItems. This currently uses
// a simple string sort on RFC3339 times.
func (ds *DataSeries) ItemsSorted() []DataItem {
	keys := ds.Keys()
	itemsSorted := []DataItem{}
	for _, key := range keys {
		item, ok := ds.ItemMap[key]
		if !ok {
			panic(fmt.Sprintf("KEY_NOT_FOUND [%s]", key))
		}
		itemsSorted = append(itemsSorted, item)
	}
	return itemsSorted
}

func (ds *DataSeries) Last() (DataItem, error) {
	items := ds.ItemsSorted()
	if len(items) == 0 {
		return DataItem{}, errors.New("E_NO_ITEMS")
	}
	return items[len(items)-1], nil
}

func (ds *DataSeries) Pop() (DataItem, error) {
	items := ds.ItemsSorted()
	if len(items) == 0 {
		return DataItem{}, errors.New("E_NO_ERROR")
	}
	last := items[len(items)-1]
	rfc := last.Time.Format(time.RFC3339)
	delete(ds.ItemMap, rfc)
	return last, nil
}

func (ds *DataSeries) LastItem(skipIfTimePartialValueLessPrev bool) (DataItem, error) {
	items := ds.ItemsSorted()
	if len(items) == 0 {
		return DataItem{}, errors.New("E_NO_ITEMS")
	}
	if len(items) == 1 {
		return items[0], nil
	}
	itemLast := items[len(items)-1]
	if skipIfTimePartialValueLessPrev {
		itemPrev := items[len(items)-2]
		dtNow := time.Now().UTC()
		if ds.Interval == timeutil.Month {
			dtNow = month.MonthBegin(dtNow, 0)
		}
		if itemLast.Time.Equal(dtNow) {
			if itemLast.ValueInt64() > itemPrev.ValueInt64() {
				return itemLast, nil
			} else {
				return itemPrev, nil
			}
		}
	}
	return itemLast, nil
}

func (ds *DataSeries) minMaxValuesInt64Only() (int64, int64) {
	int64s := []int64{}
	for _, item := range ds.ItemMap {
		int64s = append(int64s, item.Value)
	}
	if len(int64s) == 0 {
		return 0, 0
	}
	sort.Sort(sortutil.Int64Slice(int64s))
	return int64s[0], int64s[len(int64s)-1]
}

func (ds *DataSeries) minMaxValuesFloat64Only() (float64, float64) {
	float64s := []float64{}
	for _, item := range ds.ItemMap {
		float64s = append(float64s, item.ValueFloat)
	}
	if len(float64s) == 0 {
		return 0, 0
	}
	float64slice := sort.Float64Slice(float64s)
	sort.Sort(float64slice)

	return float64slice[0], float64slice[len(float64slice)-1]
}

func (ds *DataSeries) MinMaxValues() (int64, int64) {
	if ds.IsFloat {
		min, max := ds.minMaxValuesFloat64Only()
		return int64(min), int64(max)
	}
	return ds.minMaxValuesInt64Only()
}

func (ds *DataSeries) MinMaxValuesFloat64() (float64, float64) {
	if ds.IsFloat {
		return ds.minMaxValuesFloat64Only()
	}
	min, max := ds.minMaxValuesInt64Only()
	return float64(min), float64(max)
}

func (ds *DataSeries) MinValue() int64 {
	min, _ := ds.MinMaxValues()
	return min
}

func (ds *DataSeries) MaxValue() int64 {
	_, max := ds.MinMaxValues()
	return max
}

func (ds *DataSeries) OneItemMaxValue() (DataItem, error) {
	max := DataItem{}
	if len(ds.ItemMap) == 0 {
		return max, errors.New("Empty Set has no Max Value Item")
	}
	first := true
	for _, item := range ds.ItemMap {
		if first {
			max = item
			first = false
		}
		if ds.IsFloat && item.ValueFloat > max.ValueFloat {
			max = item
		} else if item.Value > max.Value {
			max = item
		}
	}
	return max, nil
}

func (ds *DataSeries) TimeSlice(sortSlice bool) timeutil.TimeSlice {
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

func (ds *DataSeries) ToMonth(inflate bool) DataSeries {
	newDataSeries := DataSeries{
		SeriesName: ds.SeriesName,
		ItemMap:    map[string]DataItem{},
		IsFloat:    ds.IsFloat,
		Interval:   timeutil.Month}
	for _, item := range ds.ItemMap {
		newDataSeries.AddItem(DataItem{
			SeriesName: item.SeriesName,
			Time:       month.MonthBegin(item.Time, 0),
			IsFloat:    item.IsFloat,
			Value:      item.Value,
			ValueFloat: item.ValueFloat})
	}
	if inflate {
		timeSeries := timeutil.TimeSeriesSlice(
			timeutil.Month,
			newDataSeries.ItemTimes())
		for _, dt := range timeSeries {
			newDataSeries.AddItem(DataItem{
				SeriesName: newDataSeries.SeriesName,
				Time:       dt,
				IsFloat:    newDataSeries.IsFloat,
				Value:      0,
				ValueFloat: 0.0})
		}
	}
	return newDataSeries
}

func (ds *DataSeries) ToMonthCumulative(inflate bool, timesInput ...time.Time) (DataSeries, error) {
	newDataSeries := DataSeries{
		SeriesName: ds.SeriesName,
		ItemMap:    map[string]DataItem{},
		IsFloat:    ds.IsFloat,
		Interval:   timeutil.Month}
	dsMonth := ds.ToMonth(inflate)
	var min time.Time
	var max time.Time
	var err error
	if len(timesInput) > 0 {
		min, max, err = timeutil.TimeSliceMinMax(timesInput)
		if err != nil {
			return newDataSeries, err
		}
	} else {
		min, max, err = timeutil.TimeSliceMinMax(dsMonth.TimeSlice(false))
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

func (ds *DataSeries) ToQuarter() DataSeries {
	newDataSeries := DataSeries{
		SeriesName: ds.SeriesName,
		ItemMap:    map[string]DataItem{},
		IsFloat:    ds.IsFloat,
		Interval:   timeutil.Quarter}
	for _, item := range ds.ItemMap {
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

func AggregateSeries(series DataSeries) DataSeries {
	aggregate := NewDataSeries()
	sortedItems := series.ItemsSorted()
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

func (ds *DataSeries) TimeSeries(interval timeutil.Interval) []time.Time {
	return timeutil.TimeSeriesSlice(interval, ds.ItemTimes())
}

func (ds *DataSeries) ItemTimes() []time.Time {
	times := []time.Time{}
	for _, item := range ds.ItemMap {
		times = append(times, item.Time)
	}
	return times
}

func (ds *DataSeries) MinMaxTimes() (time.Time, time.Time) {
	return timeutil.SliceMinMax(ds.ItemTimes())
}

func (ds *DataSeries) Stats() point.PointSet {
	ps := point.NewPointSet()
	ps.IsFloat = ds.IsFloat
	for rfc3339, item := range ds.ItemMap {
		point := point.Point{
			Name:    rfc3339,
			IsFloat: item.IsFloat}
		if item.IsFloat {
			point.AbsoluteFloat = item.ValueFloat
		} else {
			point.AbsoluteInt = item.Value
		}
		// Percentage:  float64(itemCount) / float64(totalCount) * 100}
		ps.PointsMap[rfc3339] = point
	}
	ps.Inflate()
	return ps
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
