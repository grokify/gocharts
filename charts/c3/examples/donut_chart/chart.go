package main

import (
	//"fmt"

	"io/ioutil"

	"github.com/grokify/gocharts/charts/c3"
)

func main() {
	chart := c3.C3Chart{
		Bindto: "#chart",
		Data: c3.C3ChartData{
			Columns: [][]interface{}{{"Yes", 60}, {"No", 40}},
			Type:    "donut",
		},
		Donut: c3.C3Donut{Title: "Votes"},
	}

	tmplData := c3.TemplateData{
		HeaderHTML:             "Donut Chart",
		ReportName:             "Donut Chart",
		ReportLink:             "",
		IncludeDataTable:       false,
		IncludeDataTableTotals: false,
		C3Chart:                chart}

	ioutil.WriteFile("chart.html", []byte(c3.C3DonutChartPage(tmplData)), 0644)

	//str := c3.TimeseriesHTML(c3.TimeseriesData{})
	//str := c3.AAA()
	//str = c3.TimeseriesPage(c3.TimeseriesPageData{})
	//fmt.Println(str)
}
