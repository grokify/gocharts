package c3

import (
	"github.com/grokify/elastirad-go/models/v5"
)

type C3Chart struct {
	Bindto string      `json:"bindto,omitempty"`
	Data   C3ChartData `json:"data,omitempty"`
	Donut  C3Donut     `json:"donut,omitempty"`
}

type C3ChartData struct {
	Columns [][]interface{} `json:"columns,omitempty"`
	Type    string          `json:"type,omitempty"`
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
