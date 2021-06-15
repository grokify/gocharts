package timeseries

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

// TimeSeriesToTable generates a `table.Table` given a `TimeSeries`.
func TimeSeriesToTable(ds TimeSeries, col2 string, dtFmt func(dt time.Time) string) table.Table {
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

// TimeSeriesWriteXLSX writes an XSLX file given a `TimeSeries`
func TimeSeriesWriteXLSX(filename string, ds TimeSeries, col2 string, dtFmt func(dt time.Time) string) error {
	tbl := TimeSeriesToTable(ds, col2, dtFmt)
	tbl.FormatFunc = table.FormatStringAndInts
	return table.WriteXLSX(filename, &tbl)
}
