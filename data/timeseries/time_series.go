package timeseries

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/time/timeutil"

	"github.com/grokify/gocharts/v2/data/point"
	"github.com/grokify/gocharts/v2/data/table"
)

type TimeSeries struct {
	SeriesName    string
	SeriesSetName string
	ItemMap       map[string]TimeItem
	IsFloat       bool
	Interval      timeutil.Interval
}

// NewTimeSeries instantiates a `TimeSeries` struct.
func NewTimeSeries(name string) TimeSeries {
	return TimeSeries{
		SeriesName: name,
		ItemMap:    map[string]TimeItem{}}
}

// ReadFileTimeSeries reads a time series file in JSON.
func ReadFileTimeSeries(filename string) (TimeSeries, error) {
	if data, err := os.ReadFile(filename); err != nil {
		return TimeSeries{}, err
	} else {
		var ts TimeSeries
		return ts, json.Unmarshal(data, &ts)
	}
}

// AddInt64 adds a time value, converting it to a float on the series type.
func (ts *TimeSeries) AddInt64(t time.Time, value int64) {
	item := TimeItem{
		SeriesName:    ts.SeriesName,
		SeriesSetName: ts.SeriesSetName,
		Time:          t,
		IsFloat:       ts.IsFloat}
	if ts.IsFloat {
		item.ValueFloat = float64(value)
	} else {
		item.Value = value
	}
	ts.AddItems(item)
}

// AddFloat64 adds a time value, converting it to a int64 on the series type.
func (ts *TimeSeries) AddFloat64(t time.Time, value float64) {
	item := TimeItem{
		SeriesName:    ts.SeriesName,
		SeriesSetName: ts.SeriesSetName,
		Time:          t,
		IsFloat:       ts.IsFloat}
	if ts.IsFloat {
		item.ValueFloat = value
	} else {
		item.Value = int64(value)
	}
	ts.AddItems(item)
}

// AddItems adds a `TimeItem`. It will sum values when existing time unit is encountered.
func (ts *TimeSeries) AddItems(items ...TimeItem) {
	for _, item := range items {
		if ts.ItemMap == nil {
			ts.ItemMap = map[string]TimeItem{}
		}
		if len(item.SeriesName) == 0 {
			item.SeriesName = ts.SeriesName
		}
		item.Time = item.Time.UTC()
		rfc := item.Time.Format(time.RFC3339)
		if existingItem, ok := ts.ItemMap[rfc]; ok {
			existingItem.Value += item.Value
			existingItem.ValueFloat += item.ValueFloat
			ts.ItemMap[rfc] = existingItem
		} else {
			ts.ItemMap[rfc] = item
		}
	}
}

/*
type AddTableOpts struct {
	TimeColIdx   uint
	TimeFormat   string
	CountColIdx  uint
	CountIsFloat bool
}

func DefaultAddTableOpts() *AddTableOpts {
	return &AddTableOpts{
		TimeColIdx:   0,
		TimeFormat:   time.RFC3339,
		CountColIdx:  1,
		CountIsFloat: false}
}
*/

func (ts *TimeSeries) AddTable(tbl *table.Table, timeColIdx uint, timeFormat string, countColIdx uint, countIsFloat bool) error {
	// func (ts *TimeSeries) AddTable(tbl *table.Table, opts *AddTableOpts) error {
	if tbl == nil {
		return errors.New("table cannot be nil")
	} else if timeColIdx == countColIdx {
		return errors.New("column indexes for time and count cannot be the same")
	}

	// if opts == nil {
	//	opts = DefaultAddTableOpts()
	// }
	// if opts.CountIsFloat {
	//	ts.IsFloat = true
	// }
	// dtFormat := opts.TimeFormat

	if len(timeFormat) == 0 {
		timeFormat = time.RFC3339
	}
	for _, row := range tbl.Rows {
		if len(row) == 0 {
			continue
		} else if int(timeColIdx) >= len(row) {
			return fmt.Errorf("time column doesn't exist")
		} else if int(countColIdx) >= len(row) {
			return fmt.Errorf("count column doesn't exist")
		}
		dt, err := time.Parse(timeFormat, row[int(timeColIdx)])
		if err != nil {
			return err
		}
		countStr := strings.TrimSpace(row[int(countColIdx)])
		if countStr == "" {
			ts.AddInt64(dt, 0)
		} else {
			if countIsFloat {
				if countFloat, err := strconv.ParseFloat(countStr, 64); err != nil {
					return err
				} else {
					ts.AddFloat64(dt, countFloat)
				}
			} else {
				if countInt, err := strconv.Atoi(countStr); err != nil {
					return err
				} else {
					ts.AddInt64(dt, int64(countInt))
				}
			}
		}
	}
	return nil
}

func (ts *TimeSeries) ConvertFloat64() {
	for rfc, ti := range ts.ItemMap {
		if ti.IsFloat {
			continue
		}
		ti.ValueFloat = float64(ti.Value)
		ti.IsFloat = true
		ts.ItemMap[rfc] = ti
	}
	ts.IsFloat = true
}

func (ts *TimeSeries) ConvertInt64() {
	for rfc, ti := range ts.ItemMap {
		if !ti.IsFloat {
			continue
		}
		ti.Value = int64(ti.ValueFloat)
		ti.IsFloat = false
		ts.ItemMap[rfc] = ti
	}
	ts.IsFloat = false
}

// Clone returns a copy of the `TimeSeries` struct.
func (ts *TimeSeries) Clone() TimeSeries {
	clone := TimeSeries{
		SeriesName:    ts.SeriesName,
		SeriesSetName: ts.SeriesSetName,
		ItemMap:       map[string]TimeItem{},
		IsFloat:       ts.IsFloat,
		Interval:      ts.Interval}
	for k, v := range ts.ItemMap {
		clone.ItemMap[k] = v
	}
	return clone
}

func (ts *TimeSeries) SetSeriesName(seriesName string) {
	ts.SeriesName = seriesName
	for k, v := range ts.ItemMap {
		v.SeriesName = seriesName
		ts.ItemMap[k] = v
	}
}

// Get returns a `TimeItem` given a `time.Time`.
func (ts *TimeSeries) Get(t time.Time) (TimeItem, error) {
	for _, ti := range ts.ItemMap {
		if ti.Time.Equal(t) {
			return ti, nil
		}
	}
	return TimeItem{}, fmt.Errorf("time not found [%s]", t.Format(time.RFC3339))
}

func (ts *TimeSeries) Last() (TimeItem, error) {
	items := ts.ItemsSorted()
	if len(items) == 0 {
		return TimeItem{}, ErrNoTimeItem
	}
	return items[len(items)-1], nil
}

// Pop removes the item with the chronologically last
// time. This is useful when generating interval
// charts and the last period has not concluded, thus
// providing an inaccurate projection when compared to
// previous full months of data.
func (ts *TimeSeries) Pop() (TimeItem, error) {
	items := ts.ItemsSorted()
	if len(items) == 0 {
		return TimeItem{}, ErrNoTimeItem
	}
	last := items[len(items)-1]
	rfc := last.Time.Format(time.RFC3339)
	delete(ts.ItemMap, rfc)
	return last, nil
}

func (ts *TimeSeries) Times(sortSlice bool) timeutil.Times {
	times := timeutil.Times{}
	for _, item := range ts.ItemMap {
		times = append(times, item.Time)
	}
	if sortSlice {
		sort.Sort(times)
	}
	return times
}

func (ts *TimeSeries) DeleteTime(dt time.Time) {
	delete(ts.ItemMap, dt.Format(time.RFC3339))
}

// ToMonth aggregates time values into months. `addZeroValueMonths` is used to add months with `0` values.
func (ts *TimeSeries) ToMonth(addZeroValueMonths bool, monthsFilter ...time.Month) TimeSeries {
	newTimeSeries := NewTimeSeries(ts.SeriesName)
	newTimeSeries.Interval = timeutil.IntervalYear
	newTimeSeries.IsFloat = ts.IsFloat
	monthsFilterMap := map[time.Month]int{}
	for _, m := range monthsFilter {
		monthsFilterMap[m] = 1
	}
	for _, item := range ts.ItemMap {
		if len(monthsFilterMap) > 0 {
			if _, ok := monthsFilterMap[item.Time.Month()]; !ok {
				continue
			}
		}
		newTimeSeries.AddItems(TimeItem{
			SeriesName: item.SeriesName,
			Time:       month.MonthStart(item.Time, 0),
			IsFloat:    item.IsFloat,
			Value:      item.Value,
			ValueFloat: item.ValueFloat})
	}
	if addZeroValueMonths {
		timeSeries := timeutil.TimeSeriesSlice(timeutil.IntervalMonth, newTimeSeries.ItemTimes())
		for _, dt := range timeSeries {
			newTimeSeries.AddItems(TimeItem{
				SeriesName: newTimeSeries.SeriesName,
				Time:       dt,
				IsFloat:    newTimeSeries.IsFloat,
				Value:      0,
				ValueFloat: 0.0})
		}
	}
	return newTimeSeries
}

func (ts *TimeSeries) ToMonthCumulative(inflate bool, timesInput ...time.Time) (TimeSeries, error) {
	newTimeSeries := TimeSeries{
		SeriesName: ts.SeriesName,
		ItemMap:    map[string]TimeItem{},
		IsFloat:    ts.IsFloat,
		Interval:   timeutil.IntervalMonth}
	tsMonth := ts.ToMonth(inflate)
	var min time.Time
	var max time.Time
	var err error
	if len(timesInput) > 0 {
		min, max, err = timeutil.TimeSliceMinMax(timesInput)
		if err != nil {
			return newTimeSeries, err
		}
	} else {
		min, max, err = timeutil.TimeSliceMinMax(tsMonth.Times(false))
		if err != nil {
			return newTimeSeries, err
		}
	}
	times := timeutil.TimeSeriesSlice(timeutil.IntervalMonth, []time.Time{min, max})
	cItems := []TimeItem{}
	for _, t := range times {
		rfc := t.Format(time.RFC3339)
		if item, ok := tsMonth.ItemMap[rfc]; ok {
			if len(cItems) > 0 {
				prevCItem := cItems[len(cItems)-1]
				cItems = append(cItems, TimeItem{
					SeriesName: newTimeSeries.SeriesName,
					IsFloat:    newTimeSeries.IsFloat,
					Time:       t,
					Value:      item.Value + prevCItem.Value,
					ValueFloat: item.ValueFloat + prevCItem.ValueFloat})
			} else {
				cItems = append(cItems, TimeItem{
					SeriesName: newTimeSeries.SeriesName,
					IsFloat:    newTimeSeries.IsFloat,
					Time:       t,
					Value:      item.Value,
					ValueFloat: item.ValueFloat})
			}
		} else {
			if len(cItems) > 0 {
				prevCItem := cItems[len(cItems)-1]
				cItems = append(cItems, TimeItem{
					SeriesName: newTimeSeries.SeriesName,
					IsFloat:    newTimeSeries.IsFloat,
					Time:       t,
					Value:      prevCItem.Value,
					ValueFloat: prevCItem.ValueFloat})
			} else {
				cItems = append(cItems, TimeItem{
					SeriesName: newTimeSeries.SeriesName,
					IsFloat:    newTimeSeries.IsFloat,
					Time:       t,
					Value:      0,
					ValueFloat: 0})
			}
		}
	}
	for _, cItem := range cItems {
		newTimeSeries.AddItems(cItem)
	}
	return newTimeSeries, nil
}

func (ts *TimeSeries) ToQuarter() TimeSeries {
	newTimeSeries := NewTimeSeries(ts.SeriesName)
	newTimeSeries.IsFloat = ts.IsFloat
	newTimeSeries.Interval = timeutil.IntervalQuarter
	for _, item := range ts.ItemMap {
		newTimeSeries.AddFloat64(timeutil.NewTimeMore(item.Time, 0).QuarterStart(), item.Float64())
	}
	return newTimeSeries
}

func (ts *TimeSeries) ToYear() TimeSeries {
	newTimeSeries := NewTimeSeries(ts.SeriesName)
	newTimeSeries.IsFloat = ts.IsFloat
	newTimeSeries.Interval = timeutil.IntervalYear
	for _, item := range ts.ItemMap {
		newTimeSeries.AddFloat64(timeutil.NewTimeMore(item.Time, 0).YearStart(), item.Float64())
	}
	return newTimeSeries
}

func AggregateSeries(series TimeSeries) TimeSeries {
	aggregate := NewTimeSeries(series.SeriesName)
	sortedItems := series.ItemsSorted()
	sum := int64(0)
	for _, atomicItem := range sortedItems {
		aggregateItem := TimeItem{
			SeriesName: atomicItem.SeriesName,
			Time:       atomicItem.Time,
			Value:      atomicItem.Value + sum,
		}
		sum = aggregateItem.Value
		aggregate.AddItems(aggregateItem)
	}
	return aggregate
}

func (ts *TimeSeries) Stats() point.PointSet {
	ps := point.NewPointSet()
	ps.IsFloat = ts.IsFloat
	for rfc3339, item := range ts.ItemMap {
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

func TimeSeriesDivide(numer, denom TimeSeries) (TimeSeries, error) {
	denomKeys := denom.Keys()
	ts := NewTimeSeries(denom.SeriesName)
	ts.IsFloat = true
	if numer.Interval == denom.Interval {
		ts.Interval = denom.Interval
	}
	ts.SeriesName = numer.SeriesName + " / " + denom.SeriesName
	for _, dKey := range denomKeys {
		nItem, nOk := numer.ItemMap[dKey]
		dItem, dOk := denom.ItemMap[dKey]
		if !nOk && !dOk {
			continue
		} else if !dOk || dItem.Value == 0 {
			return ts, fmt.Errorf("E_DENOM_MISSING_OR_ZERO TIME [%s] NUMERATOR [%v]",
				dKey, nItem.Value)
		}
		ts.AddItems(TimeItem{
			Time:       dItem.Time,
			ValueFloat: float64(nItem.Value) / float64(dItem.Value),
			IsFloat:    true,
		})
	}
	return ts, nil
}
