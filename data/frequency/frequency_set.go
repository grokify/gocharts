package frequency

import (
	"sort"
	"strings"
	"time"

	"github.com/grokify/gocharts/data/statictimeseries"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gotilla/time/timeutil"
)

type FrequencySet struct {
	Name         string
	FrequencyMap map[string]FrequencyStats
}

func NewFrequencySet(name string) FrequencySet {
	return FrequencySet{
		Name:         name,
		FrequencyMap: map[string]FrequencyStats{}}
}

func (fss *FrequencySet) AddDateUidCount(dt time.Time, uid string, count int) {
	fName := dt.Format(time.RFC3339)
	fss.AddStringMore(fName, uid, count)
}

func (fss *FrequencySet) AddStringMore(frequencyName, uid string, count int) {
	fs, ok := fss.FrequencyMap[frequencyName]
	if !ok {
		fs = NewFrequencyStats(frequencyName)
	}
	fs.AddStringMore(uid, count)
	fss.FrequencyMap[frequencyName] = fs
}

func (fss *FrequencySet) AddString(frequencyName, itemName string) {
	fs, ok := fss.FrequencyMap[frequencyName]
	if !ok {
		fs = NewFrequencyStats(frequencyName)
	}
	fs.AddString(itemName)
	fss.FrequencyMap[frequencyName] = fs
}

func (fss *FrequencySet) Names() []string {
	names := []string{}
	for name := range fss.FrequencyMap {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (fss *FrequencySet) TotalCount() uint64 {
	totalCount := uint64(0)
	for _, fstats := range fss.FrequencyMap {
		totalCount += fstats.TotalCount()
	}
	return totalCount
}

func (fss *FrequencySet) ToDataSeriesDistinct() (statictimeseries.DataSeries, error) {
	ds := statictimeseries.NewDataSeries()
	ds.SeriesName = fss.Name
	for rfc3339, fs := range fss.FrequencyMap {
		dt, err := time.Parse(time.RFC3339, rfc3339)
		if err != nil {
			return ds, err
		}
		ds.AddItem(statictimeseries.DataItem{
			SeriesName: fss.Name,
			Time:       dt,
			Value:      int64(len(fs.Items))})
	}
	return ds, nil
}

// FrequencySetDatetimeToQuarterUnique converts a FrequencySet
// by date to one by quarter.s.
func FrequencySetDatetimeToQuarter(name string, fsetIn FrequencySet) (FrequencySet, error) {
	fsetQtr := NewFrequencySet(name)
	for rfc3339, fstats := range fsetIn.FrequencyMap {
		dt, err := time.Parse(time.RFC3339, rfc3339)
		if err != nil {
			return fsetQtr, err
		}
		dt = timeutil.QuarterStart(dt)
		rfc3339Qtr := dt.Format(time.RFC3339)
		for item, count := range fstats.Items {
			fsetQtr.AddStringMore(rfc3339Qtr, item, count)
		}
	}
	return fsetQtr, nil
}

// FrequencySetTimeKeyCount returns a DataSeries when
// the first key is a RFC3339 time and a sum of items
// is desired per time.
func FrequencySetTimeKeyCount(fset FrequencySet) (statictimeseries.DataSeries, error) {
	ds := statictimeseries.NewDataSeries()
	ds.SeriesName = fset.Name
	for rfc3339, fstats := range fset.FrequencyMap {
		dt, err := time.Parse(time.RFC3339, rfc3339)
		if err != nil {
			return ds, err
		}
		ds.AddItem(statictimeseries.DataItem{
			SeriesName: fset.Name,
			Time:       dt,
			Value:      int64(len(fstats.Items))})
	}
	return ds, nil
}

func FrequencySetTimeKeyCountTable(fset FrequencySet, interval timeutil.Interval, countColName string) (table.Table, error) {
	ds, err := FrequencySetTimeKeyCount(fset)
	if err != nil {
		return table.NewTable(), err
	}
	ds.Interval = interval
	countColName = strings.TrimSpace(countColName)
	if len(countColName) == 0 {
		countColName = "Count"
	}
	return statictimeseries.DataSeriesToTable(ds, countColName, statictimeseries.TimeFormatRFC3339), nil
}

func FrequencySetTimeKeyCountWriteXLSX(filename string, fset FrequencySet, interval timeutil.Interval, countColName string) error {
	tbl, err := FrequencySetTimeKeyCountTable(fset, interval, countColName)
	if err != nil {
		return err
	}
	tbl.FormatFunc = table.FormatTimeAndInts
	return table.WriteXLSX(filename, &tbl)
}
