package echarts

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/grokify/mogo/type/maputil"

	"github.com/grokify/gocharts/v2/data/histogram"
)

type ChartOptions struct {
	GlobalOptions *GlobalOptions
	SeriesOptions *SeriesOptions
}

type GlobalOptions struct {
	AxisPointer *opts.AxisPointer
	BarChart    *opts.BarChart
	Colors      *opts.Colors
	Legend      *opts.Legend
	Title       *opts.Title
}

func (opts GlobalOptions) Options() []charts.GlobalOpts {
	var out []charts.GlobalOpts
	if opts.AxisPointer != nil {
		out = append(out, charts.WithAxisPointerOpts(opts.AxisPointer))
	}
	if opts.Colors != nil {
		out = append(out, charts.WithColorsOpts(*opts.Colors))
	}
	if opts.Legend != nil {
		out = append(out, charts.WithLegendOpts(*opts.Legend))
	}
	if opts.Title != nil {
		out = append(out, charts.WithTitleOpts(*opts.Title))
	}
	return out
}

type SeriesOptions struct {
	BarChart *opts.BarChart
}

func (opts SeriesOptions) Options() []charts.SeriesOpts {
	var out []charts.SeriesOpts
	if opts.BarChart != nil {
		out = append(out, charts.WithBarChartOpts(*opts.BarChart))
	}
	return out
}

type BarSeries struct {
	Name string
	Data []opts.BarData
}

func NewBarHistogramSet(opts *ChartOptions, hs *histogram.HistogramSet, histNames, binNames []string, def int, horiziontal bool) (*charts.Bar, error) {
	bar := charts.NewBar()
	if opts != nil && opts.GlobalOptions != nil {
		if gopts := opts.GlobalOptions.Options(); len(gopts) > 0 {
			bar.SetGlobalOptions(gopts...)
		}
	}

	if len(binNames) == 0 {
		binNames = hs.BinNames()
	}

	bar.SetXAxis(binNames)

	if len(histNames) == 0 {
		if len(hs.Order) > 0 {
			histNames = hs.Order
		} else {
			histNames = maputil.Keys(hs.HistogramMap)
		}
	}

	seriesSlices, err := HistogramSetBarSeriesSlice(hs, histNames, binNames, def)
	if err != nil {
		return nil, err
	}

	var seriesOpts []charts.SeriesOpts
	if opts != nil && opts.SeriesOptions != nil {
		seriesOpts = opts.SeriesOptions.Options()
	}

	for _, series := range seriesSlices {
		bar.AddSeries(series.Name, series.Data, seriesOpts...)
	}

	if horiziontal {
		bar.XYReversal()
	}

	return bar, nil
}

func HistogramSetBarSeriesSlice(hs *histogram.HistogramSet, histNames, binNames []string, def int) ([]BarSeries, error) {
	var out []BarSeries
	if hs == nil {
		return out, histogram.ErrHistogramSetCannotBeNil
	}
	for _, histName := range histNames {
		if h, ok := hs.HistogramMap[histName]; ok {
			if bs, err := HistogramBarSeries(h, binNames, def); err != nil {
				return out, err
			} else {
				out = append(out, bs)
			}
		} else {
			bs := BarSeries{
				Name: histName,
				Data: []opts.BarData{},
			}
			for _, binName := range binNames {
				bs.Data = append(bs.Data, opts.BarData{Name: binName, Value: 0})
			}
			out = append(out, bs)
		}
	}
	return out, nil
}

func HistogramBarSeries(hist *histogram.Histogram, binNames []string, def int) (BarSeries, error) {
	if hist == nil {
		return BarSeries{}, histogram.ErrHistogramCannotBeNil
	}
	bs := BarSeries{
		Name: hist.Name,
		Data: []opts.BarData{},
	}
	counts := hist.BinValuesOrDefault(binNames, def)
	for i, binName := range binNames {
		bs.Data = append(bs.Data, opts.BarData{
			Name:  binName,
			Value: counts[i],
		})
	}
	return bs, nil
}
