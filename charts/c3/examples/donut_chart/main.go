package main

import (
	"fmt"
	"io/ioutil"

	"github.com/grokify/gocharts/v2/charts/c3"
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

	filename := "output.html"

	ioutil.WriteFile(filename, []byte(c3.C3DonutChartPage(tmplData)), 0644)

	fmt.Printf("Wrote: %s\n", filename)
}
