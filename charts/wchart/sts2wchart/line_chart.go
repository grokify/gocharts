package sts2wchart

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/charts/wchart"
	"github.com/grokify/gocharts/data/statictimeseries"
	"github.com/grokify/gocharts/data/statictimeseries/interval"
	"github.com/grokify/simplego/math/ratio"
	"github.com/grokify/simplego/strconv/strconvutil"
	"github.com/grokify/simplego/time/month"
	"github.com/grokify/simplego/time/timeutil"
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
	XAxisGridInterval           timeutil.Interval
	XAxisTickFunc               func(time.Time) string
	XAxisTickInterval           timeutil.Interval // year, quarter, month
	YAxisTickFunc               func(float64) string
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
	Legend:            true,
	RegressionDegree:  0,
	XAxisGridInterval: timeutil.Quarter,
	XAxisTickFunc:     func(t time.Time) string { return t.Format("Jan '06") },
	XAxisTickInterval: timeutil.Year,
	YAxisLeft:         true,
	YAxisMinEnable:    true,
	YAxisTickFunc:     YAxisTickFormatSimple,
	NowAnnotation:     true,
	QAgoAnnotation:    true,
	YAgoAnnotation:    true,
	AgoAnnotationPct:  true,
	Interval:          timeutil.Month}

func DefaultLineChartOpts() *LineChartOpts {
	return defaultLineChartOpts
}

func DataSeriesToLineChart(ds statictimeseries.DataSeries, opts *LineChartOpts) (chart.Chart, error) {
	dss := statictimeseries.NewDataSeriesSet()
	dss.Name = ds.SeriesName
	dss.Interval = ds.Interval
	dss.IsFloat = ds.IsFloat
	dss.Series[ds.SeriesName] = ds
	dss.Inflate()
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
		dss.Inflate()
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
		XAxisTickInterval:          opts.XAxisTickInterval,
		XAxisGridInterval:          opts.XAxisGridInterval,
		XAxisTickFormatFunc:        fmtXTickFunc,
		YNumTicks:                  7,
		YAxisTickFormatFuncFloat64: opts.YAxisTickFunc}
	//YAxisTickFormatFuncFloat64: FormatYTickFunc(dss.Name)}

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
