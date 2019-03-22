package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/gotilla/time/timeutil"

	"github.com/grokify/gocharts/charts/c3"
	"github.com/grokify/gocharts/data/statictimeseries"
)

func getData(numQuarters int) []statictimeseries.DataItem {
	quarterNow := timeutil.QuarterStart(time.Now())
	quarterStart := timeutil.PrevQuarters(quarterNow, numQuarters)

	dataItems := []statictimeseries.DataItem{}

	for i := 0; i < numQuarters; i++ {
		dataItems = append(dataItems, statictimeseries.DataItem{
			SeriesName: "Data 1",
			Time:       timeutil.NextQuarters(quarterStart, i),
			Value:      int64(i),
		})
	}
	for i := 0; i < numQuarters; i++ {
		dataItems = append(dataItems, statictimeseries.DataItem{
			SeriesName: "Data 2",
			Time:       timeutil.NextQuarters(quarterStart, i),
			Value:      int64(i + 1),
		})
	}

	return dataItems
}

func main() {
	numQuarters := 5

	ds3 := statictimeseries.NewDataSeriesSetSimple()

	dataItems := getData(numQuarters)
	for _, di := range dataItems {
		ds3.AddItem(di)
	}
	ds3.Inflate()

	numCols := numQuarters
	lowFirst := true
	rep := statictimeseries.Report(ds3, numCols, lowFirst)
	fmtutil.PrintJSON(rep)
	axis := statictimeseries.ReportAxisX(ds3, numCols,
		func(t time.Time) string { return timeutil.FormatQuarterYYYYQ(t) })

	chart := c3.DataSeriesSetSimpleToC3ChartBar(rep, c3.C3Bar{})
	chart.Axis = c3.C3Axis{X: c3.C3AxisX{Type: "category", Categories: axis}}

	tmplData := c3.TemplateData{
		HeaderHTML:             "Bar Chart",
		ReportName:             "Bar Chart",
		ReportLink:             "",
		IncludeDataTable:       false,
		IncludeDataTableTotals: false,
		C3Chart:                chart}

	filename := "output.html"

	err := ioutil.WriteFile(filename, []byte(c3.C3DonutChartPage(tmplData)), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Wrote: %s\n", filename)
}
