package main

//go:generate go run main.go

import (
	"log/slog"
	"os"

	chart "github.com/go-analyze/charts/chartdraw"

	"github.com/grokify/gocharts/v2/charts/wchart"
)

func main() {
	graph := chart.BarChart{
		Title: "Test Bar Chart",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		ColorPalette: wchart.ColorsDefault(),
		Height:       512,
		//BarWidth: 60,
		Bars: []chart.Value{
			{Value: 5.25, Label: "Jan 20"},
			{Value: 4.88, Label: "Feb 20"},
			{Value: 4.74, Label: "Gray"},
			{Value: 3.22, Label: "Orange"},
			{Value: 3, Label: "Test"},
			{Value: 2.27, Label: "ABC"},
			{Value: 1, Label: "DEF"},
		},
	}

	f, err := os.Create("output.png")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer f.Close()
	if err := graph.Render(chart.PNG, f); err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}
	os.Exit(0)
}
