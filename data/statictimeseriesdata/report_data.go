// statictimeseriesdata provides tools for adding and formatting
// static time series data for reporting purposes.
package statictimeseriesdata

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/grokify/gotilla/time/timeutil"
	"github.com/grokify/gotilla/type/maputil"
)

// DataSeriesSet is used to prepare histogram data for
// static timeseries charts, adding zero values for
// time slots as necessary. Usage is to create to call:
// NewDataSeriesSet("quarter"), AddItem() and then Inflate()
type DataSeriesSet struct {
	SourceSeriesMap map[string]DataSeries
	OutputSeriesMap map[string]DataSeries
	SeriesIntervals SeriesIntervals
}

func NewDataSeriesSet(interval timeutil.Interval, weekStart time.Weekday) DataSeriesSet {
	set := DataSeriesSet{
		SourceSeriesMap: map[string]DataSeries{},
		OutputSeriesMap: map[string]DataSeries{},
		SeriesIntervals: SeriesIntervals{Interval: interval, WeekStart: weekStart}}
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
	return maputil.StringKeysSorted(set.SourceSeriesMap)
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

func (set *DataSeriesSet) Inflate() error {
	err := set.inflateSource()
	if err != nil {
		return err
	}
	return set.inflateOutput()
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
	}
	return nil
}

func (set *DataSeriesSet) BuildOutputSeries(source DataSeries) (DataSeries, error) {
	output := NewDataSeries()
	for _, item := range source.ItemMap {
		output.SeriesName = item.SeriesName
		ivalStart, err := timeutil.IntervalStart(
			item.Time, set.SeriesIntervals.Interval, set.SeriesIntervals.WeekStart)
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

// SortedItems returns sorted DataItems. This currently uses
// a simple tring sort on RFC3339 times. For dates that are not
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
