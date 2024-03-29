package timeseries

import (
	"errors"
	"sort"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

/*
func (set *TimeSeriesSet) NewTimesLowerBound(times ...time.Time) (TimeSeriesSets, error) {
	min, _ := set.MinMaxTimes()
	times = append(times, min)
	timeSlice := timeutil.Times(times)
	sort.Sort(timeSlice)
	timeSlice = timeSlice.Dedupe()
	sets := NewTimeSeriesSets("time sets by times")
	sets.KeyIsTime = true
	for seriesName, ts := range set.Series {
		for _, timeItem := range ts.ItemMap {
			timeBucket, err := timeSlice.RangeLower(timeItem.Time, true)
			if err != nil {
				return sets, err
			}
			timeItem.SeriesName = seriesName
			timeItem.SeriesSetName = timeBucket.Format(time.RFC3339)
			timeItem.Time = timeBucket
			sets.AddItems(timeItem)
		}
	}
	return sets, nil
}
*/

func (set *TimeSeriesSet) ToSetTimesRangeUpper(inclusive bool, times ...time.Time) (TimeSeriesSet, error) {
	_, max := set.MinMaxTimes()
	times = append(times, max)
	timeSlice := timeutil.Times(times)
	sort.Sort(timeSlice)
	timeSlice = timeSlice.Dedupe()
	newSet := NewTimeSeriesSet(set.Name)
	for seriesName, series := range set.Series {
		for _, timeItem := range series.ItemMap {
			if len(timeItem.SeriesName) == 0 {
				timeItem.SeriesName = seriesName
			}
			if timeItem.SeriesName != seriesName {
				return newSet, errors.New("timeItem.SeriesName != TimeSeriesSet seriesName")
			}
			tRangeUpper, err := timeSlice.RangeUpper(timeItem.Time, inclusive)
			if err != nil {
				panic("time item greater than set max time")
			}
			timeItem.Time = tRangeUpper
			newSet.AddItems(timeItem)
		}
	}

	return newSet, nil
}
