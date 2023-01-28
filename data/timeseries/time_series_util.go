package timeseries

import (
	"errors"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

func TimeSeriesMapMinMaxTimes(tsm map[string]TimeSeries) (time.Time, time.Time, error) {
	times := []time.Time{}
	for _, ts := range tsm {
		min, max := ts.MinMaxTimes()
		if !timeutil.NewTimeMore(min, 0).IsZeroAny() && !timeutil.NewTimeMore(max, 0).IsZeroAny() {
			times = append(times, min, max)
		}
	}
	return timeutil.TimeSliceMinMax(times)
}

func TimeSeriesMapMinMaxValues(tsm map[string]TimeSeries) (int64, int64, error) {
	minVal := int64(0)
	maxVal := int64(0)
	haveItems := false
	for _, ds := range tsm {
		if len(ds.ItemMap) == 0 {
			continue
		}
		minTry, maxTry := ds.MinMaxValues()
		if !haveItems {
			minVal = minTry
			maxVal = maxTry
			haveItems = true
			continue
		}
		if minTry < minVal {
			minVal = minTry
		}
	}
	if !haveItems {
		return minVal, maxVal, errors.New("timeSeriesMap no items")
	}
	return minVal, maxVal, nil
}
