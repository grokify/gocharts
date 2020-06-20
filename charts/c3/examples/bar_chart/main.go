package main

import (
	// Core Demo
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/gotilla/time/timeutil"

	"github.com/grokify/gocharts/charts/c3"
	"github.com/grokify/gocharts/data/statictimeseries"

	// More Info
	"github.com/grokify/gocharts/tables"
)

const repoLink = "https://github.com/grokify/gocharts/tree/master/charts/c3/examples/bar_chart"

func getData(numQuarters int) []statictimeseries.DataItem {
	dataItems := []statictimeseries.DataItem{}

	quarterStart := timeutil.PrevQuarters(time.Now(), uint(numQuarters))

	for i := 1; i <= 3; i++ {
		for j := 0; j < numQuarters; j++ {
			dataItems = append(dataItems, statictimeseries.DataItem{
				SeriesName: fmt.Sprintf("Data Series %d", i),
				Time:       timeutil.NextQuarters(quarterStart, uint(j)),
				Value:      int64(i + j)})
		}
	}

	return dataItems
}

func getDataSeriesSetSimple(numQuarters int) statictimeseries.DataSeriesSet {
	dataItems := getData(numQuarters)

	ds3 := statictimeseries.NewDataSeriesSetSimple()

	// Add statictimeseries.DataItem slice in 1 function call
	ds3.AddItems(dataItems...)

	// Add individual statictimeseries.DataItem items
	for _, di := range dataItems {
		ds3.AddItem(di)
	}
	ds3.Inflate()
	return ds3
}

func buildBarChart(ds3 statictimeseries.DataSeriesSet, numCols int, lowFirst bool) (c3.C3Chart, []statictimeseries.RowInt64) {
	rep := statictimeseries.Report(ds3, numCols, lowFirst)
	fmtutil.PrintJSON(rep)
	axis := statictimeseries.ReportAxisX(ds3, numCols,
		func(t time.Time) string { return timeutil.FormatQuarterYYYYQ(t) })

	chart := c3.DataSeriesSetSimpleToC3ChartBar(rep, c3.C3Bar{})
	chart.Axis = c3.C3Axis{X: c3.C3AxisX{Type: "category", Categories: axis}}
	return chart, rep
}

func buildMoreInfoHTML(ds3 statictimeseries.DataSeriesSet, c3Bar c3.C3Chart, rep []statictimeseries.RowInt64) string {
	moreInfoHTML := ""

	axis := c3Bar.Axis.X.Categories

	tableRows, qoqData, funnelData := tables.DataRowsToTableRows(rep, axis, true, true, "Count", "QoQ", "Funnel")

	addQoqChart := true
	addFunnelChart := true
	addStatsTable := true

	if addQoqChart {
		domId := "qoqChart"
		qoqChart := tables.QoqDataToChart(domId, c3Bar.Axis, qoqData)

		moreInfoHTML += "<h2>QoQ Chart</h2>" + c3.C3ChartHtmlSimple(domId, qoqChart)
	}

	if addFunnelChart {
		domId := "funnelChart"
		funChart := tables.FunnelDataToChart(domId, c3Bar.Axis, funnelData)

		moreInfoHTML += "<h2>Funnel Chart</h2>" + c3.C3ChartHtmlSimple(domId, funChart)
	}

	if addStatsTable {
		table := tables.TableData{
			Id:    "funnelpct",
			Style: tables.StyleSimple,
			Rows:  tableRows}
		moreInfoHTML += "<h2>Stats</h2>" + tables.SimpleTable(table)
	}

	return moreInfoHTML
}

func main() {
	numQuarters := 5

	ds3 := getDataSeriesSetSimple(numQuarters)

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
