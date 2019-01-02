package tables

import (
	"strconv"

	sts "github.com/grokify/gocharts/data/statictimeseries"
	scu "github.com/grokify/gotilla/strconv/strconvutil"
)

const StyleSimple = "border:1px solid #000;border-collapse:collapse"

// TableData is used to hold generic, simple table data to be generated
// by a template using `SimpleTable`. Use with `DataRowsToTableRows` to
// convert output from `statictimeseries.Report` and
// `statictimeseries.ReportAxisX`.
type TableData struct {
	Id    string
	Style string // border:1px solid #000;border-collapse:collapse
	Rows  [][]string
}

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
			rows = append(rows, r.Flatten(
				scu.FormatFloat64ToIntString,
			))
		}
	}
	if addFunnelPct {
		fun = sts.ReportFunnelPct(rep)
		rows = append(rows, axis)
		rows[len(rows)-1] = unshift(rows[len(rows)-1], funnelLabel)
		for _, r := range fun {
			rows = append(rows, r.Flatten(
				scu.FormatFloat64ToIntStringFunnel,
			))
		}
	}
	return rows, qoq, fun
}

func unshift(a []string, x string) []string { return append([]string{x}, a...) }
