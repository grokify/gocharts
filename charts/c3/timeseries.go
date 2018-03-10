package c3

import (
	"encoding/json"

	sts "github.com/grokify/gocharts/data/statictimeseries"
	tu "github.com/grokify/gotilla/time/timeutil"
)

type TimeseriesData struct {
	IncludeTitle  bool
	Title         string
	TitleLevel    string
	DivID         string
	JSDataVar     string
	JSChartVar    string
	DataSeriesSet *sts.DataSeriesSet
	JSONData      TimeseriesDataJSON
}

type TimeseriesDataJSON struct {
	Columns   [][]interface{}  `json:"columns"`
	Totals    []int64          `json:"totals"`
	TotalsMap map[string]int64 `json:"totalsMap"`
}

func (tdj *TimeseriesDataJSON) JSON() []byte {
	bytes, err := json.Marshal(tdj)
	if err != nil {
		panic(err)
	}
	return bytes
}

type TimeseriesPageData struct {
	Title     string
	URL       string
	Charts    []TimeseriesData
	Xox       sts.YoYQoQGrowth
	XoxPoints []sts.XoxPoint
}

func (data *TimeseriesData) AddDataSeriesSet(set *sts.DataSeriesSet, interval tu.Interval, seriesType sts.SeriesType) error {
	data.DataSeriesSet = set
	columns := [][]interface{}{}
	xValues := []interface{}{"x"}
	totals := []int64{}
	totalsMap := map[string]int64{}

	seriesNames := set.SeriesNamesSorted()
	for i, seriesName := range seriesNames {
		dataSeries, err := set.GetDataSeries(seriesName, seriesType)
		if err != nil {
			return err
		}
		yValues := []interface{}{seriesName}

		items := dataSeries.SortedItems()
		for j, item := range items {
			if i == 0 && interval == tu.Quarter {
				xValues = append(xValues, item.Time.Format(tu.RFC3339YMD))
			}
			yValues = append(yValues, item.Value)
			if j == len(totals) {
				totals = append(totals, int64(0))
			}
			totals[j] += item.Value
			rfc3339ym := item.Time.Format(tu.ISO8601YM)
			if _, ok := totalsMap[rfc3339ym]; ok {
				totalsMap[rfc3339ym] += item.Value
			} else {
				totalsMap[rfc3339ym] = item.Value
			}
		}
		if i == 0 {
			columns = append(columns, xValues)
		}
		columns = append(columns, yValues)
	}
	tsj := TimeseriesDataJSON{
		Columns:   columns,
		Totals:    totals,
		TotalsMap: totalsMap,
	}
	data.JSONData = tsj

	return nil
}

func (data *TimeseriesData) FormattedDataJSON() []byte {
	bytes, err := json.Marshal(data.JSONData.Columns)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (data *TimeseriesData) DataJSON() []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return bytes
}
