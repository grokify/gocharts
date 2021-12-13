package interval

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/grokify/gocharts/data/timeseries"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/maputil"
	"github.com/pkg/errors"
)

type SeriesType int

const (
	Source SeriesType = iota
	Output
	OutputAggregate
)

// TimeSeriesSet is used to prepare histogram data for
// static timeseries charts, adding zero values for
// time slots as necessary. Usage is to create to call:
// NewTimeSeriesSet("quarter"), AddItem() and then Inflate()
type TimeSeriesSet struct {
	SourceSeriesMap          map[string]timeseries.TimeSeries
	OutputSeriesMap          map[string]timeseries.TimeSeries
	OutputAggregateSeriesMap map[string]timeseries.TimeSeries
	SeriesIntervals          SeriesIntervals
	AllSeriesName            string
}

func NewTimeSeriesSet(interval timeutil.Interval, weekStart time.Weekday) TimeSeriesSet {
	return TimeSeriesSet{
		SourceSeriesMap:          map[string]timeseries.TimeSeries{},
		OutputSeriesMap:          map[string]timeseries.TimeSeries{},
		OutputAggregateSeriesMap: map[string]timeseries.TimeSeries{},
		SeriesIntervals:          SeriesIntervals{Interval: interval, WeekStart: weekStart}}
}

func (set *TimeSeriesSet) SeriesNamesSorted() []string {
	return maputil.StringKeysSorted(set.OutputSeriesMap)
}

func (set *TimeSeriesSet) AddItem(item timeseries.TimeItem) {
	item.SeriesName = strings.TrimSpace(item.SeriesName)
	if _, ok := set.SourceSeriesMap[item.SeriesName]; !ok {
		set.SourceSeriesMap[item.SeriesName] =
			timeseries.TimeSeries{
				SeriesName: item.SeriesName,
				ItemMap:    map[string]timeseries.TimeItem{}}
	}
	series := set.SourceSeriesMap[item.SeriesName]
	series.AddItems(item)
	set.SourceSeriesMap[item.SeriesName] = series
}

func (set *TimeSeriesSet) Inflate() error {
	if err := set.inflateSource(); err != nil {
		return err
	}
	if err := set.inflateOutput(); err != nil {
		return err
	}
	set.addAllSeries(set.AllSeriesName)
	return nil
}

func (set *TimeSeriesSet) inflateSource() error {
	for _, series := range set.SourceSeriesMap {
		set.SeriesIntervals.ProcItemsMap(series.ItemMap)
	}
	return set.SeriesIntervals.Inflate()
}

func (set *TimeSeriesSet) inflateOutput() error {
	for seriesName, series := range set.SourceSeriesMap {
		output, err := set.BuildOutputSeries(series)
		if err != nil {
			return err
		}
		set.OutputSeriesMap[seriesName] = output
		set.OutputAggregateSeriesMap[seriesName] = timeseries.AggregateSeries(output)
	}
	return nil
}

func (set *TimeSeriesSet) addAllSeries(allSeriesName string) {
	if len(strings.TrimSpace(allSeriesName)) == 0 {
		allSeriesName = "All"
	}
	allSeries := timeseries.NewTimeSeries(allSeriesName)

	for _, series := range set.SourceSeriesMap {
		for _, item := range series.ItemMap {
			item.SeriesName = allSeriesName
			allSeries.AddItems(item)
		}
	}

	set.OutputSeriesMap[allSeriesName] = allSeries
	set.OutputAggregateSeriesMap[allSeriesName] = timeseries.AggregateSeries(allSeries)
}

func (set *TimeSeriesSet) GetTimeSeries(seriesName string, seriesType SeriesType) (timeseries.TimeSeries, error) {
	var seriesMap map[string]timeseries.TimeSeries
	switch seriesType {
	case Source:
		seriesMap = set.SourceSeriesMap
	case Output:
		seriesMap = set.OutputSeriesMap
	case OutputAggregate:
		seriesMap = set.OutputAggregateSeriesMap
	default:
		return timeseries.TimeSeries{}, fmt.Errorf("could not find seriesName [%v] seriesType [%v]",
			seriesName,
			seriesType)
	}
	seriesData, ok := seriesMap[seriesName]
	if !ok {
		return timeseries.TimeSeries{}, fmt.Errorf("could not find seriesName [%v] seriesType [%v]",
			seriesName,
			seriesType)
	}
	return seriesData, nil
}

func (set *TimeSeriesSet) BuildOutputSeries(source timeseries.TimeSeries) (timeseries.TimeSeries, error) {
	output := timeseries.NewTimeSeries(set.AllSeriesName)
	for _, item := range source.ItemMap {
		output.SeriesName = item.SeriesName
		ivalStart, err := timeutil.IntervalStart(
			item.Time,
			set.SeriesIntervals.Interval,
			set.SeriesIntervals.WeekStart)
		if err != nil {
			return output, err
		}
		output.AddItems(timeseries.TimeItem{
			SeriesName: item.SeriesName,
			Time:       ivalStart,
			Value:      item.Value})
	}
	for _, dt := range set.SeriesIntervals.CanonicalSeries {
		output.AddItems(timeseries.TimeItem{
			SeriesName: output.SeriesName,
			Time:       dt,
			Value:      0})
	}
	return output, nil
}

func (set *TimeSeriesSet) FlattenData() map[string][]time.Time {
	out := map[string][]time.Time{}
	for seriesName, timeSeries := range set.SourceSeriesMap {
		if _, ok := out[seriesName]; !ok {
			out[seriesName] = []time.Time{}
		}
		times := out[seriesName]
		for _, timeItem := range timeSeries.ItemMap {
			for i := 0; i < int(timeItem.Value); i++ {
				times = append(times, timeItem.Time)
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

func (ival *SeriesIntervals) ProcItemsMap(itemMap map[string]timeseries.TimeItem) {
	for _, timeItem := range itemMap {
		dt := timeItem.Time
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
