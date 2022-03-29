package wchart

import (
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/time/quarter"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/wcharczuk/go-chart/v2"

	"github.com/grokify/gocharts/v2/data/timeseries"
)

func TimeSeriesMapToContinuousSeriesMonths(dsm map[string]timeseries.TimeSeries, order []string) []chart.ContinuousSeries {
	csSet := []chart.ContinuousSeries{}
	for _, seriesName := range order {
		if ds, ok := dsm[seriesName]; ok {
			csSet = append(csSet, TimeSeriesToContinuousSeries(ds))
		}
	}
	return csSet
}

func TimeSeriesToContinuousSeries(ds timeseries.TimeSeries) chart.ContinuousSeries {
	series := chart.ContinuousSeries{
		Name:    ds.SeriesName,
		XValues: []float64{},
		YValues: []float64{}}

	items := ds.ItemsSorted()
	for _, item := range items {
		switch ds.Interval {
		case timeutil.Month:
			series.XValues = append(series.XValues,
				float64(month.TimeToMonthContinuous(item.Time)))
		case timeutil.Quarter:
			series.XValues = append(series.XValues,
				float64(quarter.TimeToQuarterContinuous(item.Time)))
		default:
			series.XValues = append(series.XValues, float64(item.Time.Unix()))
		}
		if ds.IsFloat {
			series.YValues = append(series.YValues, item.ValueFloat)
		} else {
			series.YValues = append(series.YValues, float64(item.Value))
		}
	}
	return series
}

func TimeSeriesMapToContinuousSeriesQuarters(dsm map[string]timeseries.TimeSeries, order []string) []chart.ContinuousSeries {
	csSet := []chart.ContinuousSeries{}
	for _, seriesName := range order {
		if ds, ok := dsm[seriesName]; ok {
			csSet = append(csSet, TimeSeriesToContinuousSeriesQuarter(ds))
		}
	}
	return csSet
}

func TimeSeriesToContinuousSeriesQuarter(ds timeseries.TimeSeries) chart.ContinuousSeries {
	series := chart.ContinuousSeries{
		Name:    ds.SeriesName,
		XValues: []float64{},
		YValues: []float64{}}

	items := ds.ItemsSorted()
	for _, item := range items {
		series.XValues = append(
			series.XValues,
			float64(quarter.TimeToQuarterContinuous(item.Time)))
		series.YValues = append(
			series.YValues,
			float64(item.Value))
	}
	return series
}
