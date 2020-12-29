package sts2wchart

import (
	"time"

	"github.com/grokify/gocharts/charts/wchart"

	"github.com/grokify/simplego/math/mathutil"
	"github.com/grokify/simplego/time/timeutil"
	"github.com/grokify/simplego/type/number"
	"github.com/wcharczuk/go-chart"
)

type AxesCreator struct {
	GridMajorStyle             chart.Style
	GridMinorStyle             chart.Style
	PaddingTop                 int
	YNumTicks                  int
	XAxisTickInterval          timeutil.Interval // year, quarter, month
	XAxisGridInterval          timeutil.Interval
	XAxisTickFormatFunc        func(time.Time) string
	YAxisTickFormatFuncFloat64 func(float64) string
	// YAxisTickFormatFuncInt64   func(int64) string
}

func (ac *AxesCreator) AddBackground(graph chart.Chart) chart.Chart {
	graph.Background = chart.Style{Padding: chart.Box{}}
	if ac.PaddingTop > 0 {
		graph.Background.Padding.Top = ac.PaddingTop
	}
	return graph
}

func (ac *AxesCreator) AddXAxis(graph chart.Chart, interval timeutil.Interval, minTime, maxTime time.Time) chart.Chart {
	xTicks, xGridlines := wchart.TicksAndGridlinesTime(
		interval, minTime, maxTime,
		ac.GridMajorStyle, ac.GridMinorStyle, ac.XAxisTickFormatFunc, ac.XAxisTickInterval, ac.XAxisGridInterval)
	graph.XAxis.Ticks = xTicks
	graph.XAxis.GridLines = xGridlines
	graph.XAxis.GridMajorStyle = ac.GridMajorStyle
	return graph
}

func (ac *AxesCreator) AddYAxis(graph chart.Chart, minValue, maxValue int64) chart.Chart {
	tickValues := number.SliceInt64ToFloat64(mathutil.PrettyTicks(ac.YNumTicks, minValue, maxValue))
	graph.YAxis.Ticks = wchart.Ticks(tickValues, ac.YAxisTickFormatFuncFloat64)
	graph.YAxis.GridLines = wchart.GridLines(tickValues, ac.GridMinorStyle)
	graph.YAxis.GridMajorStyle = ac.GridMinorStyle
	return graph
}

func (ac *AxesCreator) AddYAxisPercent(graph chart.Chart, minValue, maxValue float64) chart.Chart {
	tickValues := mathutil.PrettyTicksPercent(ac.YNumTicks, minValue, maxValue, 2)
	//graph.YAxis.Ticks = wchart.Ticks(tickValues, ac.YTickFormatFunc)
	graph.YAxis.Ticks = wchart.Ticks(tickValues, ac.YAxisTickFormatFuncFloat64)
	graph.YAxis.GridLines = wchart.GridLines(tickValues, ac.GridMinorStyle)
	graph.YAxis.GridMajorStyle = ac.GridMinorStyle
	return graph
}

func (ac *AxesCreator) ChartAddAxesDataSeries(graph chart.Chart, interval timeutil.Interval, minTime, maxTime time.Time, minValue, maxValue int64) chart.Chart {
	graph.Background = chart.Style{Padding: chart.Box{}}
	if ac.PaddingTop > 0 {
		graph.Background.Padding.Top = ac.PaddingTop
	}

	tickValues := number.SliceInt64ToFloat64(mathutil.PrettyTicks(ac.YNumTicks, minValue, maxValue))
	graph.YAxis.Ticks = wchart.Ticks(tickValues, ac.YAxisTickFormatFuncFloat64)
	graph.YAxis.GridLines = wchart.GridLines(tickValues, ac.GridMinorStyle)
	graph.YAxis.GridMajorStyle = ac.GridMinorStyle

	xTicks, xGridlines := wchart.TicksAndGridlinesTime(
		interval, minTime, maxTime,
		ac.GridMajorStyle, ac.GridMinorStyle, ac.XAxisTickFormatFunc, ac.XAxisTickInterval, ac.XAxisGridInterval)
	graph.XAxis.Ticks = xTicks
	graph.XAxis.GridLines = xGridlines
	graph.XAxis.GridMajorStyle = ac.GridMajorStyle

	return graph
}
