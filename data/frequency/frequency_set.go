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

type FrequencySet struct {
	Name         string
	FrequencyMap map[string]FrequencyStats
	KeyIsTime    bool
}

func NewFrequencySet(name string) FrequencySet {
	return FrequencySet{
		Name:         name,
		FrequencyMap: map[string]FrequencyStats{}}
}

func (fset *FrequencySet) AddDateUidCount(dt time.Time, uid string, count int) {
	fName := dt.Format(time.RFC3339)
	fset.AddStringMore(fName, uid, count)
	if !fset.KeyIsTime {
		fset.KeyIsTime = true
	}
}

func (fset *FrequencySet) AddStringMore(frequencyName, uid string, count int) {
	fstats, ok := fset.FrequencyMap[frequencyName]
	if !ok {
		fstats = NewFrequencyStats(frequencyName)
	}
	fstats.AddStringMore(uid, count)
	fset.FrequencyMap[frequencyName] = fstats
}

func (fset *FrequencySet) AddString(frequencyName, itemName string) {
	fstats, ok := fset.FrequencyMap[frequencyName]
	if !ok {
		fstats = NewFrequencyStats(frequencyName)
	}
	fstats.AddString(itemName)
	fset.FrequencyMap[frequencyName] = fstats
}

func (fset *FrequencySet) Names() []string {
	names := []string{}
	for name := range fset.FrequencyMap {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (fset *FrequencySet) TotalCount() uint64 {
	totalCount := uint64(0)
	for _, fstats := range fset.FrequencyMap {
		totalCount += fstats.TotalCount()
	}
	return totalCount
}

func (fset *FrequencySet) ToDataSeriesDistinct() (statictimeseries.DataSeries, error) {
	ds := statictimeseries.NewDataSeries()
	ds.SeriesName = fset.Name
	for rfc3339, fs := range fset.FrequencyMap {
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

func (fset *FrequencySet) WriteXLSX(path, colName1, colName2, colNameCount string) error {
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
		for _, fstats := range fset.FrequencyMap {
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
	header := []interface{}{colName1, colName1, colNameCount}

	excelizeutil.SetRowValues(f, sheetName, 0, header)
	var err error
	rowIdx := uint(1)
	for fstatsName, fstats := range fset.FrequencyMap {
		fstatsNameDt := time.Now()
		if fset.KeyIsTime {
			fstatsNameDt, err = time.Parse(time.RFC3339, fstatsName)
			if err != nil {
				return err
			}
		}
		for itemName, itemCount := range fstats.Items {
			rowVals := []interface{}{}
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
