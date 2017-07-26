// statictimeseriesdata provides tools for adding and formatting
// static time series data for reporting purposes.
package statictimeseriesdata

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/grokify/gotilla/time/timeutil"
	"github.com/grokify/gotilla/type/maputil"
)

// DataSeriesSet is used to prepare histogram data for
// static timeseries charts, adding zero values for
// time slots as necessary. Usage is to create to call:
// NewDataSeriesSet("quarter"), AddItem() and then Inflate()
type DataSeriesSet struct {
	SourceSeriesMap map[string]DataSeries
	OutputSeriesMap map[string]DataSeries
	SeriesIntervals SeriesIntervals
}

func NewDataSeriesSet(interval string) (DataSeriesSet, error) {
	set := DataSeriesSet{
		SourceSeriesMap: map[string]DataSeries{},
		OutputSeriesMap: map[string]DataSeries{},
		SeriesIntervals: SeriesIntervals{Interval: interval}}
	err := set.SeriesIntervals.CheckInterval()
	return set, err
}

type DataSeries struct {
	SeriesName string
	ItemMap    map[string]DataItem
}

func NewDataSeries() DataSeries {
	return DataSeries{ItemMap: map[string]DataItem{}}
}

type DataItem struct {
	SeriesName string
	Time       time.Time
	Value      int64
}

func (set *DataSeriesSet) SeriesNamesSorted() []string {
	return maputil.StringKeysSorted(set.SourceSeriesMap)
}

func (set *DataSeriesSet) AddItem(item DataItem) {
	item.SeriesName = strings.TrimSpace(item.SeriesName)
	if _, ok := set.SourceSeriesMap[item.SeriesName]; !ok {
		set.SourceSeriesMap[item.SeriesName] =
			DataSeries{
				SeriesName: item.SeriesName,
				ItemMap:    map[string]DataItem{}}
	}
	series := set.SourceSeriesMap[item.SeriesName]
	series.AddItem(item)
	set.SourceSeriesMap[item.SeriesName] = series
}

func (series *DataSeries) AddItem(item DataItem) {
	if series.ItemMap == nil {
		series.ItemMap = map[string]DataItem{}
	}
	rfc := item.Time.Format(time.RFC3339)
	if _, ok := series.ItemMap[rfc]; !ok {
		series.ItemMap[rfc] = item
	} else {
		existingItem := series.ItemMap[rfc]
		existingItem.Value += item.Value
		series.ItemMap[rfc] = existingItem
	}
}

func (set *DataSeriesSet) Inflate() error {
	err := set.InflateSource()
	if err != nil {
		return err
	}
	return set.InflateOutput()
}

func (set *DataSeriesSet) InflateSource() error {
	for _, series := range set.SourceSeriesMap {
		set.SeriesIntervals.ProcItemsMap(series.ItemMap)
	}
	return set.SeriesIntervals.Inflate()
}

func (set *DataSeriesSet) InflateOutput() error {
	for seriesName, series := range set.SourceSeriesMap {
		output, err := set.BuildOutputSeries(series)
		if err != nil {
			return err
		}
		set.OutputSeriesMap[seriesName] = output
	}
	return nil
}

func (set *DataSeriesSet) BuildOutputSeries(source DataSeries) (DataSeries, error) {
	output := NewDataSeries()
	if set.SeriesIntervals.Interval == "quarter" {
		for _, item := range source.ItemMap {
			output.SeriesName = item.SeriesName
			qtrStart, err := timeutil.QuarterStart(item.Time)
			if err != nil {
				return output, err
			}
			output.AddItem(DataItem{
				SeriesName: item.SeriesName,
				Time:       qtrStart,
				Value:      item.Value})
		}
	}
	for _, dt := range set.SeriesIntervals.CanonicalSeries {
		output.AddItem(DataItem{
			SeriesName: output.SeriesName,
			Time:       dt,
			Value:      0})
	}
	return output, nil
}

// SortedItems returns sorted DataItems. This currently uses
// a simple tring sort on RFC3339 times. For dates that are not
// handled properly this way, this can be enhanced to use more
// proper comparison
func (series *DataSeries) SortedItems() []DataItem {
	itemsSorted := []DataItem{}
	timesSorted := maputil.StringKeysSorted(series.ItemMap)
	for _, rfc3339 := range timesSorted {
		itemsSorted = append(itemsSorted, series.ItemMap[rfc3339])
	}
	return itemsSorted
}

type SeriesIntervals struct {
	Interval        string
	Max             time.Time
	Min             time.Time
	CanonicalSeries []time.Time
}

func (ival *SeriesIntervals) CheckInterval() error {
	ival.Interval = strings.ToLower(strings.TrimSpace(ival.Interval))
	if ival.Interval != "quarter" {
		return errors.New(fmt.Sprintf("Unsupported interval [%v]", ival.Interval))
	}
	return nil
}

func (ival *SeriesIntervals) AreEndpointsSet() bool {
	if ival.Max.IsZero() || ival.Min.IsZero() {
		return false
	}
	return true
}

func (ival *SeriesIntervals) ProcItemsMap(itemMap map[string]DataItem) {
	for _, dataItem := range itemMap {
		dt := dataItem.Time
		if !ival.AreEndpointsSet() {
			ival.Max = dt
			ival.Min = dt
			continue
		}
		if timeutil.IsGreaterThan(dt, ival.Max) {
			ival.Max = dt
		}
		if timeutil.IsLessThan(dt, ival.Min) {
			ival.Min = dt
		}
	}
}

func (ival *SeriesIntervals) Inflate() error {
	err := ival.CheckInterval()
	if err != nil {
		return err
	}
	if ival.Interval == "quarter" {
		return ival.BuildCanonicalByQuarter()
	} else {
		return errors.New(fmt.Sprintf("Interval [%v] not found", ival.Interval))
	}
	return nil
}

func (ival *SeriesIntervals) BuildCanonicalByQuarter() error {
	if !ival.AreEndpointsSet() {
		return errors.New("Cannot build canonical dates without initialized dates.")
	}
	ival.Max = ival.Max.UTC()
	max, err := timeutil.QuarterStart(ival.Max)
	if err != nil {
		return err
	}
	ival.Max = max
	ival.Min = ival.Min.UTC()
	min, err := timeutil.QuarterStart(ival.Min)
	if err != nil {
		return err
	}
	ival.Min = min

	canonicalSeries := []time.Time{}
	curTime := ival.Min
	for timeutil.IsLessThan(curTime, ival.Max) {
		canonicalSeries = append(canonicalSeries, curTime)
		curTime = timeutil.TimeDt6AddNMonths(curTime, 3)
	}
	ival.CanonicalSeries = canonicalSeries
	return nil
}
