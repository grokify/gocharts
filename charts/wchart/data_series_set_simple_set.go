package wchart

import (
	"time"

	"github.com/grokify/mogo/time/timeutil"
	chart "github.com/go-analyze/charts/chartdraw"

	"github.com/grokify/gocharts/v2/data/timeseries"
)

type ChartOptions struct {
	LegendEnable    bool
	YAxisLeft       bool
	XAxisTimeLayout string
}

// TSSToChart converts a `TimeSeriesSet` to a `wcharczuk.Chart`.
func TSSToChart(data timeseries.TimeSeriesSet, opts ChartOptions) chart.Chart {
	formatter := TimeFormatter{Layout: opts.XAxisTimeLayout}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			ValueFormatter: formatter.FormatTime},
		Series: []chart.Series{},
	}
	for _, series := range data.Series {
		ts := chart.TimeSeries{
			Name: series.SeriesName}
		if 1 == 0 && opts.YAxisLeft {
			ts.YAxis = chart.YAxisSecondary
		}
		times := timeutil.TimeSeriesSlice(
			timeutil.IntervalMonth,
			series.ItemTimes())
		ts.XValues = times
		yvalues := []float64{}
		for _, t := range times {
			rfc := t.Format(time.RFC3339)
			if item, ok := series.ItemMap[rfc]; ok {
				yvalues = append(yvalues, float64(item.Value))
			} else {
				yvalues = append(yvalues, 0.0)
			}
		}
		ts.YValues = yvalues
		graph.Series = append(graph.Series, ts)
	}
	if opts.LegendEnable {
		graph.Elements = []chart.Renderable{
			chart.Legend(&graph),
		}
	}
	return graph
}
