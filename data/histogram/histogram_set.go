package histogram

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/grokify/gocharts/data/excelizeutil"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gocharts/data/timeseries"
	"github.com/grokify/simplego/time/timeutil"
	"github.com/grokify/simplego/type/stringsutil"
)

type HistogramSet struct {
	Name         string
	HistogramMap map[string]*Histogram
	KeyIsTime    bool
}

func NewHistogramSet(name string) *HistogramSet {
	return &HistogramSet{
		Name:         name,
		HistogramMap: map[string]*Histogram{}}
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

func (hset *HistogramSet) AddDateUidCount(dt time.Time, uid string, count int) {
	fName := dt.Format(time.RFC3339)
	hset.Add(fName, uid, count)
	if !hset.KeyIsTime {
		hset.KeyIsTime = true
	}
}

// Add provides an easy method to add a histogram bin name
// and count for an existing or new histogram in the set.
func (hset *HistogramSet) Add(histName, binName string, count int) {
	hist, ok := hset.HistogramMap[histName]
	if !ok {
		hist = NewHistogram(histName)
	}
	hist.Add(binName, count)
	hset.HistogramMap[histName] = hist
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
// Alias for `HistogramNames()`.
func (hset *HistogramSet) ItemNames() []string {
	return hset.HistogramNames()
}

// HistogramNames returns the number of histograms.
func (hset *HistogramSet) HistogramNames() []string {
	names := []string{}
	for name := range hset.HistogramMap {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// HistogramNameExists returns a boolean indicating if
// the supplied histogram name exists.
func (hset *HistogramSet) HistogramNameExists(histName string) bool {
	if _, ok := hset.HistogramMap[histName]; ok {
		return true
	}
	return false
}

// ValueSum returns the sum of all the histogram bin values.
func (hset *HistogramSet) ValueSum() int {
	valueSum := 0
	for _, hist := range hset.HistogramMap {
		valueSum += hist.ValueSum()
	}
	return valueSum
}

// BinNameExists returns a boolean indicating if a bin name
// exists in any histogram.
func (hset *HistogramSet) BinNameExists(binName string) bool {
	for _, hist := range hset.HistogramMap {
		if hist.BinNameExists(binName) {
			return true
		}
	}
	return false
}

// BinNames returns all the bin names used across all the
// histograms.
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

func (hset *HistogramSet) ToTimeSeriesDistinct() (timeseries.TimeSeries, error) {
	ds := timeseries.NewTimeSeries()
	ds.SeriesName = hset.Name
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

// WriteXLSXMatrix creates an XLSX file where the first column is the
// histogram name and the other columns are the bin names. This is
// useful for easy visualization of a table and also creating
// charts such as grouped bar charts.
func (hset *HistogramSet) WriteXLSXMatrix(filename, sheetName, histColName string) error {
	tbl, err := hset.TableMatrix(sheetName, histColName)
	if err != nil {
		return err
	}
	return tbl.WriteXLSX(filename, sheetName)
}

// TableMatrix returns a `*table.Table` where the first column is the
// histogram name and the other columns are the bin names. This is
// useful for easy visualization of a table and also creating
// charts such as grouped bar charts.
func (hset *HistogramSet) TableMatrix(tableName, histColName string) (*table.Table, error) {
	if len(strings.TrimSpace(tableName)) == 0 {
		tableName = strings.TrimSpace(hset.Name)
	}
	tbl := table.NewTable(tableName)

	if len(strings.TrimSpace(histColName)) == 0 {
		histColName = "Histogram Name"
	}

	binNames := hset.BinNames()
	tbl.Columns = append(tbl.Columns, histColName)
	tbl.Columns = append(tbl.Columns, binNames...)
	tbl.FormatMap = map[int]string{
		-1: table.FormatInt}
	if hset.KeyIsTime {
		tbl.FormatMap[0] = table.FormatTime
	} else {
		tbl.FormatMap[0] = table.FormatString
	}

	hnames := hset.HistogramNames()
	for _, hname := range hnames {
		row := []string{hname}
		hist, ok := hset.HistogramMap[hname]
		if !ok {
			return nil, fmt.Errorf("histogram name present without histogram [%s]", hname)
		}
		for _, binName := range binNames {
			if binVal, ok := hist.Bins[binName]; ok {
				row = append(row, strconv.Itoa(binVal))
			} else {
				row = append(row, "0")
			}
		}
		tbl.Rows = append(tbl.Rows, row)
	}
	return &tbl, nil
}

// WriteXLSX creates an XLSX file where the first column is the
// histogram name, the second column is the bin name and the
// third column is the bin count.
func (hset *HistogramSet) WriteXLSX(filename, sheetName, colName1, colName2, colNameCount string) error {
	// WriteXLSX writes a table as an Excel XLSX file with
	// row formatter option.
	f := excelize.NewFile()
	// Create a new sheet.

	if len(strings.TrimSpace(sheetName)) == 0 {
		sheetName = strings.TrimSpace(hset.Name)
	}
	if len(sheetName) == 0 {
		sheetName = "Sheet0"
	}
	index := f.NewSheet(sheetName)

	colName1 = strings.TrimSpace(colName1)
	if len(colName1) == 0 {
		colName1 = hset.Name
	}
	if len(colName1) == 0 {
		colName1 = "Column1"
	}
	colName2 = strings.TrimSpace(colName2)
	if len(colName1) == 0 {
		for _, fstats := range hset.HistogramMap {
			fstats.Name = strings.TrimSpace(fstats.Name)
			if len(fstats.Name) > 0 {
				colName2 = fstats.Name
				break
			}
		}
	}
	colNameCount = strings.TrimSpace(colNameCount)
	if len(colNameCount) == 0 {
		colNameCount = "Count"
	}
	header := []interface{}{colName1, colName2, colNameCount}

	excelizeutil.SetRowValues(f, sheetName, 0, header)
	var err error
	rowIdx := uint(1)
	for fstatsName, fstats := range hset.HistogramMap {
		fstatsNameDt := time.Now()
		if hset.KeyIsTime {
			fstatsNameDt, err = time.Parse(time.RFC3339, fstatsName)
			if err != nil {
				return err
			}
		}
		for binName, binCount := range fstats.Bins {
			var rowVals []interface{}
			if hset.KeyIsTime {
				rowVals = []interface{}{fstatsNameDt, binName, binCount}
			} else {
				rowVals = []interface{}{fstatsName, binName, binCount}
			}
			excelizeutil.SetRowValues(f, sheetName, rowIdx, rowVals)
			rowIdx++
		}
	}
	f.SetActiveSheet(index)
	// Delete Original Sheet
	f.DeleteSheet(f.GetSheetName(0))
	// Save xlsx file by the given path.
	return f.SaveAs(filename)

}

// DatetimeKeyToQuarter converts a HistogramSet
// by date to one by quarters.
func (hset *HistogramSet) DatetimeKeyToQuarter(name string) (*HistogramSet, error) {
	fsetQtr := NewHistogramSet(name)
	for rfc3339, hist := range hset.HistogramMap {
		dt, err := time.Parse(time.RFC3339, rfc3339)
		if err != nil {
			return fsetQtr, err
		}
		dt = timeutil.QuarterStart(dt)
		rfc3339Qtr := dt.Format(time.RFC3339)
		for binName, binCount := range hist.Bins {
			fsetQtr.Add(rfc3339Qtr, binName, binCount)
		}
	}
	return fsetQtr, nil
}

// DatetimeKeyCount returns a TimeSeries when
// the first key is a RFC3339 time and a sum of items
// is desired per time.
func (hset *HistogramSet) DatetimeKeyCount() (timeseries.TimeSeries, error) {
	ts := timeseries.NewTimeSeries()
	ts.SeriesName = hset.Name
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
	return ts.ToTable(hset.Name, "", countColName, timeseries.TimeFormatRFC3339), nil
}

func (hset *HistogramSet) HistogramSetTimeKeyCountWriteXLSX(filename string, interval timeutil.Interval, countColName string) error {
	tbl, err := hset.DatetimeKeyCountTable(interval, countColName)
	if err != nil {
		return err
	}
	tbl.FormatFunc = table.FormatTimeAndInts
	return table.WriteXLSX(filename, &tbl)
}
