package histogram

import (
	"fmt"
	"strings"
	"time"

	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/stringsutil"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gocharts/v2/data/table/format"
	"github.com/grokify/gocharts/v2/data/timeseries"
)

type HistogramSet struct {
	Name         string
	HistogramMap map[string]*Histogram
	KeyIsTime    bool
	Order        []string
}

func NewHistogramSet(name string) *HistogramSet {
	return &HistogramSet{
		Name:         name,
		HistogramMap: map[string]*Histogram{},
		Order:        []string{}}
}

func NewHistogramSetWithData(name string, data map[string]map[string]int) *HistogramSet {
	hset := &HistogramSet{
		Name:         name,
		HistogramMap: map[string]*Histogram{}}
	for statsName, statsData := range data {
		for statsItemName, statsItemValue := range statsData {
			hset.Add(statsName, statsItemName, statsItemValue)
		}
	}
	return hset
}

func (hset *HistogramSet) AddDateUIDCount(dt time.Time, uid string, count int) {
	fName := dt.Format(time.RFC3339)
	hset.Add(fName, uid, count)
	if !hset.KeyIsTime {
		hset.KeyIsTime = true
	}
}

// Add provides an easy method to add a histogram bin name
// and count for an existing or new histogram in the set.
func (hset *HistogramSet) Add(histName, binName string, binCount int) {
	hist, ok := hset.HistogramMap[histName]
	if !ok {
		hist = NewHistogram(histName)
	}
	hist.Add(binName, binCount)
	hset.HistogramMap[histName] = hist
}

// BinSetCounts returns a ap where the key is the count of bins and the string is the set name.
func (hset *HistogramSet) BinParentCounts() map[uint]map[string]uint {
	out := map[uint]map[string]uint{}
	wip := map[string]map[string]uint{} //
	for hsetName, hist := range hset.HistogramMap {
		for binName := range hist.Bins {
			if wip[binName] == nil {
				wip[binName] = map[string]uint{}
			}
			wip[binName][hsetName]++
		} // cust | aws
	}
	for binName, mapHsetCount := range wip {
		hsetCount := uint(len(mapHsetCount))
		if out[hsetCount] == nil {
			out[hsetCount] = map[string]uint{}
		}
		out[hsetCount][binName]++
	}
	return out
}

// BinValue the value of a bin.
func (hset *HistogramSet) BinValue(histName, binName string) int {
	if h, ok := hset.HistogramMap[histName]; !ok {
		return 0
	} else if h == nil {
		return 0
	} else if v, ok := h.Bins[binName]; !ok {
		return 0
	} else {
		return v
	}
}

// ItemCount returns the number of histograms.
func (hset *HistogramSet) ItemCount() uint {
	return uint(len(hset.HistogramMap))
}

// ItemCounts returns the number of histograms.
func (hset *HistogramSet) ItemCounts() *Histogram {
	histCount := NewHistogram("histogram counts counts")
	for histName, hist := range hset.HistogramMap {
		histCount.Bins[histName] = len(hist.Bins)
	}
	histCount.Inflate()
	return histCount
}

// ItemNames returns the number of histograms.
func (hset *HistogramSet) ItemNames() []string {
	return maputil.Keys(hset.HistogramMap)
}

// HistogramNameExists returns a boolean indicating if the supplied histogram name exists.
func (hset *HistogramSet) HistogramNameExists(histName string) bool {
	if _, ok := hset.HistogramMap[histName]; ok {
		return true
	}
	return false
}

// Sum returns the sum of all the histogram bin values.
func (hset *HistogramSet) Sum() int {
	valueSum := 0
	for _, hist := range hset.HistogramMap {
		valueSum += hist.Sum()
	}
	return valueSum
}

// BinNameExists returns a boolean indicating if a bin name exists in any histogram.
func (hset *HistogramSet) BinNameExists(binName string) bool {
	for _, hist := range hset.HistogramMap {
		if hist.BinNameExists(binName) {
			return true
		}
	}
	return false
}

// BinNames returns all the bin names used across all the histograms.
func (hset *HistogramSet) BinNames() []string {
	binNames := []string{}
	for _, hist := range hset.HistogramMap {
		binNames = append(binNames, hist.BinNames()...)
	}
	return stringsutil.SliceCondenseSpace(binNames, true, true)
}

// HistogramBinNames returns the bin names for a single
// histogram whose name is provided as a function parameter.
func (hset *HistogramSet) HistogramBinNames(setName string) []string {
	if hist, ok := hset.HistogramMap[setName]; ok {
		return hist.BinNames()
	}
	return []string{}
}

// LeafStats returns a histogram by combining the histogram
// bins across histograms, removing the histogram distinction.
func (hset *HistogramSet) LeafStats(name string) *Histogram {
	if len(name) == 0 {
		name = "leaf stats"
	}
	setLeafStats := NewHistogram(name)
	for _, hist := range hset.HistogramMap {
		for binName, binCount := range hist.Bins {
			setLeafStats.Add(binName, binCount)
		}
	}
	return setLeafStats
}

func (hset *HistogramSet) Map() map[string]map[string]int {
	out := map[string]map[string]int{}
	for histName, hist := range hset.HistogramMap {
		if _, ok := out[histName]; !ok {
			out[histName] = map[string]int{}
		}
		for binName, binCount := range hist.Bins {
			out[histName][binName] += binCount
		}
	}
	return out
}

func (hset *HistogramSet) MapAdd(m map[string]map[string]int) {
	for histName, histMap := range m {
		for binName, binCount := range histMap {
			hset.Add(histName, binName, binCount)
		}
	}
}

func (hset *HistogramSet) ToTimeSeriesDistinct() (timeseries.TimeSeries, error) {
	ds := timeseries.NewTimeSeries(hset.Name)
	for rfc3339, hist := range hset.HistogramMap {
		dt, err := time.Parse(time.RFC3339, rfc3339)
		if err != nil {
			return ds, err
		}
		ds.AddItems(timeseries.TimeItem{
			SeriesName: hset.Name,
			Time:       dt,
			Value:      int64(len(hist.Bins))})
	}
	return ds, nil
}

// DatetimeKeyToQuarter converts a HistogramSet by date to one by quarters.
func (hset *HistogramSet) DatetimeKeyToQuarter(name string) (*HistogramSet, error) {
	fsetQtr := NewHistogramSet(name)
	for rfc3339, hist := range hset.HistogramMap {
		dt, err := time.Parse(time.RFC3339, rfc3339)
		if err != nil {
			return fsetQtr, err
		}
		dt = timeutil.NewTimeMore(dt, 0).QuarterStart()
		rfc3339Qtr := dt.Format(time.RFC3339)
		for binName, binCount := range hist.Bins {
			fsetQtr.Add(rfc3339Qtr, binName, binCount)
		}
	}
	return fsetQtr, nil
}

// DatetimeKeyCount returns a TimeSeries when the first key is a RFC3339 time
// and a sum of items is desired per time.
func (hset *HistogramSet) DatetimeKeyCount() (timeseries.TimeSeries, error) {
	ts := timeseries.NewTimeSeries(hset.Name)
	for rfc3339, hist := range hset.HistogramMap {
		dt, err := time.Parse(time.RFC3339, rfc3339)
		if err != nil {
			return ts, err
		}
		ts.AddItems(timeseries.TimeItem{
			SeriesName: hset.Name,
			Time:       dt,
			Value:      int64(len(hist.Bins))})
	}
	return ts, nil
}

func (hset *HistogramSet) DatetimeKeyCountTable(interval timeutil.Interval, countColName string) (table.Table, error) {
	ts, err := hset.DatetimeKeyCount()
	if err != nil {
		return table.NewTable(hset.Name), err
	}
	ts.Interval = interval
	if len(strings.TrimSpace(countColName)) == 0 {
		countColName = "Count"
	}
	return ts.Table(hset.Name, "", countColName, timeseries.TimeFormatRFC3339), nil
}

func (hset *HistogramSet) HistogramSetTimeKeyCountWriteXLSX(filename string, interval timeutil.Interval, countColName string) error {
	tbl, err := hset.DatetimeKeyCountTable(interval, countColName)
	if err != nil {
		return err
	}
	tbl.FormatFunc = format.FormatTimeAndInts
	return table.WriteXLSX(filename, []*table.Table{&tbl})
}

func (hset *HistogramSet) String(histNotSet, binNotSet string) string {
	var parts []string
	for hName, h := range hset.HistogramMap {
		if strings.TrimSpace(hName) == "" {
			hName = histNotSet
		}
		bNames := h.BinNames()
		var counts []string
		for _, bName := range bNames {
			count, ok := h.Bins[bName]
			if !ok {
				count = 0
			}
			if strings.TrimSpace(bName) == "" {
				bName = binNotSet
			}
			counts = append(counts, fmt.Sprintf("%s: %d", bName, count))
		}
		if len(counts) > 0 {
			parts = append(parts, fmt.Sprintf("%s (%s)", hName, strings.Join(counts, ", ")))
		} else {
			parts = append(parts, hName)
		}
	}
	return strings.Join(parts, ", ")
}
