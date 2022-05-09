package main

import (
	"fmt"

	"github.com/grokify/gocharts/v2/charts/google"
)

// exampleData provides the example data from Google's documentation here:
// https://developers.google.com/chart/interactive/docs/gallery/linechart#examples
// `demo_reference.html` is a copy and paste from the Google documentation while
// `demo.html` was created by this program.
func exampleData() google.LineChartMaterial {
	return google.LineChartMaterial{
		Title:    "Box Office Earnings in First Two Weeks of Opening",
		Subtitle: "in millions of dollars (USD)",
		Columns: []google.Column{
			{
				Type: "number",
				Name: "Day"},
			{
				Type: "number",
				Name: "Guardians of the Galaxy'"},
			{
				Type: "number",
				Name: "The Avengers"},
			{
				Type: "number",
				Name: "Transformers: Age of Extinction"},
		},
		Data: [][]interface{}{
			[]interface{}{1, 37.8, 80.8, 41.8},
			[]interface{}{2, 30.9, 69.5, 32.4},
			[]interface{}{3, 25.4, 57, 25.7},
			[]interface{}{4, 11.7, 18.8, 10.5},
			[]interface{}{5, 11.9, 17.6, 10.4},
			[]interface{}{6, 8.8, 13.6, 7.7},
			[]interface{}{7, 7.6, 12.3, 9.6},
			[]interface{}{8, 12.3, 29.2, 10.6},
			[]interface{}{9, 16.9, 42.9, 14.8},
			[]interface{}{10, 12.8, 30.9, 11.6},
			[]interface{}{11, 5.3, 7.9, 4.7},
			[]interface{}{12, 6.6, 8.4, 5.2},
			[]interface{}{13, 4.8, 6.3, 3.6},
			[]interface{}{14, 4.2, 6.2, 3.4},
		},
	}
}

func main() {
	data := exampleData()
	pageHTML := google.LineChartMaterialPage(data)

	fmt.Println(pageHTML)
}
