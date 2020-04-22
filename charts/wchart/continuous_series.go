package wchart

import (
	"github.com/grokify/gocharts/data/statictimeseries"
	"github.com/grokify/gotilla/time/month"
	"github.com/grokify/gotilla/time/quarter"
	"github.com/wcharczuk/go-chart"
)

func DataSeriesMapToContinuousSeriesMonths(dsm map[string]statictimeseries.DataSeries, order []string) []chart.ContinuousSeries {
	csSet := []chart.ContinuousSeries{}
	for _, seriesName := range order {
		if ds, ok := dsm[seriesName]; ok {
			csSet = append(csSet, DataSeriesToContinuousSeriesMonth(ds))
		}
	}
	return csSet
}

func DataSeriesToContinuousSeriesMonth(ds statictimeseries.DataSeries) chart.ContinuousSeries {
	series := chart.ContinuousSeries{
		Name:    ds.SeriesName,
		XValues: []float64{},
		YValues: []float64{}}

	items := ds.ItemsSorted()
	for _, item := range items {
		series.XValues = append(
			series.XValues,
			float64(month.TimeToMonthContinuous(item.Time)))
		series.YValues = append(
			series.YValues,
			float64(item.Value))
	}
	return series
}

func DataSeriesMapToContinuousSeriesQuarters(dsm map[string]statictimeseries.DataSeries, order []string) []chart.ContinuousSeries {
	csSet := []chart.ContinuousSeries{}
	for _, seriesName := range order {
		if ds, ok := dsm[seriesName]; ok {
			csSet = append(csSet, DataSeriesToContinuousSeriesQuarter(ds))
		}
	}
	return csSet
}

func DataSeriesToContinuousSeriesQuarter(ds statictimeseries.DataSeries) chart.ContinuousSeries {
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
