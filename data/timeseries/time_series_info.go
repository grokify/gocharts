package timeseries

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/grokify/mogo/sort/sortutil"
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/time/timeutil"
)

// Keys returns a sorted listed of Item keys.
func (ts *TimeSeries) Keys() []string {
	keys := []string{}
	for key := range ts.ItemMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// ItemsSorted returns sorted TimeItems. This currently uses
// a simple string sort on RFC3339 times.
func (ts *TimeSeries) ItemsSorted() []TimeItem {
	keys := ts.Keys()
	itemsSorted := []TimeItem{}
	for _, key := range keys {
		if item, ok := ts.ItemMap[key]; !ok {
			panic(fmt.Sprintf("KEY_NOT_FOUND [%s]", key))
		} else {
			itemsSorted = append(itemsSorted, item)
		}
	}
	return itemsSorted
}

func (ts *TimeSeries) LastItem(skipIfTimePartialValueLessPrev bool) (TimeItem, error) {
	items := ts.ItemsSorted()
	if len(items) == 0 {
		return TimeItem{}, ErrNoTimeItem
	}
	if len(items) == 1 {
		return items[0], nil
	}
	itemLast := items[len(items)-1]
	if skipIfTimePartialValueLessPrev {
		itemPrev := items[len(items)-2]
		dtNow := time.Now().UTC()
		if ts.Interval == timeutil.IntervalMonth {
			dtNow = month.MonthStart(dtNow, 0)
		}
		if itemLast.Time.Equal(dtNow) {
			if itemLast.Int64() > itemPrev.Int64() {
				return itemLast, nil
			} else {
				return itemPrev, nil
			}
		}
	}
	return itemLast, nil
}

func (ts *TimeSeries) minMaxValuesInt64Only() (int64, int64) {
	int64s := []int64{}
	for _, item := range ts.ItemMap {
		int64s = append(int64s, item.Value)
	}
	if len(int64s) == 0 {
		return 0, 0
	}
	sortutil.Slice(int64s)
	return int64s[0], int64s[len(int64s)-1]
}

func (ts *TimeSeries) minMaxValuesFloat64Only() (float64, float64) {
	float64s := []float64{}
	for _, item := range ts.ItemMap {
		float64s = append(float64s, item.ValueFloat)
	}
	if len(float64s) == 0 {
		return 0, 0
	}
	float64slice := sort.Float64Slice(float64s)
	sort.Sort(float64slice)

	return float64slice[0], float64slice[len(float64slice)-1]
}

func (ts *TimeSeries) MinMaxValues() (int64, int64) {
	if ts.IsFloat {
		min, max := ts.minMaxValuesFloat64Only()
		return int64(min), int64(max)
	}
	return ts.minMaxValuesInt64Only()
}

func (ts *TimeSeries) MinMaxValuesFloat64() (float64, float64) {
	if ts.IsFloat {
		return ts.minMaxValuesFloat64Only()
	}
	min, max := ts.minMaxValuesInt64Only()
	return float64(min), float64(max)
}

func (ts *TimeSeries) MinValue() int64 {
	min, _ := ts.MinMaxValues()
	return min
}

func (ts *TimeSeries) MaxValue() int64 {
	_, max := ts.MinMaxValues()
	return max
}

func (ts *TimeSeries) OneItemMaxValue() (TimeItem, error) {
	max := TimeItem{}
	if len(ts.ItemMap) == 0 {
		return max, errors.New("empty set has no max value item")
	}
	first := true
	for _, item := range ts.ItemMap {
		if first {
			max = item
			first = false
		}
		if ts.IsFloat && item.ValueFloat > max.ValueFloat {
			max = item
		} else if item.Value > max.Value {
			max = item
		}
	}
	return max, nil
}

func (ts *TimeSeries) TimeSeries(interval timeutil.Interval) []time.Time {
	return timeutil.TimeSeriesSlice(interval, ts.ItemTimes())
}

func (ts *TimeSeries) ItemTimes() []time.Time {
	var times []time.Time
	for _, item := range ts.ItemMap {
		times = append(times, item.Time)
	}
	return times
}

func (ts *TimeSeries) MinMaxTimes() (time.Time, time.Time) {
	return timeutil.SliceMinMax(ts.ItemTimes())
}
