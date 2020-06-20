package interval

import (
	"sort"
	"time"

	"github.com/grokify/gocharts/data/statictimeseries"
	"github.com/grokify/gotilla/math/mathutil"
	"github.com/grokify/gotilla/time/month"
	"github.com/grokify/gotilla/time/timeutil"
	"github.com/pkg/errors"
)

type XoXGrowth struct {
	DateMap map[string]XoxPoint
	YTD     int64
	QTD     int64
}

func NewXoXDataSeries(ds statictimeseries.DataSeries) (XoXGrowth, error) {
	xox := XoXGrowth{DateMap: map[string]XoxPoint{}}
	for dateNowRfc, itemNow := range ds.ItemMap {
		dateNow, err := time.Parse(time.RFC3339, dateNowRfc)
		if err != nil {
			return xox, errors.Wrap(err, "statictimeseries.NewXoXDataSeries")
		}
		xoxPoint := XoxPoint{Time: dateNow, Value: itemNow.Value}

		quarterAgo := month.MonthBegin(dateNow, -3)
		yearAgo := month.MonthBegin(dateNow, -12)
		xoxPoint.TimeQuarterAgo = quarterAgo
		xoxPoint.TimeYearAgo = yearAgo
		if ds.Interval == timeutil.Month {
			monthAgo := month.MonthBegin(dateNow, -1)
			xoxPoint.TimeMonthAgo = monthAgo
			if itemMonthAgo, ok := ds.ItemMap[monthAgo.Format(time.RFC3339)]; ok {
				xoxPoint.MMAgoValue = itemMonthAgo.Value
				xoxPoint.MNowValue = itemNow.Value
				xoxPoint.MOldValue = itemMonthAgo.Value
				xoxPoint.MoM = mathutil.PercentChangeToXoX(float64(itemNow.Value) / float64(itemMonthAgo.Value))
				xoxPoint.MoMAggregate = mathutil.PercentChangeToXoX(float64(itemNow.Value) / float64(itemMonthAgo.Value))
			}
		}
		if itemMonthQuarterAgo, ok := ds.ItemMap[quarterAgo.Format(time.RFC3339)]; ok {
			xoxPoint.MQAgoValue = itemMonthQuarterAgo.Value
			xoxPoint.QNowValue = AggregatePriorMonths(ds, dateNow, 3)
			xoxPoint.QOldValue = AggregatePriorMonths(ds, month.MonthBegin(dateNow, -3), 3)
			xoxPoint.QoQ = mathutil.PercentChangeToXoX(float64(itemNow.Value) / float64(itemMonthQuarterAgo.Value))
			xoxPoint.QoQAggregate = mathutil.PercentChangeToXoX(
				float64(xoxPoint.QNowValue) / float64(xoxPoint.QOldValue))
		}
		if itemMonthYearAgo, ok := ds.ItemMap[yearAgo.Format(time.RFC3339)]; ok {
			xoxPoint.MYAgoValue = itemMonthYearAgo.Value
			xoxPoint.YNowValue = AggregatePriorMonths(ds, dateNow, 12)
			xoxPoint.YOldValue = AggregatePriorMonths(ds, month.MonthBegin(dateNow, -12), 12)
			xoxPoint.YoY = mathutil.PercentChangeToXoX(float64(itemNow.Value) / float64(itemMonthYearAgo.Value))
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

func AggregatePriorMonths(ds statictimeseries.DataSeries, start time.Time, months uint) int64 {
	aggregateValue := int64(0)
	monthBegin := month.MonthBegin(start, 0)
	for i := uint(1); i <= months; i++ {
		subtractMonths := i - 1
		thisMonth := monthBegin
		if subtractMonths > 0 {
			thisMonth = month.MonthBegin(monthBegin, -1*int(subtractMonths))
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
