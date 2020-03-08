// statictimeseriesdata provides tools for adding and formatting
// static time series data for reporting purposes.
package statictimeseries

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/grokify/gotilla/sort/sortutil"
	"github.com/grokify/gotilla/time/timeutil"
	"github.com/grokify/gotilla/type/maputil"
	"github.com/pkg/errors"
)

type SeriesType int

const (
	Source SeriesType = iota
	Output
	OutputAggregate
)

// DataSeriesSet is used to prepare histogram data for
// static timeseries charts, adding zero values for
// time slots as necessary. Usage is to create to call:
// NewDataSeriesSet("quarter"), AddItem() and then Inflate()
type DataSeriesSet struct {
	SourceSeriesMap          map[string]DataSeries
	OutputSeriesMap          map[string]DataSeries
	OutputAggregateSeriesMap map[string]DataSeries
	SeriesIntervals          SeriesIntervals
	AllSeriesName            string
}

func NewDataSeriesSet(interval timeutil.Interval, weekStart time.Weekday) DataSeriesSet {
	set := DataSeriesSet{
		SourceSeriesMap:          map[string]DataSeries{},
		OutputSeriesMap:          map[string]DataSeries{},
		OutputAggregateSeriesMap: map[string]DataSeries{},
		SeriesIntervals:          SeriesIntervals{Interval: interval, WeekStart: weekStart}}
	return set
}

type DataSeries struct {
	SeriesName string
	ItemMap    map[string]DataItem
}

func NewDataSeries() DataSeries {
	return DataSeries{ItemMap: map[string]DataItem{}}
}

type DataItem struct {
	SeriesName string
	Time       time.Time
	Value      int64
}

func (set *DataSeriesSet) SeriesNamesSorted() []string {
	return maputil.StringKeysSorted(set.OutputSeriesMap)
}

func (set *DataSeriesSet) AddItem(item DataItem) {
	item.SeriesName = strings.TrimSpace(item.SeriesName)
	if _, ok := set.SourceSeriesMap[item.SeriesName]; !ok {
		set.SourceSeriesMap[item.SeriesName] =
			DataSeries{
				SeriesName: item.SeriesName,
				ItemMap:    map[string]DataItem{}}
	}
	series := set.SourceSeriesMap[item.SeriesName]
	series.AddItem(item)
	set.SourceSeriesMap[item.SeriesName] = series
}

// AddItem adds data item. It will sum values when
// existing time unit is encountered.
func (series *DataSeries) AddItem(item DataItem) {
	if series.ItemMap == nil {
		series.ItemMap = map[string]DataItem{}
	}
	item.Time = item.Time.UTC()
	rfc := item.Time.Format(time.RFC3339)
	if _, ok := series.ItemMap[rfc]; !ok {
		series.ItemMap[rfc] = item
	} else {
		existingItem := series.ItemMap[rfc]
		existingItem.Value += item.Value
		series.ItemMap[rfc] = existingItem
	}
}

func (series *DataSeries) ItemsSorted() []DataItem {
	keys := []string{}
	for key := range series.ItemMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	items := []DataItem{}
	for _, key := range keys {
		item, ok := series.ItemMap[key]
		if !ok {
			panic("KEY_NOT_FOUND")
		}
		items = append(items, item)
	}
	return items
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

func (series *DataSeries) MinMaxValues() (int64, int64) {
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

func (series *DataSeries) MinValue() int64 {
	min, _ := series.MinMaxValues()
	return min
}

func (series *DataSeries) MaxValue() int64 {
	_, max := series.MinMaxValues()
	return max
}

func (set *DataSeriesSet) Inflate() error {
	if err := set.inflateSource(); err != nil {
		return err
	}
	if err := set.inflateOutput(); err != nil {
		return err
	}
	set.addAllSeries(set.AllSeriesName)
	return nil
}

func (set *DataSeriesSet) inflateSource() error {
	for _, series := range set.SourceSeriesMap {
		set.SeriesIntervals.ProcItemsMap(series.ItemMap)
	}
	return set.SeriesIntervals.Inflate()
}

func (set *DataSeriesSet) inflateOutput() error {
	for seriesName, series := range set.SourceSeriesMap {
		output, err := set.BuildOutputSeries(series)
		if err != nil {
			return err
		}
		set.OutputSeriesMap[seriesName] = output
		set.OutputAggregateSeriesMap[seriesName] = AggregateSeries(output)
	}
	return nil
}

func (set *DataSeriesSet) addAllSeries(allSeriesName string) {
	if len(strings.TrimSpace(allSeriesName)) == 0 {
		allSeriesName = "All"
	}
	allSeries := NewDataSeries()
	allSeries.SeriesName = allSeriesName

	for _, series := range set.SourceSeriesMap {
		for _, item := range series.ItemMap {
			item.SeriesName = allSeriesName
			allSeries.AddItem(item)
		}
	}

	set.OutputSeriesMap[allSeriesName] = allSeries
	set.OutputAggregateSeriesMap[allSeriesName] = AggregateSeries(allSeries)
}

func (set *DataSeriesSet) GetDataSeries(seriesName string, seriesType SeriesType) (DataSeries, error) {
	seriesMap := map[string]DataSeries{}
	switch seriesType {
	case Source:
		seriesMap = set.SourceSeriesMap
	case Output:
		seriesMap = set.OutputSeriesMap
	case OutputAggregate:
		seriesMap = set.OutputAggregateSeriesMap
	default:
		return DataSeries{}, fmt.Errorf("Could not find seriesName [%v] seriesType [%v]",
			seriesName,
			seriesType)
	}
	seriesData, ok := seriesMap[seriesName]
	if !ok {
		return DataSeries{}, fmt.Errorf("Could not find seriesName [%v] seriesType [%v]",
			seriesName,
			seriesType)
	}
	return seriesData, nil
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

func (set *DataSeriesSet) BuildOutputSeries(source DataSeries) (DataSeries, error) {
	output := NewDataSeries()
	for _, item := range source.ItemMap {
		output.SeriesName = item.SeriesName
		ivalStart, err := timeutil.IntervalStart(
			item.Time,
			set.SeriesIntervals.Interval,
			set.SeriesIntervals.WeekStart)
		if err != nil {
			return output, err
		}
		output.AddItem(DataItem{
			SeriesName: item.SeriesName,
			Time:       ivalStart,
			Value:      item.Value})
	}
	for _, dt := range set.SeriesIntervals.CanonicalSeries {
		output.AddItem(DataItem{
			SeriesName: output.SeriesName,
			Time:       dt,
			Value:      0})
	}
	return output, nil
}

func (set *DataSeriesSet) FlattenData() map[string][]time.Time {
	out := map[string][]time.Time{}
	for seriesName, dataSeries := range set.SourceSeriesMap {
		if _, ok := out[seriesName]; !ok {
			out[seriesName] = []time.Time{}
		}
		times := out[seriesName]
		for _, dataItem := range dataSeries.ItemMap {
			for i := 0; i < int(dataItem.Value); i++ {
				times = append(times, dataItem.Time)
			}
		}
		out[seriesName] = times
	}
	for seriesName, timeSlice := range out {
		sort.Slice(timeSlice, func(i, j int) bool {
			return timeSlice[i].Before(timeSlice[j])
		})
		out[seriesName] = timeSlice
	}
	return out
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

type SeriesIntervals struct {
	Interval        timeutil.Interval
	WeekStart       time.Weekday
	Max             time.Time
	Min             time.Time
	CanonicalSeries []time.Time
}

func (ival *SeriesIntervals) areEndpointsSet() bool {
	if ival.Max.IsZero() || ival.Min.IsZero() {
		return false
	}
	return true
}

func (ival *SeriesIntervals) ProcItemsMap(itemMap map[string]DataItem) {
	for _, dataItem := range itemMap {
		dt := dataItem.Time
		if !ival.areEndpointsSet() {
			ival.Max = dt
			ival.Min = dt
			continue
		}
		if timeutil.IsGreaterThan(dt, ival.Max, false) {
			ival.Max = dt
		}
		if timeutil.IsLessThan(dt, ival.Min, false) {
			ival.Min = dt
		}
	}
}

func (ival *SeriesIntervals) Inflate() error {
	err := ival.buildMinMaxEndpoints()
	if err != nil {
		return err
	}
	ival.buildCanonicalSeries()
	return nil
}

func (ival *SeriesIntervals) buildMinMaxEndpoints() error {
	if !ival.areEndpointsSet() {
		return errors.New("Cannot build canonical dates without initialized dates.")
	}
	ival.Max = ival.Max.UTC()
	ival.Min = ival.Min.UTC()
	switch ival.Interval.String() {
	case "year":
		ival.Max = timeutil.YearStart(ival.Max)
		ival.Min = timeutil.YearStart(ival.Min)
	case "quarter":
		ival.Max = timeutil.QuarterStart(ival.Max)
		ival.Min = timeutil.QuarterStart(ival.Min)
	case "month":
		ival.Max = timeutil.MonthStart(ival.Max)
		ival.Min = timeutil.MonthStart(ival.Min)
	case "week":
		max, err := timeutil.WeekStart(ival.Max, ival.WeekStart)
		if err != nil {
			return err
		}
		ival.Max = max
		min, err := timeutil.WeekStart(ival.Min, ival.WeekStart)
		if err != nil {
			return err
		}
		ival.Min = min
	default:
		panic(fmt.Sprintf("Interval [%v] not supported.", ival.Interval))
	}
	return nil
}

func (ival *SeriesIntervals) buildCanonicalSeries() {
	canonicalSeries := []time.Time{}
	curTime := ival.Min
	for timeutil.IsLessThan(curTime, ival.Max, true) {
		canonicalSeries = append(canonicalSeries, curTime)
		switch ival.Interval.String() {
		case "year":
			curTime = timeutil.TimeDt4AddNYears(curTime, 1)
		case "quarter":
			curTime = timeutil.TimeDt6AddNMonths(curTime, 3)
		case "month":
			curTime = timeutil.TimeDt6AddNMonths(curTime, 1)
		case "week":
			dur, _ := time.ParseDuration(fmt.Sprintf("%vh", 24*7))
			curTime = curTime.Add(dur)
		default:
			panic(fmt.Sprintf("Interval [%v] not supported.", ival.Interval))
		}
	}
	ival.CanonicalSeries = canonicalSeries
}
