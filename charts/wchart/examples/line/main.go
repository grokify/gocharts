package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/time/timeutil"

	"github.com/grokify/gocharts/v2/charts/wchart"
	"github.com/grokify/gocharts/v2/charts/wchart/sts2wchart"
	"github.com/grokify/gocharts/v2/data/timeseries"
)

func drawChartTSSSimple(res http.ResponseWriter, req *http.Request) {
	tss := timeseries.NewTimeSeriesSet("Example Data Series Set")

	j := 0
	for i := -10; i <= 0; i++ {
		j++
		fmt.Println(i)
		item := timeseries.TimeItem{
			SeriesName: "A Series",
			Time:       month.MonthStart(time.Now().AddDate(0, i, 0), 0),
			Value:      int64(j)}
		tss.AddItems(item)
	}
	fmtutil.MustPrintJSON(tss)
	graph, err := sts2wchart.TimeSeriesSetToLineChart(
		tss,
		&sts2wchart.LineChartOpts{
			XAxisTickFunc: func(t time.Time) string {
				return t.Format("Jan '06")
			}},
	)
	if err != nil {
		slog.Error(err.Error())
	}

	res.Header().Set(httputilmore.HeaderContentType, httputilmore.ContentTypeImagePNG)
	if err := graph.Render(chartdraw.PNG, res); err != nil {
		slog.Error(err.Error())
	}
}

func exampleTimeseries(chartName string) chartdraw.TimeSeries {
	return chartdraw.TimeSeries{
		Name: chartName,
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
	}
}

func drawChart(res http.ResponseWriter, req *http.Request) {
	// This is an example of using the `TimeSeries` to automatically coerce time.Time values into a continuous xrange.
	// Note: chartdraw.TimeSeries implements `ValueFormatterProvider` and as a result gives the XAxis the appropriate formatter to use for the ticks.
	formatter := wchart.TimeFormatter{Layout: "Jan '06"}
	graph := chartdraw.Chart{
		XAxis: chartdraw.XAxis{
			ValueFormatter: formatter.FormatTime,
		},
		Series: []chartdraw.Series{
			exampleTimeseries("ABC"),
		},
	}

	res.Header().Set(httputilmore.HeaderContentType, httputilmore.ContentTypeImagePNG)
	if err := graph.Render(chartdraw.PNG, res); err != nil {
		slog.Error(err.Error())
	}
}

func GetChartExampleDays() chartdraw.Chart {
	formatter := wchart.TimeFormatter{Layout: "Jan '06"}
	return chartdraw.Chart{
		XAxis: chartdraw.XAxis{
			ValueFormatter: formatter.FormatTime,
			GridLines: []chartdraw.GridLine{
				{
					Value: float64(time.Now().AddDate(0, 0, -6).Nanosecond()),
				},
			},
		},
		Series: []chartdraw.Series{
			exampleTimeseries("By Day"),
		},
	}
}

func GetChartExampleMonths() chartdraw.Chart {
	//formatter := wchart.TimeFormatter{Layout: "Jan '06"}
	formatter := wchart.TimeFormatter{Layout: timeutil.RFC3339FullDate}
	return chartdraw.Chart{
		XAxis: chartdraw.XAxis{
			ValueFormatter: formatter.FormatTime,
			GridLines: []chartdraw.GridLine{
				{
					Value: float64(time.Now().AddDate(0, 0, -6).Nanosecond()),
				},
			},
		},
		Series: []chartdraw.Series{
			chartdraw.TimeSeries{
				Name: "By Month",
				XValues: []time.Time{
					month.MonthStart(time.Now(), -10),
					month.MonthStart(time.Now(), -9),
					month.MonthStart(time.Now(), -8),
					month.MonthStart(time.Now(), -7),
					month.MonthStart(time.Now(), -6),
					month.MonthStart(time.Now(), -5),
					month.MonthStart(time.Now(), -4),
					month.MonthStart(time.Now(), -3),
					month.MonthStart(time.Now(), -2),
					month.MonthStart(time.Now(), -1),
					month.MonthStart(time.Now(), 0),
				},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0},
			},
		},
	}
}

func drawCustomChart(res http.ResponseWriter, req *http.Request) {
	// This is basically the other timeseries example, except we switch to hour intervals and specify a different formatter from default for the xaxis tick labels.
	graph := chartdraw.Chart{
		XAxis: chartdraw.XAxis{
			ValueFormatter: chartdraw.TimeHourValueFormatter,
		},
		Series: []chartdraw.Series{
			chartdraw.TimeSeries{
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

	res.Header().Set(httputilmore.HeaderContentType, httputilmore.ContentTypeImagePNG)
	if err := graph.Render(chartdraw.PNG, res); err != nil {
		slog.Error(err.Error())
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", drawChart)
	mux.HandleFunc("/favicon.ico", func(res http.ResponseWriter, req *http.Request) {
		if _, err := res.Write([]byte{}); err != nil {
			slog.Error(err.Error())
		}
	})
	mux.HandleFunc("/custom1", drawCustomChart)
	mux.HandleFunc("/custom2", drawChartTSSSimple)

	chart1 := GetChartExampleMonths()
	err := wchart.WritePNGFile("img_line_wrap.png", chart1)
	if err != nil {
		log.Fatal(err)
	}

	if err := writeFile("img_line_direct.png", chart1); err != nil {
		slog.Error(err.Error())
	}

	log.Fatal(httputilmore.ListenAndServeTimeouts(":8080", mux, time.Second))
}

func writeFile(filename string, ch chartdraw.Chart) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := ch.Render(chartdraw.PNG, f); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	} else {
		return nil
	}
}
