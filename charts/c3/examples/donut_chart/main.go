package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/grokify/gocharts/v2/charts/c3"
)

func main() {
	chart := c3.C3Chart{
		Bindto: "#chart",
		Data: c3.C3ChartData{
			Columns: [][]any{{"Yes", 60}, {"No", 40}},
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

	if err := os.WriteFile(filename, []byte(c3.C3DonutChartPage(tmplData)), 0600); err != nil {
		slog.Error(err.Error(), "filename", filename)
		os.Exit(1)
	}

	fmt.Printf("Wrote: %s\n", filename)
	os.Exit(0)
}
