package timeseries

import (
	"strconv"

	"github.com/grokify/mogo/time/timeutil"
)

// TimeStats is used to generate unique counts stats
// for an array of with time ane names.
type TimeStats struct {
	Items []TimeItem
}

func (ts *TimeStats) UniqueCountsByQuarter() map[string]int {
	// quarter-value-count
	wip := map[string]map[string]int{}
	for _, item := range ts.Items {
		q := strconv.Itoa(timeutil.YearQuarterForTime(item.Time))
		v := item.SeriesName
		if _, ok := wip[q]; !ok {
			wip[q] = map[string]int{}
		}
		if _, ok := wip[q][v]; !ok {
			wip[q][v] = 0
		}
		wip[q][v] += 1
	}
	return mapSSIToMSICounts(wip)
}

func (ts *TimeStats) UniqueCountsByMonth() map[string]int {
	// month-value-count
	wip := map[string]map[string]int{}
	for _, item := range ts.Items {
		m := item.Time.Format(timeutil.DT6)
		v := item.SeriesName
		if _, ok := wip[m]; !ok {
			wip[m] = map[string]int{}
		}
		if _, ok := wip[m][v]; !ok {
			wip[m][v] = 0
		}
		wip[m][v] += 1
	}
	return mapSSIToMSICounts(wip)
}

func mapSSIToMSICounts(in map[string]map[string]int) map[string]int {
	out := map[string]int{}
	for k, v := range in {
		out[k] = len(v)
	}
	return out
}
