package main

//go:generate go run main.go

import (
	"os"

	"github.com/grokify/gocharts/charts/wchart"
	"github.com/wcharczuk/go-chart"
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

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}
