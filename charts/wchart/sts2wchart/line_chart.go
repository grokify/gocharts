package sts2wchart

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/charts/wchart"
	"github.com/grokify/gocharts/data/statictimeseries"
	"github.com/grokify/gocharts/data/statictimeseries/interval"
	"github.com/grokify/gotilla/math/mathutil"
	"github.com/grokify/gotilla/math/ratio"
	"github.com/grokify/gotilla/strconv/strconvutil"
	"github.com/grokify/gotilla/time/month"
	"github.com/grokify/gotilla/time/timeutil"
	"github.com/grokify/gotilla/type/number"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

const RatioTwoCol = ratio.RatioAcademy

// LineChartOpts is used for month and quarter interval charts.
type LineChartOpts struct {
	TitleSuffixCurrentValue     bool
	TitleSuffixCurrentValueFunc func(int64) string
	TitleSuffixCurrentDateFunc  func(time.Time) string
	Legend                      bool
	RegressionDegree            int
	NowAnnotation               bool
	MAgoAnnotation              bool
	QAgoAnnotation              bool
	YAgoAnnotation              bool
	AgoAnnotationPct            bool
	YAxisLeft                   bool
	YAxisMin                    float64
	YAxisMinEnable              bool
	XAxisTickFunc               func(time.Time) string
	XAxisTickInterval           timeutil.Interval // year, quarter, month
	XAxisGridInterval           timeutil.Interval
	Width                       uint64
	Height                      uint64
	AspectRatio                 float64
	Interval                    timeutil.Interval
}

func (opts *LineChartOpts) WantAnnotations() bool {
	return opts.NowAnnotation || opts.MAgoAnnotation ||
		opts.QAgoAnnotation || opts.YAgoAnnotation
}

func (opts *LineChartOpts) WantTitleSuffix() bool {
	return opts.TitleSuffixCurrentValue ||
		opts.TitleSuffixCurrentDateFunc != nil
}

var defaultLineChartOpts = &LineChartOpts{
	Legend:           true,
	RegressionDegree: 0,
	YAxisLeft:        true,
	YAxisMinEnable:   true,
	XAxisTickFunc: func(t time.Time) string {
		return t.Format("Jan '06")
	},
	XAxisGridInterval: timeutil.Quarter,
	XAxisTickInterval: timeutil.Year,
	NowAnnotation:     true,
	QAgoAnnotation:    false,
	YAgoAnnotation:    false,
	AgoAnnotationPct:  true,
	Interval:          timeutil.Month}

func GetDefaultLineChartOpts() *LineChartOpts {
	return defaultLineChartOpts
}

func DataSeriesToLineChart(ds statictimeseries.DataSeries, opts *LineChartOpts) (chart.Chart, error) {
	dss := statictimeseries.NewDataSeriesSet()
	dss.Name = ds.SeriesName
	dss.Interval = ds.Interval
	dss.IsFloat = ds.IsFloat
	dss.Series[ds.SeriesName] = ds
	dss.Inflate(true)
	return DataSeriesSetToLineChart(dss, opts)
}

func WriteLineChartDataSeriesSet(filename string, dss statictimeseries.DataSeriesSet, opts *LineChartOpts) error {
	chart, err := DataSeriesSetToLineChart(dss, opts)
	if err != nil {
		return err
	}
	return wchart.WritePNG(filename, chart)
}

func DataSeriesSetToLineChart(dss statictimeseries.DataSeriesSet, opts *LineChartOpts) (chart.Chart, error) {
	if opts == nil {
		opts = defaultLineChartOpts
	}
	titleParts := []string{dss.Name}
	if opts.WantTitleSuffix() && len(dss.Series) == 1 {
		ds, err := dss.GetSeriesByIndex(0)
		if err != nil {
			return chart.Chart{}, err
		}
		last, err := ds.Last()
		if err == nil {
			if opts.TitleSuffixCurrentDateFunc != nil {
				str := opts.TitleSuffixCurrentDateFunc(last.Time)
				if len(str) > 0 {
					titleParts = append(titleParts, " - "+str)
				}
			}
			if opts.TitleSuffixCurrentValue {
				if opts.TitleSuffixCurrentValueFunc != nil {
					titleParts = append(titleParts, " - "+opts.TitleSuffixCurrentValueFunc(last.Value))
				} else {
					titleParts = append(titleParts, " - "+strconv.Itoa(int(last.Value)))
				}
			}
		}
	}

	graph := chart.Chart{
		Title: strings.Join(titleParts, " "),
		YAxis: chart.YAxis{
			ValueFormatter: func(v interface{}) string {
				if vf, isFloat := v.(float64); isFloat {
					return strconvutil.Commify(int64(vf))
				}
				return ""
			},
		},
		Series: []chart.Series{},
	}

	if opts.YAxisLeft {
		graph.YAxis.AxisType = chart.YAxisSecondary // move Y axis to left.
	}

	if opts.Width > 0 && opts.Height > 0 {
		graph.Width = int(opts.Width)
		graph.Height = int(opts.Height)
	} else if opts.Width > 0 && opts.AspectRatio > 0 {
		graph.Width = int(opts.Width)
		graph.Height = int(ratio.WidthToHeight(float64(opts.Width), opts.AspectRatio))
	} else if opts.Height > 0 && opts.AspectRatio > 0 {
		graph.Height = int(opts.Height)
		graph.Width = int(ratio.HeightToWidth(float64(opts.Height), opts.AspectRatio))
	}

	mainSeries := chart.ContinuousSeries{}
	if opts.Interval == timeutil.Quarter || opts.Interval == timeutil.Month {
		if opts.Interval != dss.Interval {
			return chart.Chart{}, fmt.Errorf("E_INTERVAL_MISMATCH INPUT_INTERVAL [%s]", dss.Interval)
			//panic("opts.Interval dss.Interval mismatch")
		}
	}

	if len(dss.Order) == 0 {
		dss.Inflate(true)
	}
	for _, seriesName := range dss.Order {
		if ds, ok := dss.Series[seriesName]; ok {
			mainSeries = wchart.DataSeriesToContinuousSeries(ds)

			mainSeries.Style = chart.Style{StrokeWidth: float64(3)}

			graph.Series = append(graph.Series, mainSeries)
			//fmtutil.PrintJSON(mainSeries)
			//panic("Q")
			if opts.RegressionDegree == 1 {
				linRegSeries := &chart.LinearRegressionSeries{
					InnerSeries: mainSeries,
					Style: chart.Style{
						StrokeWidth: float64(2),
						StrokeColor: wchart.ColorOrange}}
				graph.Series = append(graph.Series, linRegSeries)
			} else if opts.RegressionDegree > 1 {
				polyRegSeries := &chart.PolynomialRegressionSeries{
					Degree:      opts.RegressionDegree,
					InnerSeries: mainSeries,
					Style: chart.Style{
						StrokeWidth: float64(2),
						StrokeColor: wchart.ColorOrange}}
				graph.Series = append(graph.Series, polyRegSeries)
			}
		} else {
			return chart.Chart{}, fmt.Errorf("E_SERIES_NAME_NOT_FOUND [%s]", seriesName)
		}
	}

	if opts.Legend {
		//note we have to do this as a separate step because we need a reference to graph
		graph.Elements = []chart.Renderable{
			chart.Legend(&graph),
		}
	}

	fmtXTickFunc := opts.XAxisTickFunc
	if fmtXTickFunc == nil {
		fmtXTickFunc = FormatXTickTimeFunc(dss.Interval)
	}

	axesCreator := AxesCreator{
		PaddingTop: 50,
		GridMajorStyle: chart.Style{
			StrokeWidth: float64(1),
			StrokeColor: drawing.ColorFromHex("000000")},
		GridMinorStyle: chart.Style{
			StrokeWidth: float64(1),
			StrokeColor: drawing.ColorFromHex("aaaaaa")},
		XAxisTickInterval: opts.XAxisTickInterval,
		XAxisGridInterval: opts.XAxisGridInterval,
		XTickFormatFunc:   fmtXTickFunc,
		YNumTicks:         7,
		YTickFormatFunc:   FormatYTickFunc(dss.Name)}

	minTime, maxTime := dss.MinMaxTimes()
	if !dss.IsFloat {
		minValue, maxValue := dss.MinMaxValues()
		if opts.YAxisMinEnable {
			minValue = int64(opts.YAxisMin)
		}
		graph = axesCreator.ChartAddAxesDataSeries(
			graph, dss.Interval, minTime, maxTime, minValue, maxValue)
	} else {
		graph = axesCreator.AddBackground(graph)
		graph = axesCreator.AddXAxis(graph, dss.Interval, minTime, maxTime)
		axesCreator.YTickFormatFunc = func(raw float64) string {
			return fmt.Sprintf("%.1f%%", raw*100)
		}
		minValue, maxValue := dss.MinMaxValuesFloat64()
		if opts.YAxisMinEnable {
			minValue = opts.YAxisMin
		}
		graph = axesCreator.AddYAxisPercent(graph, minValue, maxValue)
	}

	if opts.Interval == timeutil.Month {
		for _, ds := range dss.Series {
			annoSeries, err := dataSeriesMonthToAnnotations(ds, *opts)
			if err == nil && len(annoSeries.Annotations) > 0 {
				graph.Series = append(graph.Series, annoSeries)
			}
		}
	}
	return graph, nil
}

func dataSeriesMonthToAnnotations(ds statictimeseries.DataSeries, opts LineChartOpts) (chart.AnnotationSeries, error) {
	annoSeries := chart.AnnotationSeries{
		Annotations: []chart.Value2{},
		Style: chart.Style{
			StrokeWidth: float64(2),
			StrokeColor: wchart.MustParseColor("limegreen")},
	}

	if !opts.WantAnnotations() {
		return annoSeries, nil
	}

	xox, err := interval.NewXoXDataSeries(ds)
	if err != nil {
		return annoSeries, err
	}
	xoxLast := xox.Last()

	if opts.NowAnnotation {
		annoSeries.Annotations = append(annoSeries.Annotations, chart.Value2{
			XValue: float64(month.TimeToMonthContinuous(xoxLast.Time)),
			YValue: float64(xoxLast.Value),
			Label:  strconvutil.Int64Abbreviation(xoxLast.Value)})
	}
	if opts.MAgoAnnotation {
		annoSeries.Annotations = append(annoSeries.Annotations, chart.Value2{
			XValue: float64(month.TimeToMonthContinuous(xoxLast.TimeMonthAgo)),
			YValue: float64(xoxLast.MMAgoValue),
			Label:  "M: " + strconvutil.Int64Abbreviation(xoxLast.MMAgoValue)})
	}
	if opts.QAgoAnnotation {
		suffix := ""
		if opts.AgoAnnotationPct {
			suffix = fmt.Sprintf(", %d%%", int(xoxLast.QoQ))
		}
		annoSeries.Annotations = append(annoSeries.Annotations, chart.Value2{
			XValue: float64(month.TimeToMonthContinuous(xoxLast.TimeQuarterAgo)),
			YValue: float64(xoxLast.MQAgoValue),
			Label:  "Q: " + strconvutil.Int64Abbreviation(xoxLast.MQAgoValue) + suffix})
	}
	if opts.YAgoAnnotation {
		suffix := ""
		if opts.AgoAnnotationPct {
			suffix = fmt.Sprintf(", %d%%", int(xoxLast.YoY))
		}
		annoSeries.Annotations = append(annoSeries.Annotations, chart.Value2{
			XValue: float64(month.TimeToMonthContinuous(xoxLast.TimeYearAgo)),
			YValue: float64(xoxLast.MYAgoValue),
			Label:  "Y: " + strconvutil.Int64Abbreviation(xoxLast.MYAgoValue) + suffix})
	}
	return annoSeries, nil
}

func dataSeriesQuarterToAnnotations(ds statictimeseries.DataSeries, opts LineChartOpts) (chart.AnnotationSeries, error) {
	annoSeries := chart.AnnotationSeries{
		Annotations: []chart.Value2{},
		Style: chart.Style{
			StrokeWidth: float64(2),
			StrokeColor: wchart.MustParseColor("limegreen")},
	}

	if !opts.WantAnnotations() {
		return annoSeries, nil
	}

	xox, err := interval.NewXoXDataSeries(ds)
	if err != nil {
		return annoSeries, err
	}
	xoxLast := xox.Last()

	if opts.NowAnnotation {
		annoSeries.Annotations = append(annoSeries.Annotations, chart.Value2{
			XValue: float64(month.TimeToMonthContinuous(xoxLast.Time)),
			YValue: float64(xoxLast.Value),
			Label:  strconvutil.Int64Abbreviation(xoxLast.Value)})
	}
	if opts.MAgoAnnotation {
		annoSeries.Annotations = append(annoSeries.Annotations, chart.Value2{
			XValue: float64(month.TimeToMonthContinuous(xoxLast.TimeMonthAgo)),
			YValue: float64(xoxLast.MMAgoValue),
			Label:  "M: " + strconvutil.Int64Abbreviation(xoxLast.MMAgoValue)})
	}
	if opts.QAgoAnnotation {
		suffix := ""
		if opts.AgoAnnotationPct {
			suffix = fmt.Sprintf(", %d%%", int(xoxLast.QoQ))
		}
		annoSeries.Annotations = append(annoSeries.Annotations, chart.Value2{
			XValue: float64(month.TimeToMonthContinuous(xoxLast.TimeQuarterAgo)),
			YValue: float64(xoxLast.MQAgoValue),
			Label:  "Q: " + strconvutil.Int64Abbreviation(xoxLast.MQAgoValue) + suffix})
	}
	if opts.YAgoAnnotation {
		suffix := ""
		if opts.AgoAnnotationPct {
			suffix = fmt.Sprintf(", %d%%", int(xoxLast.YoY))
		}
		annoSeries.Annotations = append(annoSeries.Annotations, chart.Value2{
			XValue: float64(month.TimeToMonthContinuous(xoxLast.TimeYearAgo)),
			YValue: float64(xoxLast.MYAgoValue),
			Label:  "Y: " + strconvutil.Int64Abbreviation(xoxLast.MYAgoValue) + suffix})
	}
	return annoSeries, nil
}

func FormatXTickTimeFunc(interval timeutil.Interval) func(time.Time) string {
	if interval == timeutil.Month {
		return func(dt time.Time) string {
			return dt.Format("1/06")
			//	return dt.Format("Jan '06")
		}
	} else if interval == timeutil.Quarter {
		return func(dt time.Time) string {
			return timeutil.FormatQuarterYYYYQ(dt)
		}
	}
	return func(dt time.Time) string {
		return dt.Format("1/06")
	}
}

var rxMrr = regexp.MustCompile(`(?i)\bmrr\b`)

func FormatYTickFunc(seriesName string) func(float64) string {
	return func(val float64) string {
		abbr := strconvutil.Int64Abbreviation(int64(val))
		if rxMrr.MatchString(seriesName) {
			return "$" + abbr
		}
		return abbr
	}
}

type AxesCreator struct {
	GridMajorStyle     chart.Style
	GridMinorStyle     chart.Style
	PaddingTop         int
	YNumTicks          int
	XAxisTickInterval  timeutil.Interval // year, quarter, month
	XAxisGridInterval  timeutil.Interval
	XTickFormatFunc    func(time.Time) string
	YTickFormatFunc    func(float64) string
	YTickFormatFuncInt func(int64) string
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
		ac.GridMajorStyle, ac.GridMinorStyle, ac.XTickFormatFunc, ac.XAxisTickInterval, ac.XAxisGridInterval)
	graph.XAxis.Ticks = xTicks
	graph.XAxis.GridLines = xGridlines
	graph.XAxis.GridMajorStyle = ac.GridMajorStyle
	return graph
}

func (ac *AxesCreator) AddYAxis(graph chart.Chart, minValue, maxValue int64) chart.Chart {
	tickValues := number.SliceInt64ToFloat64(mathutil.PrettyTicks(ac.YNumTicks, minValue, maxValue))
	graph.YAxis.Ticks = wchart.Ticks(tickValues, ac.YTickFormatFunc)
	graph.YAxis.GridLines = wchart.GridLines(tickValues, ac.GridMinorStyle)
	graph.YAxis.GridMajorStyle = ac.GridMinorStyle
	return graph
}

func (ac *AxesCreator) AddYAxisPercent(graph chart.Chart, minValue, maxValue float64) chart.Chart {
	tickValues := mathutil.PrettyTicksPercent(ac.YNumTicks, minValue, maxValue, 2)
	graph.YAxis.Ticks = wchart.Ticks(tickValues, ac.YTickFormatFunc)
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
	graph.YAxis.Ticks = wchart.Ticks(tickValues, ac.YTickFormatFunc)
	graph.YAxis.GridLines = wchart.GridLines(tickValues, ac.GridMinorStyle)
	graph.YAxis.GridMajorStyle = ac.GridMinorStyle

	xTicks, xGridlines := wchart.TicksAndGridlinesTime(
		interval, minTime, maxTime,
		ac.GridMajorStyle, ac.GridMinorStyle, ac.XTickFormatFunc, ac.XAxisTickInterval, ac.XAxisGridInterval)
	graph.XAxis.Ticks = xTicks
	graph.XAxis.GridLines = xGridlines
	graph.XAxis.GridMajorStyle = ac.GridMajorStyle

	return graph
}
