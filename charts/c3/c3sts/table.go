// c3sts is a package for creating C3 charts from statictimeseries data.
// see c3/examples/funnel_chart for usage.
package c3sts

import (
	//"math"
	"strconv"

	"github.com/grokify/mogo/math/mathutil"
	"github.com/grokify/mogo/strconv/strconvutil"

	"github.com/grokify/gocharts/v2/charts/c3"
	"github.com/grokify/gocharts/v2/data/timeseries"
)

// const StyleSimple = "border:1px solid #000;border-collapse:collapse"

// TableData is used to hold generic, simple table data to be generated
// by a template using `SimpleTable`. Use with `DataRowsToTableRows` to
// convert output from `timeseries.Report` and
// `timeseries.ReportAxisX`. This is used in the C3 Bar Chart
// example.
/*
type TableData struct {
	Id    string
	Style string // border:1px solid #000;border-collapse:collapse
	Rows  [][]string
}*/

// DataRowsToTableRows Builds rows from the output of timeseries.Report,
// array of []timeseries.RowInt64.
func DataRowsToTableRows(rep []timeseries.RowInt64, axis []string, addQoQPct, addFunnelPct bool, countLabel, qoqLabel, funnelLabel string) ([][]string, []timeseries.RowFloat64, []timeseries.RowFloat64) {
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
	qoq := []timeseries.RowFloat64{}
	fun := []timeseries.RowFloat64{}
	if addQoQPct {
		qoq = timeseries.ReportGrowthPct(rep)
		rows = append(rows, axis)
		rows[len(rows)-1] = unshift(rows[len(rows)-1], qoqLabel)
		for _, r := range qoq {
			rows = append(rows,
				r.Flatten(strconvutil.FormatFloat64ToIntString, 1, "0%"))
		}
	}
	if addFunnelPct {
		fun = timeseries.ReportFunnelPct(rep)
		rows = append(rows, axis)
		rows[len(rows)-1] = unshift(rows[len(rows)-1], funnelLabel)
		for _, r := range fun {
			rows = append(rows,
				r.Flatten(strconvutil.FormatFloat64ToIntStringFunnel, 0, ""))
		}
	}
	return rows, qoq, fun
}

func unshift(a []string, x string) []string { return append([]string{x}, a...) }

func QoqDataToChart(domID string, axis c3.C3Axis, qoqData []timeseries.RowFloat64) c3.C3Chart {
	qoqChart := c3.C3Chart{
		Bindto: "#" + domID,
		Data: c3.C3ChartData{
			Columns: [][]interface{}{},
		},
		Axis: axis,
		Grid: c3.C3Grid{Y: c3.C3GridLines{Show: true}}}

	for _, r := range qoqData {
		r2 := []interface{}{}
		r2 = append(r2, r.Name)
		for _, v := range r.Values {
			r2 = append(r2, int(mathutil.Round(strconvutil.ChangeToXoXPct(v))))
		}
		qoqChart.Data.Columns = append(qoqChart.Data.Columns, r2)
	}
	return qoqChart
}

func FunnelDataToChart(domID string, axis c3.C3Axis, funnelData []timeseries.RowFloat64) c3.C3Chart {
	funnelChart := c3.C3Chart{
		Bindto: "#" + domID,
		Data: c3.C3ChartData{
			Columns: [][]interface{}{},
		},
		Axis: axis,
		Grid: c3.C3Grid{Y: c3.C3GridLines{Show: true}}}

	for _, r := range funnelData {
		r2 := []interface{}{}
		r2 = append(r2, r.Name)
		for _, v := range r.Values {
			r2 = append(r2, int(mathutil.Round(strconvutil.ChangeToFunnelPct(v))))
		}
		funnelChart.Data.Columns = append(funnelChart.Data.Columns, r2)
	}
	return funnelChart
}
