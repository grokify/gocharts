package timeseries

import (
	"errors"
	"time"

	tu "github.com/grokify/mogo/time/timeutil"
)

type TimeSeriesSimple struct {
	Name        string
	DisplayName string
	Times       []time.Time
}

func NewTimeSeriesSimple(name, displayName string) TimeSeriesSimple {
	return TimeSeriesSimple{
		Name:        name,
		DisplayName: displayName,
		Times:       []time.Time{}}
}

func (tss *TimeSeriesSimple) ToTimeSeriesQuarter() TimeSeries {
	ts := NewTimeSeries(tss.Name)
	ts.SeriesName = tss.Name
	for _, t := range tss.Times {
		ts.AddItems(TimeItem{
			SeriesName: tss.Name,
			Time:       tu.QuarterStart(t),
			Value:      int64(1)})
	}
	return ts
}

type TimeSeriesFunnel struct {
	Series map[string]TimeSeriesSimple
	Order  []string
}

func (tsf *TimeSeriesFunnel) Times() []time.Time {
	times := []time.Time{}
	for _, s := range tsf.Series {
		times = append(times, s.Times...)
	}
	return times
}

func (tsf *TimeSeriesFunnel) TimesSorted() []time.Time {
	times := tsf.Times()
	return tu.Sort(times)
}

func (tsf *TimeSeriesFunnel) TimeSeriesSetByQuarter() (TimeSeriesSet, error) {
	dss := TimeSeriesSet{Order: tsf.Order}
	seriesMapQuarter := map[string]TimeSeries{}

	allTimes := []time.Time{}
	for _, s := range tsf.Series {
		allTimes = append(allTimes, s.Times...)
	}

	if len(allTimes) == 0 {
		return dss, errors.New("no times")
	}
	earliest, err := tu.Earliest(allTimes, false)
	if err != nil {
		return dss, err
	}
	latest, err := tu.Latest(allTimes, false)
	if err != nil {
		return dss, err
	}
	earliestQuarter := tu.QuarterStart(earliest)
	latestQuarter := tu.QuarterStart(latest)

	sliceQuarter := tu.QuarterSlice(earliestQuarter, latestQuarter)
	dss.Times = sliceQuarter

	for name, tss := range tsf.Series {
		timeSeries := tss.ToTimeSeriesQuarter()
		timeSeries.SeriesName = tss.Name
		for _, q := range sliceQuarter {
			q = q.UTC()
			rfc := q.Format(time.RFC3339)
			if _, ok := timeSeries.ItemMap[rfc]; !ok {
				timeSeries.AddItems(TimeItem{
					SeriesName: tss.Name,
					Time:       q,
					Value:      int64(0)})
			}
		}
		seriesMapQuarter[name] = timeSeries
	}
	dss.Series = seriesMapQuarter
	return dss, nil
}
