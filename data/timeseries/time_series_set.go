package timeseries

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/grokify/mogo/sort/sortutil"
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/time/year"
	"github.com/grokify/mogo/type/stringsutil"
)

// TimeSeriesSet is a data structure to manage a set of similar `TimeSeries`.
// It is necessary for all `TimeSeries` to have the same value of `IsFloat`.
type TimeSeriesSet struct {
	Name              string
	Series            map[string]TimeSeries
	Times             []time.Time
	Order             []string
	ActualTargetPairs []ActualTargetPair
	IsFloat           bool
	Interval          timeutil.Interval
}

// ActualTargetPair provides metadata on associating two series names that
// represent actual and target data. This can be used to product additional
// data in charts and tables.
type ActualTargetPair struct {
	ActualSeriesName string
	TargetSeriesName string
}

// NewTimeSeriesSet returns an initialized `TimeSeriesSet`.
func NewTimeSeriesSet(name string) TimeSeriesSet {
	return TimeSeriesSet{
		Name:   name,
		Series: map[string]TimeSeries{},
		Times:  []time.Time{},
		Order:  []string{}}
}

// ReadFileTimeSeriesSet reads a time series set file in JSON.
func ReadFileTimeSeriesSet(filename string) (TimeSeriesSet, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return TimeSeriesSet{}, err
	}
	var tset TimeSeriesSet
	return tset, json.Unmarshal(data, &tset)
}

// AddInt64 adds an `int64` value, converting it to a `float64` if necssary based on
// set definition.
func (set *TimeSeriesSet) AddInt64(seriesName string, dt time.Time, value int64) {
	item := TimeItem{
		SeriesSetName: set.Name,
		SeriesName:    seriesName,
		Time:          dt,
		IsFloat:       set.IsFloat}
	if set.IsFloat {
		item.ValueFloat = float64(value)
	} else {
		item.Value = value
	}
	set.AddItems(item)
}

// AddFloat64 adds an `int64` value, converting it to a `int64` if necssary based on
// set definition.
func (set *TimeSeriesSet) AddFloat64(seriesName string, dt time.Time, value float64) {
	item := TimeItem{
		SeriesSetName: set.Name,
		SeriesName:    seriesName,
		Time:          dt,
		IsFloat:       set.IsFloat}
	if set.IsFloat {
		item.ValueFloat = value
	} else {
		item.Value = int64(value)
	}
	set.AddItems(item)
}

func (set *TimeSeriesSet) AddItems(items ...TimeItem) {
	for _, item := range items {
		if _, ok := set.Series[item.SeriesName]; !ok {
			set.Series[item.SeriesName] =
				TimeSeries{
					SeriesSetName: set.Name,
					SeriesName:    item.SeriesName,
					ItemMap:       map[string]TimeItem{},
					IsFloat:       item.IsFloat,
					Interval:      set.Interval}
		}
		ts := set.Series[item.SeriesName]
		ts.AddItems(item)
		set.Series[item.SeriesName] = ts
		set.Times = append(set.Times, item.Time)
	}
}

func (set *TimeSeriesSet) AddSeries(timeSeries ...TimeSeries) error {
	for _, ts := range timeSeries {
		ts.SeriesName = strings.TrimSpace(ts.SeriesName)
		if len(ts.SeriesName) == 0 {
			return errors.New("E_TImeSeriesSet.AddTimeSeries_NO_DataSeries.SeriesName")
		}
		for _, item := range ts.ItemMap {
			if len(item.SeriesName) == 0 || item.SeriesName != ts.SeriesName {
				item.SeriesName = ts.SeriesName
			}
			set.AddItems(item)
		}
	}
	return nil
}

func (set *TimeSeriesSet) GetInt64WithDefault(seriesName, rfc3339fulldate string, def int64) int64 {
	if series, ok := set.Series[seriesName]; !ok {
		return def
	} else if ti, ok := series.ItemMap[rfc3339fulldate]; !ok {
		return def
	} else {
		return ti.Int64()
	}
}

func (set *TimeSeriesSet) Inflate() {
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

func (set *TimeSeriesSet) SeriesNames() []string {
	seriesNames := []string{}
	for seriesName := range set.Series {
		seriesNames = append(seriesNames, seriesName)
	}
	sort.Strings(seriesNames)
	return seriesNames
}

func (set *TimeSeriesSet) GetSeriesByIndex(index int) (TimeSeries, error) {
	if len(set.Order) == 0 && len(set.Series) > 0 {
		set.Inflate()
	}
	if index < len(set.Order) {
		name := set.Order[index]
		if ds, ok := set.Series[name]; ok {
			return ds, nil
		}
	}
	return TimeSeries{}, fmt.Errorf("E_CANNOT_FIND_INDEX_[%d]_SET_COUNT_[%d]", index, len(set.Order))
}

func (set *TimeSeriesSet) Item(seriesName, rfc3339 string) (TimeItem, error) {
	di := TimeItem{}
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

func (set *TimeSeriesSet) SetInterval(interval timeutil.Interval, recursive bool) {
	set.Interval = interval
	if recursive {
		for k, ts := range set.Series {
			ts.Interval = interval
			set.Series[k] = ts
		}
	}
}

func (set *TimeSeriesSet) TimeSlice(sortAsc bool) timeutil.Times {
	var times []time.Time
	for _, ts := range set.Series {
		for _, item := range ts.ItemMap {
			tm := timeutil.NewTimeMore(item.Time, 0)
			if set.Interval == timeutil.IntervalYear && !tm.IsYearStart() {
				panic("timeitem for TimeSeriesSet year is not a year start")
			} else if set.Interval == timeutil.IntervalMonth && !tm.IsMonthStart() {
				panic("timeitem for TimeSeriesSet month is not a month start")
			}
			times = append(times, item.Time)
		}
	}
	times = timeutil.Sort(timeutil.Distinct(times))
	switch set.Interval {
	case timeutil.IntervalMonth:
		return month.TimesMonthStarts(times...)
	case timeutil.IntervalYear:
		return year.TimesYearStarts(times...)
	default:
		return times
	}
}

func (set *TimeSeriesSet) TimeStrings() []string {
	times := []string{}
	for _, ds := range set.Series {
		for rfc3339 := range ds.ItemMap {
			times = append(times, rfc3339)
		}
	}
	return stringsutil.SliceCondenseSpace(times, true, true)
}

func (set *TimeSeriesSet) MinMaxTimes() (time.Time, time.Time) {
	values := timeutil.Times{}
	for _, ds := range set.Series {
		min, max := ds.MinMaxTimes()
		values = append(values, min, max)
	}
	sort.Sort(values)
	return values[0], values[len(values)-1]
}

func (set *TimeSeriesSet) MinMaxValues() (int64, int64) {
	values := []int64{} // sortutil.Int64Slice{}
	for _, ds := range set.Series {
		min, max := ds.MinMaxValues()
		values = append(values, min, max)
	}
	sortutil.Slice(values)
	return values[0], values[len(values)-1]
}

func (set *TimeSeriesSet) MinMaxValuesFloat64() (float64, float64) {
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
