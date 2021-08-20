package timeseries

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/grokify/simplego/sort/sortutil"
	"github.com/grokify/simplego/time/month"
	"github.com/grokify/simplego/time/timeslice"
	"github.com/grokify/simplego/time/timeutil"
	"github.com/grokify/simplego/type/stringsutil"
	"github.com/pkg/errors"
)

type TimeSeriesSet struct {
	Name     string
	Series   map[string]TimeSeries
	Times    []time.Time
	Order    []string
	IsFloat  bool
	Interval timeutil.Interval
}

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

// AddFloat64 adds a time value, converting it to a int64 on
// the series type.
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

func (set *TimeSeriesSet) TimeSlice(sortAsc bool) timeslice.TimeSlice {
	times := []time.Time{}
	for _, ds := range set.Series {
		for _, item := range ds.ItemMap {
			times = append(times, item.Time)
		}
	}
	times = timeutil.Sort(timeutil.Distinct(times))
	return month.TimeSeriesMonth(sortAsc, times...)
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
	values := timeslice.TimeSlice{}
	for _, ds := range set.Series {
		min, max := ds.MinMaxTimes()
		values = append(values, min, max)
	}
	sort.Sort(values)
	return values[0], values[len(values)-1]
}

func (set *TimeSeriesSet) MinMaxValues() (int64, int64) {
	values := sortutil.Int64Slice{}
	for _, ds := range set.Series {
		min, max := ds.MinMaxValues()
		values = append(values, min, max)
	}
	sort.Sort(values)
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
