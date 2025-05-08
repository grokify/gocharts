package wchart

import (
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/time/quarter"
	"github.com/grokify/mogo/time/timeutil"
	chart "github.com/wcharczuk/go-chart/v2"

	"github.com/grokify/gocharts/v2/data/timeseries"
)

func TimeSeriesMapToContinuousSeriesMonths(dsm map[string]timeseries.TimeSeries, order []string) ([]chart.ContinuousSeries, error) {
	csSet := []chart.ContinuousSeries{}
	for _, seriesName := range order {
		if ds, ok := dsm[seriesName]; ok {
			if cs, err := TimeSeriesToContinuousSeries(ds); err != nil {
				return csSet, err
			} else {
				csSet = append(csSet, cs)
			}
		}
	}
	return csSet, nil
}

func TimeSeriesToContinuousSeries(ds timeseries.TimeSeries) (chart.ContinuousSeries, error) {
	series := chart.ContinuousSeries{
		Name:    ds.SeriesName,
		XValues: []float64{},
		YValues: []float64{}}

	items := ds.ItemsSorted()
	for _, item := range items {
		switch ds.Interval {
		case timeutil.IntervalMonth:
			if dtC, err := month.TimeToMonthContinuous(item.Time); err != nil {
				return series, err
			} else {
				series.XValues = append(series.XValues, float64(dtC))
			}
		case timeutil.IntervalQuarter:
			if dtC, err := quarter.TimeToQuarterContinuous(item.Time); err != nil {
				return series, err
			} else {
				series.XValues = append(series.XValues, float64(dtC))
			}
		default:
			series.XValues = append(series.XValues, float64(item.Time.Unix()))
		}
		if ds.IsFloat {
			series.YValues = append(series.YValues, item.ValueFloat)
		} else {
			series.YValues = append(series.YValues, float64(item.Value))
		}
	}
	return series, nil
}

func TimeSeriesMapToContinuousSeriesQuarters(dsm map[string]timeseries.TimeSeries, order []string) ([]chart.ContinuousSeries, error) {
	csSet := []chart.ContinuousSeries{}
	for _, seriesName := range order {
		if ds, ok := dsm[seriesName]; ok {
			if cs, err := TimeSeriesToContinuousSeriesQuarter(ds); err != nil {
				return csSet, err
			} else {
				csSet = append(csSet, cs)
			}
		}
	}
	return csSet, nil
}

func TimeSeriesToContinuousSeriesQuarter(ds timeseries.TimeSeries) (chart.ContinuousSeries, error) {
	series := chart.ContinuousSeries{
		Name:    ds.SeriesName,
		XValues: []float64{},
		YValues: []float64{}}

	items := ds.ItemsSorted()
	for _, item := range items {
		dtQuarterContinuous, err := quarter.TimeToQuarterContinuous(item.Time)
		if err != nil {
			return series, err
		}
		series.XValues = append(
			series.XValues,
			float64(dtQuarterContinuous))
		series.YValues = append(
			series.YValues,
			float64(item.Value))
	}
	return series, nil
}
