package c3

import (
	"fmt"

	"github.com/grokify/gocharts/v2/data/slot"
	"github.com/grokify/gocharts/v2/data/timeseries"
)

func TimeSeriesSetSimpleToC3ChartBar(data []timeseries.RowInt64, c3BarInfo C3Bar) C3Chart {
	c3Chart := C3Chart{
		Data: C3ChartData{
			Columns: [][]interface{}{},
			Type:    "bar"},
		Bar: c3BarInfo}

	for _, r := range data {
		row := []interface{}{}
		row = append(row, r.Name)
		for _, v := range r.Values {
			row = append(row, v)
		}
		c3Chart.Data.Columns = append(c3Chart.Data.Columns, row)
	}

	return c3Chart
}

func SlotDataSeriesSetSimpleToC3ChartBar(input slot.SlotDataSeriesSetSimple, c3BarInfo C3Bar, hardMax int64) (C3Chart, error) {
	output := C3Chart{
		Data: C3ChartData{
			Columns: [][]interface{}{},
			Type:    "bar"},
		Bar: c3BarInfo}

	columns := [][]interface{}{}
	min, max := input.MinMaxX()
	if hardMax > 0 {
		max = hardMax
	}
	seriesNames := input.KeysSorted()

	for _, seriesName := range seriesNames {
		slotDataSeries, ok := input.SeriesSet[seriesName]
		if !ok {
			return output, fmt.Errorf("series name not found [%v]", seriesName)
		}

		column := []interface{}{seriesName}
		for i := min; i <= max; i++ {
			//fmt.Printf("%v ", i)
			x := i
			if y, ok := slotDataSeries.SeriesData[x]; ok {
				column = append(column, y)
			} else {
				column = append(column, 0)
			}
		}
		columns = append(columns, column)
	}
	output.Data.Columns = columns

	return output, nil
}
