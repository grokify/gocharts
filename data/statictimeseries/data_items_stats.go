// statictimeseriesdata provides tools for adding and formatting
// static time series data for reporting purposes.
package statictimeseries

import (
	"strconv"

	tu "github.com/grokify/gotilla/time/timeutil"
)

// DataItemsStats is used to generate unique counts stats
// for an array of with time ane names.
type TimeStats struct {
	Items []DataItem
}

func (ts *TimeStats) UniqueCountsByQuarter() map[string]int {
	// quarter-value-count
	wip := map[string]map[string]int{}
	for _, item := range ts.Items {
		q := strconv.Itoa(int(tu.QuarterInt32ForTime(item.Time)))
		v := item.SeriesName
		if _, ok := wip[q]; !ok {
			wip[q] = map[string]int{}
		}
		if _, ok := wip[q][v]; !ok {
			wip[q][v] = 0
		}
		wip[q][v] += 1
	}
	return MSS2MS(wip)
}

func (ts *TimeStats) UniqueCountsByMonth() map[string]int {
	// month-value-count
	wip := map[string]map[string]int{}
	for _, item := range ts.Items {
		m := item.Time.Format(tu.DT6)
		v := item.SeriesName
		if _, ok := wip[m]; !ok {
			wip[m] = map[string]int{}
		}
		if _, ok := wip[m][v]; !ok {
			wip[m][v] = 0
		}
		wip[m][v] += 1
	}
	return MSS2MS(wip)
}

func MSS2MS(in map[string]map[string]int) map[string]int {
	out := map[string]int{}
	for k, v := range in {
		out[k] = len(v)
	}
	return out
}
