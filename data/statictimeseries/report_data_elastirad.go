// statictimeseriesdata provides tools for adding and formatting
// static time series data for reporting purposes.
package statictimeseriesdata

import (
	"time"

	"github.com/grokify/gotilla/time/timeutil"

	"github.com/grokify/elastirad-go/models/v5"
)

func EsAggsToDataSeriesSet(aggs []v5.AggregationResRad, interval timeutil.Interval, weekStart time.Weekday) DataSeriesSet {
	set := NewDataSeriesSet(interval, weekStart)

	for _, agg := range aggs {
		seriesName := agg.AggregationName
		for _, bucket := range agg.AggregationData.Buckets {
			set.AddItem(DataItem{
				SeriesName: seriesName,
				Time:       timeutil.UnixMillis(int64(bucket.Key.(float64))),
				Value:      int64(bucket.DocCount)})
		}
	}
	set.Inflate()
	return set
}
