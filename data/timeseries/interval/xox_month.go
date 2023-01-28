package interval

import (
	"sort"
	"time"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/math/mathutil"
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/time/timeutil"

	"github.com/grokify/gocharts/v2/data/timeseries"
)

type XoXGrowth struct {
	DateMap map[string]XoxPoint
	YTD     int64
	QTD     int64
}

func NewXoXTimeSeries(ds timeseries.TimeSeries) (XoXGrowth, error) {
	xox := XoXGrowth{DateMap: map[string]XoxPoint{}}
	for dateNowRfc, itemNow := range ds.ItemMap {
		dateNow, err := time.Parse(time.RFC3339, dateNowRfc)
		if err != nil {
			return xox, errorsutil.Wrap(err, "timeseries.NewXoXTimeSeries")
		}
		xoxPoint := XoxPoint{Time: dateNow, Value: itemNow.Int64()}

		quarterAgo := month.MonthStart(dateNow, -3)
		yearAgo := month.MonthStart(dateNow, -12)
		xoxPoint.TimeQuarterAgo = quarterAgo
		xoxPoint.TimeYearAgo = yearAgo
		if ds.Interval == timeutil.Month {
			monthAgo := month.MonthStart(dateNow, -1)
			xoxPoint.TimeMonthAgo = monthAgo
			if itemMonthAgo, ok := ds.ItemMap[monthAgo.Format(time.RFC3339)]; ok {
				xoxPoint.MMAgoValue = itemMonthAgo.Int64()
				xoxPoint.MNowValue = itemNow.Int64()
				xoxPoint.MOldValue = itemMonthAgo.Int64()
				xoxPoint.MoM = mathutil.PercentChangeToXoX(itemNow.Float64() / itemMonthAgo.Float64())
				xoxPoint.MoMAggregate = mathutil.PercentChangeToXoX(itemNow.Float64() / itemMonthAgo.Float64())
			}
		}
		if itemMonthQuarterAgo, ok := ds.ItemMap[quarterAgo.Format(time.RFC3339)]; ok {
			xoxPoint.MQAgoValue = itemMonthQuarterAgo.Int64()
			xoxPoint.QNowValue = AggregatePriorMonths(ds, dateNow, 3)
			xoxPoint.QOldValue = AggregatePriorMonths(ds, month.MonthStart(dateNow, -3), 3)
			xoxPoint.QoQ = mathutil.PercentChangeToXoX(itemNow.Float64() / itemMonthQuarterAgo.Float64())
			xoxPoint.QoQAggregate = mathutil.PercentChangeToXoX(
				float64(xoxPoint.QNowValue) / float64(xoxPoint.QOldValue))
		}
		if itemMonthYearAgo, ok := ds.ItemMap[yearAgo.Format(time.RFC3339)]; ok {
			xoxPoint.MYAgoValue = itemMonthYearAgo.Int64()
			xoxPoint.YNowValue = AggregatePriorMonths(ds, dateNow, 12)
			xoxPoint.YOldValue = AggregatePriorMonths(ds, month.MonthStart(dateNow, -12), 12)
			xoxPoint.YoY = mathutil.PercentChangeToXoX(itemNow.Float64() / itemMonthYearAgo.Float64())
			xoxPoint.YoYAggregate = mathutil.PercentChangeToXoX(
				float64(xoxPoint.YNowValue) / float64(xoxPoint.YOldValue))
			/*
				xoxPoint.YAgoValue = itemYear.Value
				xoxPoint.YoY = mathutil.PercentChangeToXoX(float64(itemNow.Value) / float64(itemYear.Value))
			*/
		}
		xox.DateMap[dateNowRfc] = xoxPoint
	}
	return xox, nil
}

func AggregatePriorMonths(ds timeseries.TimeSeries, start time.Time, months uint) int64 {
	aggregateValue := int64(0)
	monthStart := month.MonthStart(start, 0)
	for i := uint(1); i <= months; i++ {
		subtractMonths := i - 1
		thisMonth := monthStart
		if subtractMonths > 0 {
			thisMonth = month.MonthStart(monthStart, -1*int(subtractMonths))
		}
		key := thisMonth.Format(time.RFC3339)
		if item, ok := ds.ItemMap[key]; ok {
			aggregateValue += item.Value
		}
	}
	return aggregateValue
}

func (xg *XoXGrowth) Last() XoxPoint {
	dates := []string{}
	for date := range xg.DateMap {
		dates = append(dates, date)
	}
	if len(dates) == 0 {
		return XoxPoint{}
	}
	sort.Strings(dates)
	lastDate := dates[len(dates)-1]
	if lastItem, ok := xg.DateMap[lastDate]; ok {
		return lastItem
	}
	return XoxPoint{}
}
