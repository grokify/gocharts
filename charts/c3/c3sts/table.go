// c3sts is a package for creating C3 charts from statictimeseries data.
// see c3/examples/funnel_chart for usage.
package c3sts

import (
	//"math"
	"strconv"

	math "github.com/grokify/gotilla/math/mathutil"

	sts "github.com/grokify/gocharts/data/statictimeseries"
	scu "github.com/grokify/gotilla/strconv/strconvutil"

	"github.com/grokify/gocharts/charts/c3"
)

// const StyleSimple = "border:1px solid #000;border-collapse:collapse"

// TableData is used to hold generic, simple table data to be generated
// by a template using `SimpleTable`. Use with `DataRowsToTableRows` to
// convert output from `statictimeseries.Report` and
// `statictimeseries.ReportAxisX`. This is used in the C3 Bar Chart
// example.
/*
type TableData struct {
	Id    string
	Style string // border:1px solid #000;border-collapse:collapse
	Rows  [][]string
}*/

// DataRowsToTableRows Builds rows from the output of statictimeseries.Report,
// array of []statictimeseries.RowInt64.
func DataRowsToTableRows(rep []sts.RowInt64, axis []string, addQoQPct, addFunnelPct bool, countLabel, qoqLabel, funnelLabel string) ([][]string, []sts.RowFloat64, []sts.RowFloat64) {
	rows := [][]string{}
	rows = append(rows, axis)
	rows[len(rows)-1] = unshift(rows[len(rows)-1], countLabel)
	for _, r := range rep {
		rows = append(rows, r.Flatten(
			func(i int64) string {
				return strconv.Itoa(int(i))
			},
		))
	}
	qoq := []sts.RowFloat64{}
	fun := []sts.RowFloat64{}
	if addQoQPct {
		qoq = sts.ReportGrowthPct(rep)
		rows = append(rows, axis)
		rows[len(rows)-1] = unshift(rows[len(rows)-1], qoqLabel)
		for _, r := range qoq {
			rows = append(rows,
				r.Flatten(scu.FormatFloat64ToIntString, 1, "0%"))
		}
	}
	if addFunnelPct {
		fun = sts.ReportFunnelPct(rep)
		rows = append(rows, axis)
		rows[len(rows)-1] = unshift(rows[len(rows)-1], funnelLabel)
		for _, r := range fun {
			rows = append(rows,
				r.Flatten(scu.FormatFloat64ToIntStringFunnel, 0, ""))
		}
	}
	return rows, qoq, fun
}

func unshift(a []string, x string) []string { return append([]string{x}, a...) }

func QoqDataToChart(domId string, axis c3.C3Axis, qoqData []sts.RowFloat64) c3.C3Chart {
	qoqChart := c3.C3Chart{
		Bindto: "#" + domId,
		Data: c3.C3ChartData{
			Columns: [][]interface{}{},
		},
		Axis: axis,
		Grid: c3.C3Grid{Y: c3.C3GridLines{Show: true}}}

	for _, r := range qoqData {
		r2 := []interface{}{}
		r2 = append(r2, r.Name)
		for _, v := range r.Values {
			r2 = append(r2, int(math.Round(scu.ChangeToXoXPct(v))))
		}
		qoqChart.Data.Columns = append(qoqChart.Data.Columns, r2)
	}
	return qoqChart
}

func FunnelDataToChart(domId string, axis c3.C3Axis, funnelData []sts.RowFloat64) c3.C3Chart {
	funnelChart := c3.C3Chart{
		Bindto: "#" + domId,
		Data: c3.C3ChartData{
			Columns: [][]interface{}{},
		},
		Axis: axis,
		Grid: c3.C3Grid{Y: c3.C3GridLines{Show: true}}}

	for _, r := range funnelData {
		r2 := []interface{}{}
		r2 = append(r2, r.Name)
		for _, v := range r.Values {
			r2 = append(r2, int(math.Round(scu.ChangeToFunnelPct(v))))
		}
		funnelChart.Data.Columns = append(funnelChart.Data.Columns, r2)
	}
	return funnelChart
}
