package frequency

import (
	"sort"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/grokify/gocharts/data/excelizeutil"
	"github.com/grokify/gocharts/data/statictimeseries"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/simplego/time/timeutil"
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
	fset := &HistogramSet{
		Name:         name,
		HistogramMap: map[string]*Histogram{}}
	for statsName, statsData := range data {
		for statsItemName, statsItemValue := range statsData {
			fset.Add(statsName, statsItemName, statsItemValue)
		}
	}
	return fset
}

func (fset *HistogramSet) AddDateUidCount(dt time.Time, uid string, count int) {
	fName := dt.Format(time.RFC3339)
	fset.Add(fName, uid, count)
	if !fset.KeyIsTime {
		fset.KeyIsTime = true
	}
}

func (fset *HistogramSet) Add(setName, binName string, count int) {
	fstats, ok := fset.HistogramMap[setName]
	if !ok {
		fstats = NewHistogram(setName)
	}
	fstats.Add(binName, count)
	fset.HistogramMap[setName] = fstats
}

/*
func (fset *HistogramSet) AddString(frequencyName, itemName string) {
	fstats, ok := fset.HistogramMap[frequencyName]
	if !ok {
		fstats = NewHistogram(frequencyName)
	}
	fstats.Add(itemName, 1)
	fset.HistogramMap[frequencyName] = fstats
}*/

func (fset *HistogramSet) Names() []string {
	names := []string{}
	for name := range fset.HistogramMap {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (fset *HistogramSet) TotalCount() uint64 {
	totalCount := uint64(0)
	for _, fstats := range fset.HistogramMap {
		totalCount += fstats.TotalCount()
	}
	return totalCount
}

func (fset *HistogramSet) LeafStats(name string) *Histogram {
	if len(name) == 0 {
		name = "leaf stats"
	}
	setLeafStats := NewHistogram(name)
	for _, fstats := range fset.HistogramMap {
		for k, v := range fstats.Items {
			setLeafStats.Add(k, v)
		}
	}
	return setLeafStats
}

func (fset *HistogramSet) ToDataSeriesDistinct() (statictimeseries.DataSeries, error) {
	ds := statictimeseries.NewDataSeries()
	ds.SeriesName = fset.Name
	for rfc3339, fs := range fset.HistogramMap {
		dt, err := time.Parse(time.RFC3339, rfc3339)
		if err != nil {
			return ds, err
		}
		ds.AddItem(statictimeseries.DataItem{
			SeriesName: fset.Name,
			Time:       dt,
			Value:      int64(len(fs.Items))})
	}
	return ds, nil
}

func (fset *HistogramSet) WriteXLSX(path, colName1, colName2, colNameCount string) error {
	// WriteXLSX writes a table as an Excel XLSX file with
	// row formatter option.
	f := excelize.NewFile()
	// Create a new sheet.

	sheetName := strings.TrimSpace(fset.Name)
	if len(sheetName) == 0 {
		sheetName = "Sheet0"
	}
	index := f.NewSheet(sheetName)

	colName1 = strings.TrimSpace(colName1)
	if len(colName1) == 0 {
		colName1 = fset.Name
	}
	if len(colName1) == 0 {
		colName1 = "Column1"
	}
	colName2 = strings.TrimSpace(colName2)
	if len(colName1) == 0 {
		for _, fstats := range fset.HistogramMap {
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
	for fstatsName, fstats := range fset.HistogramMap {
		fstatsNameDt := time.Now()
		if fset.KeyIsTime {
			fstatsNameDt, err = time.Parse(time.RFC3339, fstatsName)
			if err != nil {
				return err
			}
		}
		for itemName, itemCount := range fstats.Items {
			var rowVals []interface{}
			if fset.KeyIsTime {
				rowVals = []interface{}{fstatsNameDt, itemName, itemCount}
			} else {
				rowVals = []interface{}{fstatsName, itemName, itemCount}
			}
			excelizeutil.SetRowValues(f, sheetName, rowIdx, rowVals)
			rowIdx++
		}
	}
	f.SetActiveSheet(index)
	// Delete Original Sheet
	f.DeleteSheet(f.GetSheetName(0))
	// Save xlsx file by the given path.
	return f.SaveAs(path)

}

// HistogramSetDatetimeToQuarter converts a HistogramSet
// by date to one by quarter.s.
func HistogramSetDatetimeToQuarter(name string, fsetIn *HistogramSet) (*HistogramSet, error) {
	fsetQtr := NewHistogramSet(name)
	for rfc3339, fstats := range fsetIn.HistogramMap {
		dt, err := time.Parse(time.RFC3339, rfc3339)
		if err != nil {
			return fsetQtr, err
		}
		dt = timeutil.QuarterStart(dt)
		rfc3339Qtr := dt.Format(time.RFC3339)
		for item, count := range fstats.Items {
			fsetQtr.Add(rfc3339Qtr, item, count)
		}
	}
	return fsetQtr, nil
}

// HistogramSetTimeKeyCount returns a DataSeries when
// the first key is a RFC3339 time and a sum of items
// is desired per time.
func HistogramSetTimeKeyCount(fset HistogramSet) (statictimeseries.DataSeries, error) {
	ds := statictimeseries.NewDataSeries()
	ds.SeriesName = fset.Name
	for rfc3339, fstats := range fset.HistogramMap {
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

func HistogramSetTimeKeyCountTable(fset HistogramSet, interval timeutil.Interval, countColName string) (table.Table, error) {
	ds, err := HistogramSetTimeKeyCount(fset)
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

func HistogramSetTimeKeyCountWriteXLSX(filename string, fset HistogramSet, interval timeutil.Interval, countColName string) error {
	tbl, err := HistogramSetTimeKeyCountTable(fset, interval, countColName)
	if err != nil {
		return err
	}
	tbl.FormatFunc = table.FormatTimeAndInts
	return table.WriteXLSX(filename, &tbl)
}
