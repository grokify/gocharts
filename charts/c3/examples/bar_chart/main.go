package main

import (
	// Core Demo
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/time/timeutil"

	"github.com/grokify/gocharts/v2/charts/c3"
	"github.com/grokify/gocharts/v2/charts/c3/c3sts"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gocharts/v2/data/timeseries"
	// More Info
)

const repoLink = "https://github.com/grokify/gocharts/tree/master/charts/c3/examples/bar_chart"

func getData(numQuarters uint) []timeseries.TimeItem {
	timeItems := []timeseries.TimeItem{}

	quarterStart := timeutil.QuarterAdd(time.Now(), -1*int(numQuarters))

	for i := 1; i <= 3; i++ {
		for j := 0; j < int(numQuarters); j++ {
			timeItems = append(timeItems, timeseries.TimeItem{
				SeriesName: fmt.Sprintf("Data Series %d", i),
				Time:       timeutil.QuarterAdd(quarterStart, j),
				Value:      int64(i + j)})
		}
	}

	return timeItems
}

func getTimeSeriesSetSimple(numQuarters uint) timeseries.TimeSeriesSet {
	timeItems := getData(numQuarters)

	tss := timeseries.NewTimeSeriesSet("Bar Chart Data")

	// Add timeseries.DataItem slice in 1 function call
	tss.AddItems(timeItems...)

	// Add individual timeseries.DataItem items
	for _, ti := range timeItems {
		tss.AddItems(ti)
	}
	tss.Inflate()
	return tss
}

func buildBarChart(ds3 timeseries.TimeSeriesSet, numCols int, lowFirst bool) (c3.C3Chart, []timeseries.RowInt64) {
	rep := timeseries.Report(ds3, numCols, lowFirst)
	fmtutil.PrintJSON(rep)
	axis := timeseries.ReportAxisX(ds3, numCols,
		func(t time.Time) string { return timeutil.FormatQuarterYYYYQ(t) })

	chart := c3.TimeSeriesSetSimpleToC3ChartBar(rep, c3.C3Bar{})
	chart.Axis = c3.C3Axis{X: c3.C3AxisX{Type: "category", Categories: axis}}
	return chart, rep
}

func buildMoreInfoHTML(ds3 timeseries.TimeSeriesSet, c3Bar c3.C3Chart, rep []timeseries.RowInt64) string {
	moreInfoHTML := ""

	axis := c3Bar.Axis.X.Categories

	tableRows, qoqData, funnelData := c3sts.DataRowsToTableRows(rep, axis, true, true, "Count", "QoQ", "Funnel")

	addQoqChart := true
	addFunnelChart := true
	addStatsTable := true

	if addQoqChart {
		domId := "qoqChart"
		qoqChart := c3sts.QoqDataToChart(domId, c3Bar.Axis, qoqData)

		moreInfoHTML += "<h2>QoQ Chart</h2>" + c3.C3ChartHTMLSimple(domId, qoqChart)
	}

	if addFunnelChart {
		domId := "funnelChart"
		funChart := c3sts.FunnelDataToChart(domId, c3Bar.Axis, funnelData)

		moreInfoHTML += "<h2>Funnel Chart</h2>" + c3.C3ChartHTMLSimple(domId, funChart)
	}

	if addStatsTable {
		tbl := table.Table{
			ID:    "funnelpct",
			Style: table.StyleSimple,
			Rows:  tableRows}
		moreInfoHTML += "<h2>Stats</h2>" + table.SimpleTable(tbl)
	}

	return moreInfoHTML
}

func main() {
	numQuarters := 5

	ds3 := getTimeSeriesSetSimple(uint(numQuarters))

	chart, rep := buildBarChart(ds3, numQuarters, true)

	moreInfoHTML := buildMoreInfoHTML(ds3, chart, rep)

	tmplData := c3.TemplateData{
		HeaderHTML:             "Bar Chart",
		ReportName:             "Bar Chart",
		ReportLink:             repoLink,
		IncludeDataTable:       false,
		IncludeDataTableTotals: false,
		C3Chart:                chart,
		FooterHTML:             moreInfoHTML}

	filename := "output.html"

	err := ioutil.WriteFile(filename, []byte(c3.C3DonutChartPage(tmplData)), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Wrote: %s\n", filename)
}
