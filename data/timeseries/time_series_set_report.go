package timeseries

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/mogo/os/osutil"
	"github.com/grokify/mogo/time/timeslice"
	"github.com/grokify/mogo/time/timeutil"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gocharts/v2/data/table/format"
)

// ReportAxisX generates data for use with `C3Chart.C3Axis.C3AxisX.Categories`.
func ReportAxisX(set TimeSeriesSet, cols int, conv func(time.Time) string) []string {
	var times timeslice.TimeSlice
	if cols < len(set.Times) {
		min := len(set.Times) - cols
		times = set.Times[min:]
	} else { // cols >= len(set.Times)
		times = set.Times
	}
	cats := []string{}
	for _, t := range times {
		cats = append(cats, conv(t))
	}
	return cats
}

// Report generates data for use with `C3Chart.C3ChartData.Columns`.
func Report(set TimeSeriesSet, cols int, lowFirst bool) []RowInt64 {
	rows := []RowInt64{}
	var times timeslice.TimeSlice
	var timePlus1 time.Time
	havePlus1 := false
	if cols < len(set.Times) {
		min := len(set.Times) - cols
		prev := min - 1
		times = set.Times[min:]
		timePlus1 = set.Times[prev]
		havePlus1 = true
	} else { // cols >= len(set.Times)
		times = set.Times
		if cols > len(set.Times) {
			timePlus1 = set.Times[len(set.Times)-cols-1]
			havePlus1 = true
		}
	}
	timePlus1Rfc := timePlus1.UTC().Format(time.RFC3339)
	if !lowFirst {
		times = sort.Reverse(times).(timeslice.TimeSlice)
	}
	for _, seriesName := range set.Order {
		row := RowInt64{
			Name:        seriesName + " Count",
			HavePlusOne: havePlus1,
		}
		if ds, ok := set.Series[seriesName]; !ok {
			for i := 0; i < cols; i++ {
				row.Values = append(row.Values, 0)
			}
			if havePlus1 {
				row.ValuePlusOne = 0
			}
		} else {
			for _, t := range times {
				rfc := t.UTC().Format(time.RFC3339)
				if item, ok := ds.ItemMap[rfc]; ok {
					row.Values = append(row.Values, item.Value)
				} else {
					row.Values = append(row.Values, 0)
				}
			}
			if havePlus1 {
				if item, ok := ds.ItemMap[timePlus1Rfc]; ok {
					row.ValuePlusOne = item.Value
				} else {
					row.ValuePlusOne = 0
				}
			}
		}
		rows = append(rows, row)
	}
	return rows
}

func ReportFunnelPct(rows []RowInt64) []RowFloat64 {
	pcts := []RowFloat64{}
	if len(rows) < 2 {
		return pcts
	}
	for i := 0; i < len(rows)-1; i++ {
		r := RowFloat64{Name: fmt.Sprintf("Success Pct #%v", i)}
		j := i + 1
		for k := 0; k < len(rows[0].Values); k++ {
			v1 := rows[i].Values[k]
			v2 := rows[j].Values[k]
			pct := float64(v2) / float64(v1)
			r.Values = append(r.Values, pct)
		}
		pcts = append(pcts, r)
	}
	return pcts
}

func ReportGrowthPct(rows []RowInt64) []RowFloat64 {
	grows := []RowFloat64{}
	if len(rows) == 0 {
		return grows
	}
	for i := 0; i < len(rows); i++ {
		r := rows[i]
		grow := RowFloat64{Name: fmt.Sprintf("%v XoX", r.Name)}
		if r.HavePlusOne {
			pct := float64(r.Values[0]) / float64(r.ValuePlusOne)
			grow.Values = append(grow.Values, pct)
		}
		for j := 0; j < len(r.Values)-1; j++ {
			k := j + 1
			pct := float64(r.Values[k]) / float64(r.Values[j])
			grow.Values = append(grow.Values, pct)
		}
		grows = append(grows, grow)
	}
	return grows
}

type TimeSeriesSetTableOpts struct {
	TimeColumnTitle string
	FuncFormatTime  func(time.Time) string
	TotalInclude    bool
	TotalTitle      string
	PercentInclude  bool
	PercentSuffix   string
}

func (opts *TimeSeriesSetTableOpts) TotalTitleOrDefault() string {
	if len(opts.TotalTitle) > 0 {
		return opts.TotalTitle
	}
	return "Total"
}

func (opts *TimeSeriesSetTableOpts) PercentSuffixOrDefault() string {
	if len(opts.PercentSuffix) > 0 {
		return opts.PercentSuffix
	}
	return "%"
}

// Table returns a `table.Table`.
func (set *TimeSeriesSet) Table(opts *TimeSeriesSetTableOpts) (table.Table, error) {
	if opts == nil {
		opts = &TimeSeriesSetTableOpts{}
	}
	tbl := table.NewTable(set.Name)
	seriesNames := set.SeriesNames()
	timeColumnTitle := strings.TrimSpace(opts.TimeColumnTitle)
	if len(timeColumnTitle) == 0 {
		timeColumnTitle = "Time"
	}
	tbl.Columns = []string{timeColumnTitle}
	tbl.Columns = append(tbl.Columns, seriesNames...)
	tbl.FormatMap = map[int]string{0: table.FormatTime}
	for i := range seriesNames {
		if set.IsFloat {
			tbl.FormatMap[i+1] = table.FormatFloat
		} else {
			tbl.FormatMap[i+1] = table.FormatInt
		}
	}
	if opts.TotalInclude {
		tbl.Columns = append(tbl.Columns, opts.TotalTitleOrDefault())
		if set.IsFloat {
			tbl.FormatMap[len(tbl.Columns)-1] = table.FormatFloat
		} else {
			tbl.FormatMap[len(tbl.Columns)-1] = table.FormatInt
		}
	}
	if opts.PercentInclude {
		for _, seriesName := range seriesNames {
			tbl.Columns = append(
				tbl.Columns,
				seriesName+" "+opts.PercentSuffixOrDefault())
			tbl.FormatMap[len(tbl.Columns)-1] = table.FormatFloat
		}
	}
	timeStrings := set.TimeStrings()
	for _, rfc3339 := range timeStrings {
		line := []string{}
		if opts.FuncFormatTime == nil {
			line = append(line, rfc3339)
		} else {
			dt, err := time.Parse(time.RFC3339, rfc3339)
			if err != nil {
				return tbl, err
			}
			line = append(line, opts.FuncFormatTime(dt))
		}
		lineTotal := float64(0)
		seriesValues := []float64{}
		for _, seriesName := range seriesNames {
			item, err := set.Item(seriesName, rfc3339)
			if err != nil {
				line = append(line, "0")
				seriesValues = append(seriesValues, 0)
			} else {
				if item.IsFloat {
					line = append(line, fmt.Sprintf("%.10f", item.ValueFloat))
				} else {
					line = append(line, strconv.Itoa(int(item.Value)))
				}
				lineTotal += item.Float64()
				seriesValues = append(seriesValues, item.Float64())
			}
		}
		if opts.TotalInclude {
			if set.IsFloat {
				line = append(line, fmt.Sprintf("%.10f", lineTotal))
			} else {
				line = append(line, strconv.Itoa(int(lineTotal)))
			}
		}
		if opts.PercentInclude {
			for _, seriesValue := range seriesValues {
				if lineTotal == 0 {
					line = append(line, fmt.Sprintf("%.10f", lineTotal))
				} else {
					line = append(line, fmt.Sprintf("%.10f", seriesValue/lineTotal))
				}
			}
		}
		tbl.Rows = append(tbl.Rows, line)
	}
	return tbl, nil
}

func (set *TimeSeriesSet) TableYearYOY(seriesColName, valuesSuffix, yoySuffix string) table.Table {
	if set.Interval != timeutil.Year {
		panic("interval is not year")
	}
	if len(strings.TrimSpace(seriesColName)) == 0 {
		seriesColName = "Series"
	}
	if len(strings.TrimSpace(valuesSuffix)) == 0 {
		valuesSuffix = "Values"
	}
	if len(strings.TrimSpace(yoySuffix)) == 0 {
		yoySuffix = "YoY"
	}

	tbl := table.NewTable(set.Name)
	cols := []string{seriesColName}
	times := set.TimeSlice(true)
	for _, dt := range times {
		cols = append(cols, dt.Format("2006"))
	}
	tbl.Columns = cols
	tbl.FormatMap = map[int]string{
		-1: table.FormatFloat,
		0:  table.FormatString}

	order := set.Order
	if len(order) != len(set.Series) {
		names := []string{}
		for name := range set.Series {
			names = append(names, name)
		}
		sort.Strings(names)
	}

	for _, seriesName := range order {
		ts, ok := set.Series[seriesName]
		if !ok {
			panic(fmt.Sprintf("set not found [%s]", seriesName))
		}
		tblSeries := ts.TableYearYOY(seriesColName, seriesName+" "+valuesSuffix, seriesName+" "+yoySuffix)
		if len(tbl.Columns) != len(tblSeries.Columns) {
			panic("table length mismatch")
		}
		tbl.Rows = append(tbl.Rows, tblSeries.Rows...)
	}

	return tbl
}

// tblYoy := tsAnn.TableYearYOY("Series", "MAU", "YOY")

// WriteJSON writes the TimeSeriesSet to a JSON file. To write a minimized JSON
// file use an empty string for `prefix` and `indent`.
func (set *TimeSeriesSet) WriteJSON(filename string, perm os.FileMode, prefix, indent string) error {
	return osutil.WriteFileJSON(filename, set, perm, prefix, indent)
}

// WriteXLSX writes the TimeSeriesSet as a XLSX spreadsheet file.
func (set *TimeSeriesSet) WriteXLSX(filename string, opts *TimeSeriesSetTableOpts) error {
	tbl, err := set.Table(opts)
	if err != nil {
		return err
	}
	if set.Interval == timeutil.Month {
		tbl.FormatFunc = format.FormatMonthAndFloats
	} else {
		tbl.FormatFunc = format.FormatDateAndFloats
	}
	return table.WriteXLSX(filename, &tbl)
}
