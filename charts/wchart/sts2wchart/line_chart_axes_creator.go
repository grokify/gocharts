package sts2wchart

import (
	"time"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/grokify/mogo/math/mathutil"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/number"

	"github.com/grokify/gocharts/v2/charts/wchart"
)

type AxesCreator struct {
	GridMajorStyle             chartdraw.Style
	GridMinorStyle             chartdraw.Style
	PaddingTop                 int
	YNumTicks                  int
	XAxisTickInterval          timeutil.Interval // year, quarter, month
	XAxisGridInterval          timeutil.Interval
	XAxisTickFormatFunc        func(time.Time) string
	YAxisTickFormatFuncFloat64 func(float64) string
	// YAxisTickFormatFuncInt64   func(int64) string
}

func (ac *AxesCreator) AddBackground(graph chartdraw.Chart) chartdraw.Chart {
	graph.Background = chartdraw.Style{Padding: chartdraw.Box{}}
	if ac.PaddingTop > 0 {
		graph.Background.Padding.Top = ac.PaddingTop
	}
	return graph
}

func (ac *AxesCreator) AddXAxis(graph chartdraw.Chart, interval timeutil.Interval, minTime, maxTime time.Time) (chartdraw.Chart, error) {
	xTicks, xGridlines, err := wchart.TicksAndGridlinesTime(
		interval, minTime, maxTime,
		ac.GridMajorStyle, ac.GridMinorStyle, ac.XAxisTickFormatFunc, ac.XAxisTickInterval, ac.XAxisGridInterval)
	if err != nil {
		return graph, err
	}
	graph.XAxis.Ticks = xTicks
	graph.XAxis.GridLines = xGridlines
	graph.XAxis.GridMajorStyle = ac.GridMajorStyle
	return graph, nil
}

func (ac *AxesCreator) AddYAxis(graph chartdraw.Chart, minValue, maxValue int64) chartdraw.Chart {
	tickValues := number.SliceToFloat64(mathutil.PrettyTicks(ac.YNumTicks, minValue, maxValue))
	graph.YAxis.Ticks = wchart.Ticks(tickValues, ac.YAxisTickFormatFuncFloat64)
	graph.YAxis.GridLines = wchart.GridLines(tickValues, ac.GridMinorStyle)
	graph.YAxis.GridMajorStyle = ac.GridMinorStyle
	return graph
}

func (ac *AxesCreator) AddYAxisPercent(graph chartdraw.Chart, minValue, maxValue float64) chartdraw.Chart {
	tickValues := mathutil.PrettyTicksPercent(ac.YNumTicks, minValue, maxValue, 2)
	//graph.YAxis.Ticks = wchart.Ticks(tickValues, ac.YTickFormatFunc)
	graph.YAxis.Ticks = wchart.Ticks(tickValues, ac.YAxisTickFormatFuncFloat64)
	graph.YAxis.GridLines = wchart.GridLines(tickValues, ac.GridMinorStyle)
	graph.YAxis.GridMajorStyle = ac.GridMinorStyle
	return graph
}

func (ac *AxesCreator) ChartAddAxesDataSeries(graph chartdraw.Chart, interval timeutil.Interval, minTime, maxTime time.Time, minValue, maxValue int64) (chartdraw.Chart, error) {
	graph.Background = chartdraw.Style{Padding: chartdraw.Box{}}
	if ac.PaddingTop > 0 {
		graph.Background.Padding.Top = ac.PaddingTop
	}

	tickValues := number.SliceToFloat64(mathutil.PrettyTicks(ac.YNumTicks, minValue, maxValue))
	graph.YAxis.Ticks = wchart.Ticks(tickValues, ac.YAxisTickFormatFuncFloat64)
	graph.YAxis.GridLines = wchart.GridLines(tickValues, ac.GridMinorStyle)
	graph.YAxis.GridMajorStyle = ac.GridMinorStyle

	xTicks, xGridlines, err := wchart.TicksAndGridlinesTime(
		interval, minTime, maxTime,
		ac.GridMajorStyle, ac.GridMinorStyle, ac.XAxisTickFormatFunc, ac.XAxisTickInterval, ac.XAxisGridInterval)
	if err != nil {
		return graph, err
	}
	graph.XAxis.Ticks = xTicks
	graph.XAxis.GridLines = xGridlines
	graph.XAxis.GridMajorStyle = ac.GridMajorStyle

	return graph, nil
}
