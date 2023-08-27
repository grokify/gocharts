package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/time/timeutil"

	"github.com/grokify/gocharts/v2/charts/c3"
	"github.com/grokify/gocharts/v2/charts/c3/c3sts"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gocharts/v2/data/timeseries"
)

func TestData() timeseries.TimeSeriesSet {
	tss := timeseries.NewTimeSeriesSet("Funnel Chart Data")
	tss.Interval = timeutil.Month
	for i := 0; i < 12; i++ {
		dt := month.MonthStart(time.Now().UTC(), i)
		//dt := timeutil.DeltaQuarters(time.Now().UTC(), i)
		val1 := int64(float64(i*20)*rand.Float64() + 10)
		val2 := int64(float64(val1)*1 - (float64(5) * rand.Float64()))
		val3 := int64(float64(val1)*1 - (float64(2) * rand.Float64()))

		tss.AddItems(timeseries.TimeItem{
			SeriesName: "Funnel Stage 1",
			Time:       dt,
			Value:      val1})
		tss.AddItems(timeseries.TimeItem{
			SeriesName: "Funnel Stage 2",
			Time:       dt,
			Value:      val2})
		tss.AddItems(timeseries.TimeItem{
			SeriesName: "Funnel Stage 3",
			Time:       dt,
			Value:      val3})
	}
	return tss
}

func main() {
	dss := TestData()
	dss.Inflate()
	numCols := len(dss.Times)
	lowFirst := true
	rep := timeseries.Report(dss, numCols, lowFirst)

	fmtutil.PrintJSON(rep)
	axis := timeseries.ReportAxisX(dss, numCols,
		func(t time.Time) string { return timeutil.FormatQuarterYYYYQ(t) })

	c3Bar := c3.TimeSeriesSetSimpleToC3ChartBar(rep, c3.C3Bar{})
	c3Bar.Axis = c3.C3Axis{X: c3.C3AxisX{Type: "category", Categories: axis}}

	if 1 == 0 {
		out := c3.C3BarChartJS(c3Bar)
		err := ioutil.WriteFile("bar_old.html", []byte(out), 0644)
		if err != nil {
			log.Fatal(err)
			fmt.Printf("DONE")
		}
	}

	footerHTML := ""
	if 1 == 1 {
		tableRows, qoqData, funData := c3sts.DataRowsToTableRows(rep, axis, true, true, "Count", "QoQ", "Funnel")

		if 1 == 1 {
			domID := "qoqchart"
			qoqChart := c3.C3Chart{
				Bindto: "#" + domID,
				Data: c3.C3ChartData{
					Columns: [][]any{},
				},
				Axis: c3Bar.Axis,
				Grid: c3.C3Grid{Y: c3.C3GridLines{Show: true}}}

			for _, r := range qoqData {
				r2 := []any{}
				r2 = append(r2, r.Name)
				for _, v := range r.Values {
					r2 = append(r2, int(math.Round(strconvutil.ChangeToXoXPct(v))))
				}
				qoqChart.Data.Columns = append(qoqChart.Data.Columns, r2)
			}

			footerHTML += "<h2>QoQ Chart</h2>" + c3.C3ChartHTMLSimple(domID, qoqChart)
		}
		if 1 == 1 {
			domID := "funchart"
			funChart := c3.C3Chart{
				Bindto: "#" + domID,
				Data: c3.C3ChartData{
					Columns: [][]any{},
				},
				Axis: c3Bar.Axis,
				Grid: c3.C3Grid{Y: c3.C3GridLines{Show: true}}}

			for _, r := range funData {
				r2 := []any{}
				r2 = append(r2, r.Name)
				for _, v := range r.Values {
					r2 = append(r2, int(math.Round(strconvutil.ChangeToFunnelPct(v))))
				}
				funChart.Data.Columns = append(funChart.Data.Columns, r2)
			}

			footerHTML += "<h2>Funnel Chart</h2>" + c3.C3ChartHTMLSimple(domID, funChart)
		}

		tbl := table.Table{
			ID:    "funnelpct",
			Style: table.StyleSimple,
			Rows:  tableRows}
		footerHTML += "<h2>Stats</h2>" + table.SimpleTable(tbl)
	}

	templateData := c3.TemplateData{
		HeaderHTML: "Developer Funnel",
		ReportName: "Developer Funnel",
		C3Chart:    c3Bar,
		FooterHTML: footerHTML}

	out := c3.C3DonutChartPage(templateData)
	err := ioutil.WriteFile("bar.html", []byte(out), 0644)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("DONE")
	}

	fmt.Println("DONE")
}
