package interval

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/grokify/gocharts/data/statictimeseries"
	"github.com/grokify/simplego/time/timeutil"
	"github.com/grokify/simplego/type/maputil"
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
	SourceSeriesMap          map[string]statictimeseries.DataSeries
	OutputSeriesMap          map[string]statictimeseries.DataSeries
	OutputAggregateSeriesMap map[string]statictimeseries.DataSeries
	SeriesIntervals          SeriesIntervals
	AllSeriesName            string
}

func NewDataSeriesSet(interval timeutil.Interval, weekStart time.Weekday) DataSeriesSet {
	return DataSeriesSet{
		SourceSeriesMap:          map[string]statictimeseries.DataSeries{},
		OutputSeriesMap:          map[string]statictimeseries.DataSeries{},
		OutputAggregateSeriesMap: map[string]statictimeseries.DataSeries{},
		SeriesIntervals:          SeriesIntervals{Interval: interval, WeekStart: weekStart}}
}

func (set *DataSeriesSet) SeriesNamesSorted() []string {
	return maputil.StringKeysSorted(set.OutputSeriesMap)
}

func (set *DataSeriesSet) AddItem(item statictimeseries.DataItem) {
	item.SeriesName = strings.TrimSpace(item.SeriesName)
	if _, ok := set.SourceSeriesMap[item.SeriesName]; !ok {
		set.SourceSeriesMap[item.SeriesName] =
			statictimeseries.DataSeries{
				SeriesName: item.SeriesName,
				ItemMap:    map[string]statictimeseries.DataItem{}}
	}
	series := set.SourceSeriesMap[item.SeriesName]
	series.AddItem(item)
	set.SourceSeriesMap[item.SeriesName] = series
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
		set.OutputAggregateSeriesMap[seriesName] = statictimeseries.AggregateSeries(output)
	}
	return nil
}

func (set *DataSeriesSet) addAllSeries(allSeriesName string) {
	if len(strings.TrimSpace(allSeriesName)) == 0 {
		allSeriesName = "All"
	}
	allSeries := statictimeseries.NewDataSeries()
	allSeries.SeriesName = allSeriesName

	for _, series := range set.SourceSeriesMap {
		for _, item := range series.ItemMap {
			item.SeriesName = allSeriesName
			allSeries.AddItem(item)
		}
	}

	set.OutputSeriesMap[allSeriesName] = allSeries
	set.OutputAggregateSeriesMap[allSeriesName] = statictimeseries.AggregateSeries(allSeries)
}

func (set *DataSeriesSet) GetDataSeries(seriesName string, seriesType SeriesType) (statictimeseries.DataSeries, error) {
	seriesMap := map[string]statictimeseries.DataSeries{}
	switch seriesType {
	case Source:
		seriesMap = set.SourceSeriesMap
	case Output:
		seriesMap = set.OutputSeriesMap
	case OutputAggregate:
		seriesMap = set.OutputAggregateSeriesMap
	default:
		return statictimeseries.DataSeries{}, fmt.Errorf("Could not find seriesName [%v] seriesType [%v]",
			seriesName,
			seriesType)
	}
	seriesData, ok := seriesMap[seriesName]
	if !ok {
		return statictimeseries.DataSeries{}, fmt.Errorf("Could not find seriesName [%v] seriesType [%v]",
			seriesName,
			seriesType)
	}
	return seriesData, nil
}

func (set *DataSeriesSet) BuildOutputSeries(source statictimeseries.DataSeries) (statictimeseries.DataSeries, error) {
	output := statictimeseries.NewDataSeries()
	for _, item := range source.ItemMap {
		output.SeriesName = item.SeriesName
		ivalStart, err := timeutil.IntervalStart(
			item.Time,
			set.SeriesIntervals.Interval,
			set.SeriesIntervals.WeekStart)
		if err != nil {
			return output, err
		}
		output.AddItem(statictimeseries.DataItem{
			SeriesName: item.SeriesName,
			Time:       ivalStart,
			Value:      item.Value})
	}
	for _, dt := range set.SeriesIntervals.CanonicalSeries {
		output.AddItem(statictimeseries.DataItem{
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

func (ival *SeriesIntervals) ProcItemsMap(itemMap map[string]statictimeseries.DataItem) {
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
