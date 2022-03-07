package c3

import (
	"encoding/json"

	"github.com/grokify/gocharts/data/timeseries/interval"
	"github.com/grokify/mogo/time/timeutil"
)

type TimeseriesData struct {
	IncludeTitle  bool
	Title         string
	TitleLevel    string
	DivID         string
	JSDataVar     string
	JSChartVar    string
	TimeSeriesSet *interval.TimeSeriesSet
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
	Xox       interval.YoYQoQGrowth
	XoxPoints []interval.XoxPoint
}

func (data *TimeseriesData) AddTimeSeriesSet(set *interval.TimeSeriesSet, interval timeutil.Interval, seriesType interval.SeriesType) error {
	data.TimeSeriesSet = set
	columns := [][]interface{}{}
	xValues := []interface{}{"x"}
	totals := []int64{}
	totalsMap := map[string]int64{}

	seriesNames := set.SeriesNamesSorted()
	for i, seriesName := range seriesNames {
		timeSeries, err := set.GetTimeSeries(seriesName, seriesType)
		if err != nil {
			return err
		}
		yValues := []interface{}{seriesName}

		items := timeSeries.ItemsSorted()
		for j, item := range items {
			if i == 0 && interval == timeutil.Quarter {
				xValues = append(xValues, item.Time.Format(timeutil.RFC3339FullDate))
			}
			yValues = append(yValues, item.Value)
			if j == len(totals) {
				totals = append(totals, int64(0))
			}
			totals[j] += item.Value
			rfc3339ym := item.Time.Format(timeutil.ISO8601YM)
			totalsMap[rfc3339ym] += item.Value
		}
		if i == 0 {
			columns = append(columns, xValues)
		}
		columns = append(columns, yValues)
	}
	tsj := TimeseriesDataJSON{
		Columns:   columns,
		Totals:    totals,
		TotalsMap: totalsMap}
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
