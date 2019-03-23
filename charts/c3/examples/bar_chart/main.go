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
	"math"

	scu "github.com/grokify/gotilla/strconv/strconvutil"

	"github.com/grokify/gocharts/tables"
)

func getData(numQuarters int) []statictimeseries.DataItem {
	dataItems := []statictimeseries.DataItem{}

	quarterStart := timeutil.PrevQuarters(time.Now(), numQuarters)

	for i := 1; i <= 3; i++ {
		for j := 0; j < numQuarters; j++ {
			dataItems = append(dataItems, statictimeseries.DataItem{
				SeriesName: fmt.Sprintf("Data Series %d", i),
				Time:       timeutil.NextQuarters(quarterStart, j),
				Value:      int64(i + j)})
		}
	}

	return dataItems
}

func getDataSeriesSetSimple(numQuarters int) statictimeseries.DataSeriesSetSimple {
	ds3 := statictimeseries.NewDataSeriesSetSimple()

	dataItems := getData(numQuarters)

	ds3.AddItems(dataItems...)

	for _, di := range dataItems {
		ds3.AddItem(di)
	}
	ds3.Inflate()
	return ds3
}

func buildChart(ds3 statictimeseries.DataSeriesSetSimple, numCols int, lowFirst bool) (c3.C3Chart, []statictimeseries.RowInt64) {
	rep := statictimeseries.Report(ds3, numCols, lowFirst)
	fmtutil.PrintJSON(rep)
	axis := statictimeseries.ReportAxisX(ds3, numCols,
		func(t time.Time) string { return timeutil.FormatQuarterYYYYQ(t) })

	chart := c3.DataSeriesSetSimpleToC3ChartBar(rep, c3.C3Bar{})
	chart.Axis = c3.C3Axis{X: c3.C3AxisX{Type: "category", Categories: axis}}
	return chart, rep
}

func buildMoreInfoHTML(ds3 statictimeseries.DataSeriesSetSimple, c3Bar c3.C3Chart, rep []statictimeseries.RowInt64) string {
	moreInfoHTML := ""

	if 1 == 1 {
		axis := c3Bar.Axis.X.Categories

		tableRows, qoqData, funData := tables.DataRowsToTableRows(rep, axis, true, true, "Count", "QoQ", "Funnel")

		if 1 == 1 {
			domId := "qoqchart"
			qoqChart := c3.C3Chart{
				Bindto: "#" + domId,
				Data: c3.C3ChartData{
					Columns: [][]interface{}{},
				},
				Axis: c3Bar.Axis,
				Grid: c3.C3Grid{Y: c3.C3GridLines{Show: true}}}

			for _, r := range qoqData {
				r2 := []interface{}{}
				r2 = append(r2, r.Name)
				for _, v := range r.Values {
					r2 = append(r2, int(math.Round(scu.ChangeToXoXPct(v))))
				}
				qoqChart.Data.Columns = append(qoqChart.Data.Columns, r2)
			}

			moreInfoHTML += "<h2>QoQ Chart</h2>" + c3.C3ChartHtmlSimple(domId, qoqChart)
		}
		if 1 == 1 {
			domId := "funchart"
			funChart := c3.C3Chart{
				Bindto: "#" + domId,
				Data: c3.C3ChartData{
					Columns: [][]interface{}{},
				},
				Axis: c3Bar.Axis,
				Grid: c3.C3Grid{Y: c3.C3GridLines{Show: true}}}

			for _, r := range funData {
				r2 := []interface{}{}
				r2 = append(r2, r.Name)
				for _, v := range r.Values {
					r2 = append(r2, int(math.Round(scu.ChangeToFunnelPct(v))))
				}
				funChart.Data.Columns = append(funChart.Data.Columns, r2)
			}

			moreInfoHTML += "<h2>Funnel Chart</h2>" + c3.C3ChartHtmlSimple(domId, funChart)
		}

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

	chart, rep := buildChart(ds3, numQuarters, true)

	footerHTML := buildMoreInfoHTML(ds3, chart, rep)

	tmplData := c3.TemplateData{
		HeaderHTML:             "Bar Chart",
		ReportName:             "Bar Chart",
		ReportLink:             "",
		IncludeDataTable:       false,
		IncludeDataTableTotals: false,
		C3Chart:                chart,
		FooterHTML:             footerHTML}

	filename := "output.html"

	err := ioutil.WriteFile(filename, []byte(c3.C3DonutChartPage(tmplData)), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Wrote: %s\n", filename)
}
