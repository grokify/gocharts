package timeseries

import (
	"strconv"
	"strings"
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

// ToTable generates a `table.Table` given a `TimeSeries`.
func (ts *TimeSeries) ToTable(tableName, dateColumnName, countColumnName string, dtFmt func(dt time.Time) string) table.Table {
	// previously only took dateColumnName as a parameter.
	if len(strings.TrimSpace(dateColumnName)) == 0 {
		switch ts.Interval {
		case timeutil.Month:
			dateColumnName = "Month"
		case timeutil.Quarter:
			dateColumnName = "Quarter"
		default:
			dateColumnName = "Date"
		}
	}
	if len(strings.TrimSpace(countColumnName)) == 0 {
		countColumnName = "Count"
	}
	tbl := table.NewTable()
	tbl.Name = tableName
	tbl.Columns = []string{dateColumnName, countColumnName}
	tbl.FormatMap = map[int]string{}
	if ts.IsFloat {
		tbl.FormatMap[1] = table.FormatFloat
	} else {
		tbl.FormatMap[1] = table.FormatInt
	}
	if dtFmt == nil {
		dtFmt = func(dt time.Time) string {
			return dt.Format(time.RFC3339)
		}
		tbl.FormatMap[0] = table.FormatTime
	}
	itemsSorted := ts.ItemsSorted()
	for _, item := range itemsSorted {
		row := []string{
			dtFmt(item.Time)}
		if ts.IsFloat {
			row = append(row, strconv.FormatFloat(item.ValueFloat, 'f', -1, 64))
		} else {
			row = append(row, strconv.Itoa(int(item.Value)))
		}
		tbl.Rows = append(tbl.Rows, row)
	}
	return tbl
}

// WriteXLSX writes an XSLX file given a `TimeSeries`
func (ts *TimeSeries) WriteXLSX(filename string, sheetName, dateColumnName, countColumnName string) error {
	tbl := ts.ToTable(sheetName, dateColumnName, countColumnName, nil)
	return table.WriteXLSX(filename, &tbl)
}
