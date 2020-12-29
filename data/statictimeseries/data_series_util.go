package statictimeseries

import (
	"time"

	"github.com/grokify/simplego/time/timeutil"
	"github.com/pkg/errors"
)

func DataSeriesMapMinMaxTimes(dsm map[string]DataSeries) (time.Time, time.Time, error) {
	times := []time.Time{}
	for _, ds := range dsm {
		min, max := ds.MinMaxTimes()
		if !timeutil.TimeIsZeroAny(min) && !timeutil.TimeIsZeroAny(max) {
			times = append(times, min, max)
		}
	}
	return timeutil.TimeSliceMinMax(times)
}

func DataSeriesMapMinMaxValues(dsm map[string]DataSeries) (int64, int64, error) {
	minVal := int64(0)
	maxVal := int64(0)
	haveItems := false
	for _, ds := range dsm {
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
		return minVal, maxVal, errors.New("dataSeriesMap no items")
	}
	return minVal, maxVal, nil
}
