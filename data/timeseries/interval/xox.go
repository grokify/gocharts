package interval

import (
	"sort"
	"time"

	tu "github.com/grokify/mogo/time/timeutil"
)

type XoxPoint struct {
	Time           time.Time
	TimeMonthAgo   time.Time
	TimeQuarterAgo time.Time
	TimeYearAgo    time.Time
	Value          int64
	YOldValue      int64
	QOldValue      int64
	MOldValue      int64
	YNowValue      int64
	QNowValue      int64
	MNowValue      int64
	MYAgoValue     int64
	MQAgoValue     int64
	MMAgoValue     int64
	AggregateValue int64
	YoY            float64
	QoQ            float64
	MoM            float64
	YoYAggregate   float64
	QoQAggregate   float64
	MoMAggregate   float64
}

type YoYQoQGrowth struct {
	DateMap map[string]XoxPoint
	YTD     int64
	QTD     int64
}

func NewYoYQoQGrowth(set TimeSeriesSet) (YoYQoQGrowth, error) {
	yoy := YoYQoQGrowth{DateMap: map[string]XoxPoint{}}

	seriesNames := set.SeriesNamesSorted()
	for _, seriesName := range seriesNames {
		if seriesName == set.AllSeriesName {
			continue
		}
		outputDataSeries, err := set.GetTimeSeries(seriesName, Output)
		if err != nil {
			return yoy, err
		}
		outputItems := outputDataSeries.ItemsSorted()

		aggregateTimeSeries, err := set.GetTimeSeries(seriesName, OutputAggregate)
		if err != nil {
			return yoy, err
		}
		aggregateItems := aggregateTimeSeries.ItemsSorted()

		for j, item := range outputItems {
			aggregateItem := aggregateItems[j]
			point := XoxPoint{
				Time:           item.Time,
				Value:          item.Value,
				AggregateValue: aggregateItem.Value,
				YoY:            0.0,
				QoQ:            0.0,
			}
			key := item.Time.Format(time.RFC3339)
			if existingPoint, ok := yoy.DateMap[key]; ok {
				trap := false
				if existingPoint.Value > 0 && point.Value > 0 {
					trap = false
				}
				point.Value += existingPoint.Value
				point.AggregateValue += existingPoint.AggregateValue
				yoy.DateMap[key] = point
				if trap {
					panic("GOT")
				}
			} else {
				yoy.DateMap[key] = point
			}
		}
	}

	for key, point := range yoy.DateMap {
		yearAgo := tu.PrevQuarters(point.Time, 4)
		yearKey := yearAgo.Format(time.RFC3339)
		quarterAgo := tu.PrevQuarter(point.Time)
		quarterKey := quarterAgo.Format(time.RFC3339)
		if yearPoint, ok := yoy.DateMap[yearKey]; ok {
			if yearPoint.Value > 0 {
				point.YoY = (float64(point.Value) - float64(yearPoint.Value)) / float64(yearPoint.Value)
				point.YoYAggregate = (float64(point.AggregateValue) - float64(yearPoint.AggregateValue)) / float64(yearPoint.AggregateValue)
			}
		}
		if quarterPoint, ok := yoy.DateMap[quarterKey]; ok {
			if quarterPoint.Value > 0 {
				point.QoQ = (float64(point.Value) - float64(quarterPoint.Value)) / float64(quarterPoint.Value)
				point.QoQAggregate = (float64(point.AggregateValue) - float64(quarterPoint.AggregateValue)) / float64(quarterPoint.AggregateValue)
			}
		}
		yoy.DateMap[key] = point
	}
	yoy = AddYtdAndQtd(yoy)
	return yoy, nil
}

func AddYtdAndQtd(yoy YoYQoQGrowth) YoYQoQGrowth {
	ytd := int64(0)
	qtd := int64(0)
	now := time.Now()
	qt := tu.QuarterStart(now)
	yr := tu.YearStart(now)
	for _, point := range yoy.DateMap {
		if tu.IsGreaterThan(point.Time, qt, true) {
			qtd += point.Value
		}
		if tu.IsGreaterThan(point.Time, yr, true) {
			ytd += point.Value
		}
	}
	yoy.YTD = ytd
	yoy.QTD = qtd
	return yoy
}

func (yoy *YoYQoQGrowth) ItemsSorted() []XoxPoint {
	keys := []string{}
	for key := range yoy.DateMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	points := []XoxPoint{}
	for _, key := range keys {
		if point, ok := yoy.DateMap[key]; ok {
			points = append(points, point)
		}
	}
	return points
}
