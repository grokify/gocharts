package timeseries

import (
	"errors"
	"time"

	tu "github.com/grokify/simplego/time/timeutil"
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

func (tss *TimeSeriesSimple) ToDataSeriesQuarter() DataSeries {
	ds := NewDataSeries()
	ds.SeriesName = tss.Name
	for _, t := range tss.Times {
		ds.AddItem(TimeItem{
			SeriesName: tss.Name,
			Time:       tu.QuarterStart(t),
			Value:      int64(1)})
	}
	return ds
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

func (tsf *TimeSeriesFunnel) DataSeriesSetByQuarter() (DataSeriesSet, error) {
	dss := DataSeriesSet{Order: tsf.Order}
	seriesMapQuarter := map[string]DataSeries{}

	allTimes := []time.Time{}
	for _, s := range tsf.Series {
		allTimes = append(allTimes, s.Times...)
	}

	if len(allTimes) == 0 {
		return dss, errors.New("No times")
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
		dataSeries := tss.ToDataSeriesQuarter()
		dataSeries.SeriesName = tss.Name
		for _, q := range sliceQuarter {
			q = q.UTC()
			rfc := q.Format(time.RFC3339)
			if _, ok := dataSeries.ItemMap[rfc]; !ok {
				dataSeries.AddItem(TimeItem{
					SeriesName: tss.Name,
					Time:       q,
					Value:      int64(0)})
			}
		}
		seriesMapQuarter[name] = dataSeries
	}
	dss.Series = seriesMapQuarter
	return dss, nil
}
