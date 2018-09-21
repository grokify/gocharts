package c3

import (
	"fmt"

	"github.com/grokify/gocharts/data"
)

func SlotDataSeriesSetSimpleToC3ChartBar(input data.SlotDataSeriesSetSimple, c3BarInfo C3Bar) (C3Chart, error) {
	output := C3Chart{
		Data: C3ChartData{
			Columns: [][]interface{}{},
			Type:    "bar"},
		Bar: c3BarInfo,
	}

	columns := [][]interface{}{}
	min, max := input.MinMaxX()
	seriesNames := input.KeysSorted()

	for _, seriesName := range seriesNames {
		slotDataSeries, ok := input.SeriesSet[seriesName]
		if !ok {
			return output, fmt.Errorf("series name not found [%v]", seriesName)
		}
		fmt.Println(seriesName)
		column := []interface{}{seriesName}
		for i := min; i <= max; i++ {
			fmt.Printf("%v ", i)
			x := int64(i)
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
