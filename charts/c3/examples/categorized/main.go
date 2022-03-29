package main

import (
	"fmt"
	"log"

	"io/ioutil"

	"github.com/grokify/gocharts/v2/charts/c3"
)

/*

URL: https://c3js.org/samples/categorized.html

```js
	var chart = c3.generate({
    data: {
        columns: [
            ['data1', 30, 200, 100, 400, 150, 250, 50, 100, 250]
        ]
    },
    axis: {
        x: {
            type: 'category',
            categories: ['cat1', 'cat2', 'cat3', 'cat4', 'cat5', 'cat6', 'cat7', 'cat8', 'cat9']
        }
    }
});
```
*/

func main() {
	chart := c3.C3Chart{
		Bindto: "#chart",
		Data: c3.C3ChartData{
			Columns: [][]interface{}{{"data1", 30, 200, 100, 400, 150, 250, 50, 100, 250}},
		},
		Axis: c3.C3Axis{
			X: c3.C3AxisX{
				Type:       "Category",
				Categories: []string{"cat1", "cat2", "cat3", "cat4", "cat5", "cat6", "cat7", "cat8", "cat9"},
			},
		},
	}

	tmplData := c3.TemplateData{
		HeaderHTML:             "Category Axis",
		ReportName:             "Category Axis",
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
