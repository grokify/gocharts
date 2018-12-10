package tables

import (
	"strconv"

	"github.com/grokify/gocharts/data/statictimeseries"
	scu "github.com/grokify/gotilla/strconv/strconvutil"
)

const StyleSimple = "border:1px solid #000;border-collapse:collapse"

type TableData struct {
	Id    string
	Style string // border:1px solid #000;border-collapse:collapse
	Rows  [][]string
}

// DataRowsToTableRows Builds rows from the output of statictimeseries.Report,
// array of []statictimeseries.RowInt64
func DataRowsToTableRows(rep []statictimeseries.RowInt64, axis []string, addQoQPct, addFunnelPct bool, countLabel, qoqLabel, funnelLabel string) [][]string {
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
	if addQoQPct {
		qoq := statictimeseries.ReportGrowthPct(rep)
		rows = append(rows, axis)
		rows[len(rows)-1] = unshift(rows[len(rows)-1], qoqLabel)
		for _, r := range qoq {
			rows = append(rows, r.Flatten(
				scu.FormatFloat64ToIntString,
			))
		}
	}
	if addFunnelPct {
		pcts := statictimeseries.ReportFunnelPct(rep)
		rows = append(rows, axis)
		rows[len(rows)-1] = unshift(rows[len(rows)-1], funnelLabel)
		for _, r := range pcts {
			rows = append(rows, r.Flatten(
				scu.FormatFloat64ToIntStringFunnel,
			))
		}
	}
	return rows
}

func unshift(a []string, x string) []string { return append([]string{x}, a...) }
