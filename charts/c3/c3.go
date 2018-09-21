package c3

import (
	"encoding/json"

	"github.com/grokify/elastirad-go/models/v5"
)

const (
	ChartTypeBar = "bar"
)

type C3Chart struct {
	Bindto string      `json:"bindto,omitempty"`
	Data   C3ChartData `json:"data,omitempty"`
	Donut  C3Donut     `json:"donut,omitempty"`
	Bar    C3Bar       `json:"bar,omitempty"`
}

type C3ChartData struct {
	Columns [][]interface{} `json:"columns,omitempty"`
	Type    string          `json:"type,omitempty"`
}

func (data *C3ChartData) MustJSON() []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return bytes
}

type C3Donut struct {
	Title string `json:"title,omitempty"`
}

func C3ChartForEsAggregationSimple(agg v5.AggregationResRad) C3Chart {
	c3Chart := C3Chart{
		Data: C3ChartData{
			Columns: [][]interface{}{},
		},
	}
	for _, bucket := range agg.AggregationData.Buckets {
		c3Column := []interface{}{bucket.Key, bucket.DocCount}
		c3Chart.Data.Columns = append(c3Chart.Data.Columns, c3Column)
	}
	return c3Chart
}

/*
var chart = c3.generate({
    data: {
        columns: [
            ['data1', 30, 200, 100, 400, 150, 250],
            ['data2', 130, 100, 140, 200, 150, 50]
        ],
        type: 'bar'
    },
    bar: {
        width: {
            ratio: 0.5 // this makes bar width 50% of length between ticks
        }
        // or
        //width: 100 // this makes bar width 100px
    }
});
*/

type C3Bar struct {
	WidthRatio float64
	Width      int
}

type C3ColumnsInt struct {
	Columns []C3ColumnInt
}

func (cols *C3ColumnsInt) ToC3ChartData(chartType string) C3ChartData {
	columns := [][]interface{}{}
	for _, col := range cols.Columns {
		row := []interface{}{}
		row = append(row, col.Name)
		for _, val := range col.Values {
			row = append(row, val)
		}
	}

	return C3ChartData{
		Columns: columns,
		Type:    chartType}
}

type C3ColumnInt struct {
	Name   string
	Values []int
}
