// statictimeseriesdata provides tools for adding and formatting
// static time series data for reporting purposes.
package statictimeseries

import (
	"strconv"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gotilla/type/stringsutil"
)

func DataSeriesSliceTimes(dsSlice []DataSeries) []string {
	timeStrings := []string{}
	for _, ds := range dsSlice {
		keys := ds.Keys()
		timeStrings = append(timeStrings, keys...)
	}
	return stringsutil.SliceCondenseSpace(timeStrings, true, true)
}

func DataSeriesSliceNames(dsSlice []DataSeries) []string {
	names := []string{}
	for _, ds := range dsSlice {
		names = append(names, ds.SeriesName)
	}
	return stringsutil.SliceCondenseSpace(names, false, false)
}

func DataSeriesSliceTable(dsSlice []DataSeries) table.TableData {
	tbl := table.NewTableData()
	names := DataSeriesSliceNames(dsSlice)
	tbl.Columns = []string{"Date"}
	tbl.Columns = append(tbl.Columns, names...)
	timeStrings := DataSeriesSliceTimes(dsSlice)
	for _, timeStr := range timeStrings {
		row := []string{timeStr}
		for _, ds := range dsSlice {
			if item, ok := ds.ItemMap[timeStr]; ok {
				if item.IsFloat {
					row = append(row, strconv.FormatFloat(item.ValueFloat, 'f', -1, 64))
				} else {
					row = append(row, strconv.Itoa(int(item.Value)))
				}
			} else {
				row = append(row, "0")
			}
		}
		tbl.Records = append(tbl.Records, row)
	}
	return tbl
}

// DataSeriesSliceWriteXLSX writes a slice of DataSeries to an
// Excel XLSX file for easy consumption.
func DataSeriesSliceWriteXLSX(filename string, dsSlice []DataSeries) error {
	tbl := DataSeriesSliceTable(dsSlice)
	tbl.FormatFunc = table.FormatTimeAndFloats
	return table.WriteXLSX(filename, &tbl)
}
