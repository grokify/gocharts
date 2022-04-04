package timeseries

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/mogo/io/ioutilmore"
	"github.com/grokify/mogo/time/timeutil"

	"github.com/grokify/gocharts/v2/data/table"
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

// Table generates a `table.Table` given a `TimeSeries`.
func (ts *TimeSeries) Table(tableName, dateColumnName, countColumnName string, dtFmt func(dt time.Time) string) table.Table {
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
	tbl := table.NewTable(tableName)
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

func (ts *TimeSeries) TableMonthXOX(timeFmtColName, seriesName, valuesName, yoyName, qoqName, momName string) table.Table {
	if len(strings.TrimSpace(seriesName)) == 0 {
		seriesName = "Series"
	}
	if len(strings.TrimSpace(valuesName)) == 0 {
		valuesName = "Values"
	}
	if len(strings.TrimSpace(yoyName)) == 0 {
		yoyName = "YoY"
	}
	if len(strings.TrimSpace(qoqName)) == 0 {
		qoqName = "QoQ"
	}
	if len(strings.TrimSpace(momName)) == 0 {
		momName = "MoM"
	}
	tsm := ts.ToMonth(true)
	tbl := table.NewTable("")
	cols := []string{seriesName}
	times := tsm.TimeSlice(true)
	for _, dt := range times {
		cols = append(cols, dt.Format(timeFmtColName))
	}
	tbl.Columns = cols
	tbl.FormatMap = map[int]string{
		-1: table.FormatFloat,
		0:  table.FormatString}

	yoy := tsm.TimeSeriesMonthYOY()
	qoq := tsm.TimeSeriesMonthQOQ()
	mom := tsm.TimeSeriesMonthMOM()

	valData := []string{valuesName}
	yoyData := []string{yoyName}
	qoqData := []string{qoqName}
	momData := []string{momName}
	for _, dt := range times {
		tiVal, err := tsm.Get(dt)
		if err != nil {
			panic("internal time not found")
		}
		valData = append(valData, strconv.FormatFloat(tiVal.Float64(), 'f', -1, 64))
		tiYOY, err := yoy.Get(dt)
		if err != nil {
			yoyData = append(yoyData, "0")
		} else {
			yoyData = append(yoyData, strconv.FormatFloat(tiYOY.Float64(), 'f', -1, 64))
		}
		tiQOQ, err := qoq.Get(dt)
		if err != nil {
			qoqData = append(qoqData, "0")
		} else {
			qoqData = append(qoqData, strconv.FormatFloat(tiQOQ.Float64(), 'f', -1, 64))
		}
		tiMOM, err := mom.Get(dt)
		if err != nil {
			momData = append(momData, "0")
		} else {
			momData = append(momData, strconv.FormatFloat(tiMOM.Float64(), 'f', -1, 64))
		}
	}
	tbl.Rows = [][]string{valData, yoyData, qoqData, momData}
	return tbl
}

// WriteJSON writes the data to a JSON file. To write a minimized JSON
// file use an empty string for `prefix` and `indent`.
func (ts *TimeSeries) WriteJSON(filename string, perm os.FileMode, prefix, indent string) error {
	return ioutilmore.WriteFileJSON(filename, ts, perm, prefix, indent)
}

// WriteXLSX writes an XSLX file given a `TimeSeries`
func (ts *TimeSeries) WriteXLSX(filename string, sheetName, dateColumnName, countColumnName string) error {
	tbl := ts.Table(sheetName, dateColumnName, countColumnName, nil)
	return table.WriteXLSX(filename, &tbl)
}
