package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/grokify/gocharts/charts/wcharczuk"
	"github.com/grokify/gocharts/data/statictimeseries"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/gotilla/time/timeutil"
	chart "github.com/wcharczuk/go-chart"
)

func drawChartDSSSimple(res http.ResponseWriter, req *http.Request) {
	ds3 := statictimeseries.NewDataSeriesSetSimple()

	j := 0
	for i := -10; i <= 0; i++ {
		j++
		fmt.Println(i)
		item := statictimeseries.DataItem{
			SeriesName: "A Series",
			Time:       timeutil.MonthStart(time.Now().AddDate(0, i, 0)),
			Value:      int64(j)}
		ds3.AddItem(item)
	}
	fmtutil.PrintJSON(ds3)
	graph := wcharczuk.DSSSimpleToChart(ds3, "Jan '06")

	res.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, res)
}

func drawChart(res http.ResponseWriter, req *http.Request) {
	/*
	   This is an example of using the `TimeSeries` to automatically coerce time.Time values into a continuous xrange.
	   Note: chart.TimeSeries implements `ValueFormatterProvider` and as a result gives the XAxis the appropriate formatter to use for the ticks.
	*/
	formatter := wcharczuk.TimeFormatter{Layout: "Jan '06"}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			ValueFormatter: formatter.FormatTime,
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name: "ABC",
				XValues: []time.Time{
					time.Now().AddDate(0, 0, -10),
					time.Now().AddDate(0, 0, -9),
					time.Now().AddDate(0, 0, -8),
					time.Now().AddDate(0, 0, -7),
					time.Now().AddDate(0, 0, -6),
					time.Now().AddDate(0, 0, -5),
					time.Now().AddDate(0, 0, -4),
					time.Now().AddDate(0, 0, -3),
					time.Now().AddDate(0, 0, -2),
					time.Now().AddDate(0, 0, -1),
					time.Now(),
				},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0},
			},
		},
	}

	res.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, res)
}

func drawCustomChart(res http.ResponseWriter, req *http.Request) {
	/*
	   This is basically the other timeseries example, except we switch to hour intervals and specify a different formatter from default for the xaxis tick labels.
	*/
	graph := chart.Chart{
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeHourValueFormatter,
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: []time.Time{
					time.Now().Add(-10 * time.Hour),
					time.Now().Add(-9 * time.Hour),
					time.Now().Add(-8 * time.Hour),
					time.Now().Add(-7 * time.Hour),
					time.Now().Add(-6 * time.Hour),
					time.Now().Add(-5 * time.Hour),
					time.Now().Add(-4 * time.Hour),
					time.Now().Add(-3 * time.Hour),
					time.Now().Add(-2 * time.Hour),
					time.Now().Add(-1 * time.Hour),
					time.Now(),
				},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0},
			},
		},
	}

	res.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, res)
}

func main() {
	http.HandleFunc("/", drawChart)
	http.HandleFunc("/favicon.ico", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte{})
	})
	http.HandleFunc("/custom1", drawCustomChart)
	http.HandleFunc("/custom2", drawChartDSSSimple)
	http.ListenAndServe(":8080", nil)
}
