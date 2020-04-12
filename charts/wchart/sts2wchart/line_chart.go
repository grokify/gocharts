package sts2wchart

import (
	"time"

	"github.com/grokify/gocharts/charts/wchart"
	"github.com/grokify/gocharts/data/statictimeseries"
	"github.com/grokify/gotilla/math/mathutil"
	"github.com/grokify/gotilla/strconv/strconvutil"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func DataSeriesMonthToLineChart(ds statictimeseries.DataSeries) chart.Chart {
	mainSeries := wchart.DataSeriesToContinuousSeriesMonth(ds)
	polyRegSeries := &chart.PolynomialRegressionSeries{
		Degree:      3,
		InnerSeries: mainSeries}

	mainSeries.Style = chart.Style{
		StrokeWidth: float64(3)}

	graph := chart.Chart{
		Title: ds.SeriesName,
		YAxis: chart.YAxis{
			ValueFormatter: func(v interface{}) string {
				if vf, isFloat := v.(float64); isFloat {
					return strconvutil.Commify(int64(vf))
				}
				return ""
			},
		},
		Series: []chart.Series{
			mainSeries,
			polyRegSeries,
		},
	}

	//note we have to do this as a separate step because we need a reference to graph
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	ac := AxesCreator{
		PaddingTop: 50,
		GridMajorStyle: chart.Style{
			StrokeWidth: float64(1),
			StrokeColor: drawing.ColorFromHex("000000")},
		GridMinorStyle: chart.Style{
			StrokeWidth: float64(1),
			StrokeColor: drawing.ColorFromHex("aaaaaa")},
		XTickMonthFormat: "1/06",
		YNumTicks:        10,
		YTickFormatFunc:  strconvutil.Int64Abbreviation}

	minValue, maxValue := ds.MinMaxValues()
	minTime, maxTime := ds.MinMaxTimes()

	graph = ac.ChartAddAxesMonthDataSeries(
		graph, minTime, maxTime, minValue, maxValue)

	if 1 == 0 {
		tickValues := mathutil.PrettyTicks(10.0, minValue, maxValue)
		graph.YAxis.Ticks = wchart.Ticks(tickValues, strconvutil.Int64Abbreviation)

		style := chart.Style{
			StrokeWidth: float64(1),
			StrokeColor: drawing.ColorFromHex("aaaaaa")}
		styleMajor := chart.Style{
			StrokeWidth: float64(1),
			StrokeColor: drawing.ColorFromHex("000000")}

		graph.YAxis.GridLines = wchart.GridLines(tickValues, style)
		graph.YAxis.GridMajorStyle = style

		minTime, maxTime := ds.MinMaxTimes()
		xTicks, xGridlines := wchart.TicksAndGridlinesMonths(
			minTime, maxTime, styleMajor, style, "1/06", true)
		graph.XAxis.Ticks = xTicks
		graph.XAxis.GridLines = xGridlines
		graph.XAxis.GridMajorStyle = style
	}
	return graph
}

type AxesCreator struct {
	GridMajorStyle   chart.Style
	GridMinorStyle   chart.Style
	PaddingTop       int
	YNumTicks        int
	YTickFormatFunc  func(int64) string
	XTickMonthFormat string
}

func (ac *AxesCreator) ChartAddAxesMonthDataSeries(graph chart.Chart, minTime, maxTime time.Time, minValue, maxValue int64) chart.Chart {
	graph.Background = chart.Style{
		Padding: chart.Box{}}
	if ac.PaddingTop > 0 {
		graph.Background.Padding.Top = ac.PaddingTop
	}

	tickValues := mathutil.PrettyTicks(ac.YNumTicks, minValue, maxValue)
	graph.YAxis.Ticks = wchart.Ticks(tickValues, ac.YTickFormatFunc)
	graph.YAxis.GridLines = wchart.GridLines(tickValues, ac.GridMinorStyle)
	graph.YAxis.GridMajorStyle = ac.GridMinorStyle

	xTicks, xGridlines := wchart.TicksAndGridlinesMonths(
		minTime, maxTime, ac.GridMajorStyle, ac.GridMinorStyle, ac.XTickMonthFormat, true)
	graph.XAxis.Ticks = xTicks
	graph.XAxis.GridLines = xGridlines
	graph.XAxis.GridMajorStyle = ac.GridMajorStyle

	return graph
}
