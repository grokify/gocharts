package google

import (
	"errors"
	"time"

	"github.com/grokify/gocharts/v2/data/timeseries"
	"github.com/grokify/mogo/type/stringsutil"
)

func TimeSeriesSetToLineChartMaterial(tss *timeseries.TimeSeriesSet, fn func(t time.Time) string) ([]Column, [][]any, error) {
	/*
		Returns data in Google Charts format
			Ref: https://developers.google.com/chart/interactive/docs/gallery/linechart#examples

			   var data = google.visualization.arrayToDataTable([
			     ['Year', 'Sales', 'Expenses'],
			     ['2004',  1000,      400],
			     ['2005',  1170,      460],
			     ['2006',  660,       1120],
			     ['2007',  1030,      540]
			   ]);
	*/

	if tss == nil {
		return nil, nil, errors.New("timeseries.TimeSeries cannot be nil")
	}
	var cols []Column
	var rows [][]any
	cols = append(cols, Column{
		Name: stringsutil.ToUpperFirst(tss.Interval.String(), false),
		Type: TypeString,
	})
	names := tss.SeriesNames()
	for _, seriesName := range names {
		cols = append(cols, Column{Name: seriesName, Type: TypeNumber})
	}
	times := tss.TimeSlice(true)
	for _, t := range times {
		var row []any
		if fn == nil {
			row = append(row, t.Format(time.RFC3339))
		} else {
			row = append(row, fn(t))
		}
		for _, seriesName := range names {
			row = append(row, tss.GetInt64WithDefault(seriesName, t.Format(time.RFC3339), 0))
		}
		rows = append(rows, row)
	}

	return cols, rows, nil
}
