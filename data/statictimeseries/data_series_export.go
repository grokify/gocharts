package statictimeseries

import (
	"strconv"
	"time"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gotilla/time/timeutil"
)

// DataSeriesToTable generates a `table.TableData` given
// a `DataSeries`.
func DataSeriesToTable(ds DataSeries, interval timeutil.Interval, col2 string) table.TableData {
	tbl := table.NewTableData()
	colDt := "Date"
	dtFmt := func(dt time.Time) string {
		return dt.Format(time.RFC3339)
	}
	if interval == timeutil.Month {
		colDt = "Month"
		dtFmt = func(dt time.Time) string {
			return dt.Format("Jan '06")
		}
	} else if interval == timeutil.Quarter {
		colDt = "Quarter"
		dtFmt = func(dt time.Time) string {
			return timeutil.FormatQuarterYYQ(dt)
		}
	}
	tbl.Columns = []string{colDt, col2}
	itemsSorted := ds.ItemsSorted()
	for _, item := range itemsSorted {
		row := []string{
			dtFmt(item.Time),
			strconv.Itoa(int(item.Value)),
		}
		tbl.Records = append(tbl.Records, row)
	}
	return tbl
}

// DataSeriesWriteXLSX writes an XSLX file given a
// `DataSeries`
func DataSeriesWriteXLSX(filename string, ds DataSeries, interval timeutil.Interval, col2 string) error {
	tbl := DataSeriesToTable(ds, interval, col2)
	tf := &table.TableFormatter{
		Table:     &tbl,
		Formatter: table.FormatStringAndInts}
	return table.WriteXLSXFormatted(filename, tf)
}
