package frequency

import (
	"strings"
	"time"

	"github.com/grokify/gocharts/data/statictimeseries"
	"github.com/grokify/gotilla/time/timeutil"
)

type FrequencySets struct {
	FrequencySetMap map[string]FrequencySet
}

func NewFrequencySets() FrequencySets {
	return FrequencySets{FrequencySetMap: map[string]FrequencySet{}}
}

func (fsets *FrequencySets) Add(key1, key2, uid string, trimSpace bool) {
	if trimSpace {
		key1 = strings.TrimSpace(key1)
		key2 = strings.TrimSpace(key2)
		uid = strings.TrimSpace(uid)
	}
	fset, ok := fsets.FrequencySetMap[key1]
	if !ok {
		fset = NewFrequencySet(key1)
	}
	fset.AddString(key2, uid)
	fsets.FrequencySetMap[key1] = fset
}

func (fsets *FrequencySets) Flatten(name string) FrequencySet {
	fsetFlat := NewFrequencySet(name)
	for _, fset := range fsets.FrequencySetMap {
		for k2, fstats := range fset.FrequencyMap {
			for item, count := range fstats.Items {
				fsetFlat.AddStringMore(k2, item, count)
			}
		}
	}
	return fsetFlat
}

type FrequencySet struct {
	Name         string
	FrequencyMap map[string]FrequencyStats
}

func NewFrequencySet(name string) FrequencySet {
	return FrequencySet{
		Name:         name,
		FrequencyMap: map[string]FrequencyStats{}}
}

func (fss *FrequencySet) AddStringMore(frequencyName, itemName string, count int) {
	fs, ok := fss.FrequencyMap[frequencyName]
	if !ok {
		fs = NewFrequencyStats(frequencyName)
	}
	fs.AddStringMore(itemName, count)
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

// FrequencySetTimeKeyCounts returns a DataSeries when
// the first key is a RFC3339 time and a count of items
// is desired per time.
func FrequencySetTimeKeyCounts(fset FrequencySet) (statictimeseries.DataSeries, error) {
	ds := statictimeseries.NewDataSeries()
	ds.SeriesName = fset.Name
	for rfc3339, fstats := range fset.FrequencyMap {
		dtCount := len(fstats.Items)
		dt, err := time.Parse(time.RFC3339, rfc3339)
		if err != nil {
			return ds, err
		}
		di := statictimeseries.DataItem{
			SeriesName: fset.Name,
			Time:       dt,
			Value:      int64(dtCount)}
		ds.AddItem(di)
	}
	return ds, nil
}
