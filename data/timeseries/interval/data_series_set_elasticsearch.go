package interval

import (
	"time"

	"github.com/grokify/gocharts/data/timeseries"
	"github.com/grokify/simplego/time/timeutil"

	v5 "github.com/grokify/elastirad-go/models/v5"
)

func EsAggsToDataSeriesSet(aggs []v5.AggregationResRad, interval timeutil.Interval, weekStart time.Weekday) DataSeriesSet {
	set := NewDataSeriesSet(interval, weekStart)

	for _, agg := range aggs {
		seriesName := agg.AggregationName
		for _, bucket := range agg.AggregationData.Buckets {
			set.AddItem(timeseries.TimeItem{
				SeriesName: seriesName,
				Time:       timeutil.UnixMillis(int64(bucket.Key.(float64))),
				Value:      int64(bucket.DocCount)})
		}
	}
	set.Inflate()
	return set
}
