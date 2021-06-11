package statictimeseries

import (
	"strconv"
	"time"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/simplego/time/timeutil"
)

func TimeFormatRFC3339(dt time.Time) string {
	return dt.Format(time.RFC3339)
}

func TimeFormatNiceMonth(dt time.Time) string {
	return dt.Format("Jan '06")
}

func TimeFormatNiceQuarter(dt time.Time) string {
	return timeutil.FormatQuarterYYQ(dt)
}

// DataSeriesToTable generates a `table.Table` given a `DataSeries`.
func DataSeriesToTable(ds DataSeries, col2 string, dtFmt func(dt time.Time) string) table.Table {
	tbl := table.NewTable()
	colDt := "Date"
	/*dtFmt := func(dt time.Time) string {
		return dt.Format(time.RFC3339)
	}*/
	if ds.Interval == timeutil.Month {
		colDt = "Month"
		/*dtFmt = func(dt time.Time) string {
			return dt.Format("Jan '06")
		}*/
	} else if ds.Interval == timeutil.Quarter {
		colDt = "Quarter"
		/*dtFmt = func(dt time.Time) string {
			return timeutil.FormatQuarterYYQ(dt)
		}*/
	}
	tbl.Columns = []string{colDt, col2}
	itemsSorted := ds.ItemsSorted()
	for _, item := range itemsSorted {
		row := []string{
			dtFmt(item.Time),
			strconv.Itoa(int(item.Value)),
		}
		tbl.Rows = append(tbl.Rows, row)
	}
	return tbl
}

// DataSeriesWriteXLSX writes an XSLX file given a `DataSeries`
func DataSeriesWriteXLSX(filename string, ds DataSeries, col2 string, dtFmt func(dt time.Time) string) error {
	tbl := DataSeriesToTable(ds, col2, dtFmt)
	tbl.FormatFunc = table.FormatStringAndInts
	return table.WriteXLSX(filename, &tbl)
}
