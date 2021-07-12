package timeseries

import (
	"strconv"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/simplego/type/stringsutil"
)

func TimeSeriesSliceTimes(dsSlice []TimeSeries) []string {
	timeStrings := []string{}
	for _, ds := range dsSlice {
		keys := ds.Keys()
		timeStrings = append(timeStrings, keys...)
	}
	return stringsutil.SliceCondenseSpace(timeStrings, true, true)
}

func TimeSeriesSliceNames(dsSlice []TimeSeries) []string {
	names := []string{}
	for _, ds := range dsSlice {
		names = append(names, ds.SeriesName)
	}
	return stringsutil.SliceCondenseSpace(names, false, false)
}

func TimeSeriesSliceTable(dsSlice []TimeSeries) table.Table {
	tbl := table.NewTable("")
	names := TimeSeriesSliceNames(dsSlice)
	tbl.Columns = []string{"Date"}
	tbl.Columns = append(tbl.Columns, names...)
	timeStrings := TimeSeriesSliceTimes(dsSlice)
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
		tbl.Rows = append(tbl.Rows, row)
	}
	return tbl
}

// TimeSeries writes a slice of TimeSeries to an
// Excel XLSX file for easy consumption.
func TimeSeriesSliceWriteXLSX(filename string, dsSlice []TimeSeries) error {
	tbl := TimeSeriesSliceTable(dsSlice)
	tbl.FormatFunc = table.FormatTimeAndFloats
	return table.WriteXLSX(filename, &tbl)
}
