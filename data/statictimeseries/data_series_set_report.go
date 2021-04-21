package statictimeseries

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/simplego/time/timeutil"
)

// ReportAxisX generates data for use with `C3Chart.C3Axis.C3AxisX.Categories`.
func ReportAxisX(dss DataSeriesSet, cols int, conv func(time.Time) string) []string {
	var times timeutil.TimeSlice
	if cols < len(dss.Times) {
		min := len(dss.Times) - cols
		times = dss.Times[min:]
	} else { // cols >= len(dss.Times)
		times = dss.Times
	}
	cats := []string{}
	for _, t := range times {
		cats = append(cats, conv(t))
	}
	return cats
}

// Report generates data for use with `C3Chart.C3ChartData.Columns`.
func Report(dss DataSeriesSet, cols int, lowFirst bool) []RowInt64 {
	rows := []RowInt64{}
	var times timeutil.TimeSlice
	var timePlus1 time.Time
	havePlus1 := false
	if cols < len(dss.Times) {
		min := len(dss.Times) - cols
		prev := min - 1
		times = dss.Times[min:]
		timePlus1 = dss.Times[prev]
		havePlus1 = true
	} else { // cols >= len(dss.Times)
		times = dss.Times
		if cols > len(dss.Times) {
			timePlus1 = dss.Times[len(dss.Times)-cols-1]
			havePlus1 = true
		}
	}
	timePlus1Rfc := timePlus1.UTC().Format(time.RFC3339)
	if !lowFirst {
		times = sort.Reverse(times).(timeutil.TimeSlice)
	}
	for _, seriesName := range dss.Order {
		row := RowInt64{
			Name:        seriesName + " Count",
			HavePlusOne: havePlus1,
		}
		if ds, ok := dss.Series[seriesName]; !ok {
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

type DssTableOpts struct {
	FuncFormatTime func(time.Time) string
	TotalInclude   bool
	TotalTitle     string
	PercentInclude bool
	PercentSuffix  string
}

func (opts *DssTableOpts) TotalTitleOrDefault() string {
	if len(opts.TotalTitle) > 0 {
		return opts.TotalTitle
	}
	return "Total"
}

func (opts *DssTableOpts) PercentSuffixOrDefault() string {
	if len(opts.PercentSuffix) > 0 {
		return opts.PercentSuffix
	}
	return "%"
}

// ToTable returns a `DataSeriesSet` as a `table.TableData`.
// func DssToTable(dss DataSeriesSet, fmtTime func(time.Time) string) (table.Table, error) {
func (dss *DataSeriesSet) ToTable(opts *DssTableOpts) (table.Table, error) {
	if opts == nil {
		opts = &DssTableOpts{}
	}
	tbl := table.NewTable()
	seriesNames := dss.SeriesNames()
	tbl.Columns = []string{"Time"}
	tbl.Columns = append(tbl.Columns, seriesNames...)
	if opts.TotalInclude {
		tbl.Columns = append(tbl.Columns, opts.TotalTitleOrDefault())
	}
	if opts.PercentInclude {
		for _, seriesName := range seriesNames {
			tbl.Columns = append(
				tbl.Columns,
				seriesName+" "+opts.PercentSuffixOrDefault())
		}
	}
	timeStrings := dss.TimeStrings()
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
			item, err := dss.GetItem(seriesName, rfc3339)
			if err == nil {
				if item.IsFloat {
					line = append(line, fmt.Sprintf("%.10f", item.ValueFloat))
				} else {
					line = append(line, strconv.Itoa(int(item.Value)))
				}
				lineTotal += item.ValueFloat64()
				seriesValues = append(seriesValues, item.ValueFloat64())
			} else {
				line = append(line, "0")
				seriesValues = append(seriesValues, 0)
			}
		}
		if opts.TotalInclude {
			line = append(line, fmt.Sprintf("%.10f", lineTotal))
		}
		if opts.PercentInclude {
			for _, seriesValue := range seriesValues {
				line = append(line, fmt.Sprintf("%.10f", seriesValue/lineTotal))
			}
		}
		tbl.Records = append(tbl.Records, line)
	}
	return tbl, nil
}

func (dss *DataSeriesSet) WriteXLSX(filename string, opts *DssTableOpts) error {
	tbl, err := dss.ToTable(opts)
	if err != nil {
		return err
	}
	if dss.Interval == timeutil.Month {
		tbl.FormatFunc = table.FormatMonthAndFloats
	} else {
		tbl.FormatFunc = table.FormatDateAndFloats
	}
	/*tbl.FormatMap = map[int]string{
	0:  "string",
	-1: "float"}*/
	return table.WriteXLSX(filename, &tbl)
}
