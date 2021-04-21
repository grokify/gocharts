package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/grokify/gocharts/charts/c3"
	"github.com/grokify/gocharts/charts/c3/c3sts"
	"github.com/grokify/gocharts/data/statictimeseries"
	"github.com/grokify/gocharts/data/table"

	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/strconv/strconvutil"
	"github.com/grokify/simplego/time/month"
	"github.com/grokify/simplego/time/timeutil"
)

func TestData() statictimeseries.DataSeriesSet {
	dss := statictimeseries.NewDataSeriesSet("Funnel Chart Data")
	dss.Interval = timeutil.Month
	for i := 0; i < 12; i++ {
		dt := month.MonthBegin(time.Now().UTC(), i)
		//dt := timeutil.DeltaQuarters(time.Now().UTC(), i)
		val1 := int64(float64(i*20)*rand.Float64() + 10)
		val2 := int64(float64(val1)*1 - (float64(5) * rand.Float64()))
		val3 := int64(float64(val1)*1 - (float64(2) * rand.Float64()))

		dss.AddItem(statictimeseries.DataItem{
			SeriesName: "Funnel Stage 1",
			Time:       dt,
			Value:      val1})
		dss.AddItem(statictimeseries.DataItem{
			SeriesName: "Funnel Stage 2",
			Time:       dt,
			Value:      val2})
		dss.AddItem(statictimeseries.DataItem{
			SeriesName: "Funnel Stage 3",
			Time:       dt,
			Value:      val3})
	}
	return dss
}

func main() {
	dss := TestData()
	dss.Inflate()
	numCols := len(dss.Times)
	lowFirst := true
	rep := statictimeseries.Report(dss, numCols, lowFirst)

	fmtutil.PrintJSON(rep)
	axis := statictimeseries.ReportAxisX(dss, numCols,
		func(t time.Time) string { return timeutil.FormatQuarterYYYYQ(t) })

	c3Bar := c3.DataSeriesSetSimpleToC3ChartBar(rep, c3.C3Bar{})
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
					Columns: [][]interface{}{},
				},
				Axis: c3Bar.Axis,
				Grid: c3.C3Grid{Y: c3.C3GridLines{Show: true}}}

			for _, r := range qoqData {
				r2 := []interface{}{}
				r2 = append(r2, r.Name)
				for _, v := range r.Values {
					r2 = append(r2, int(math.Round(strconvutil.ChangeToXoXPct(v))))
				}
				qoqChart.Data.Columns = append(qoqChart.Data.Columns, r2)
			}

			footerHTML += "<h2>QoQ Chart</h2>" + c3.C3ChartHtmlSimple(domID, qoqChart)
		}
		if 1 == 1 {
			domID := "funchart"
			funChart := c3.C3Chart{
				Bindto: "#" + domID,
				Data: c3.C3ChartData{
					Columns: [][]interface{}{},
				},
				Axis: c3Bar.Axis,
				Grid: c3.C3Grid{Y: c3.C3GridLines{Show: true}}}

			for _, r := range funData {
				r2 := []interface{}{}
				r2 = append(r2, r.Name)
				for _, v := range r.Values {
					r2 = append(r2, int(math.Round(strconvutil.ChangeToFunnelPct(v))))
				}
				funChart.Data.Columns = append(funChart.Data.Columns, r2)
			}

			footerHTML += "<h2>Funnel Chart</h2>" + c3.C3ChartHtmlSimple(domID, funChart)
		}

		tbl := table.Table{
			ID:      "funnelpct",
			Style:   table.StyleSimple,
			Records: tableRows}
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
