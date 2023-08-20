package timeseries

import (
	"strconv"

	"github.com/grokify/mogo/type/stringsutil"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gocharts/v2/data/table/format"
)

func TimeSeriesSliceTimes(tsSlice []TimeSeries) []string {
	timeStrings := []string{}
	for _, ds := range tsSlice {
		keys := ds.Keys()
		timeStrings = append(timeStrings, keys...)
	}
	return stringsutil.SliceCondenseSpace(timeStrings, true, true)
}

func TimeSeriesSliceNames(tsSlice []TimeSeries) []string {
	names := []string{}
	for _, ds := range tsSlice {
		names = append(names, ds.SeriesName)
	}
	return stringsutil.SliceCondenseSpace(names, false, false)
}

func TimeSeriesSliceTable(tsSlice []TimeSeries) table.Table {
	tbl := table.NewTable("")
	// names := TimeSeriesSliceNames(tsSlice)
	tbl.Columns = []string{"Date"}
	tbl.Columns = append(tbl.Columns, TimeSeriesSliceNames(tsSlice)...)
	timeStrings := TimeSeriesSliceTimes(tsSlice)
	for _, timeStr := range timeStrings {
		row := []string{timeStr}
		for _, ds := range tsSlice {
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
func TimeSeriesSliceWriteXLSX(filename string, tsSlice []TimeSeries) error {
	tbl := TimeSeriesSliceTable(tsSlice)
	tbl.FormatFunc = format.FormatTimeAndFloats
	return table.WriteXLSX(filename, []*table.Table{&tbl})
}
