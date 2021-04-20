// statictimeseriesdata provides tools for adding and formatting
// static time series data for reporting purposes.
package statictimeseries

import (
	"time"
)

type DataItem struct {
	SeriesName string
	Time       time.Time
	IsFloat    bool
	Value      int64
	ValueFloat float64
}

func (item *DataItem) ValueInt64() int64 {
	if item.IsFloat {
		return int64(item.ValueFloat)
	}
	return item.Value
}

func (item *DataItem) ValueFloat64() float64 {
	if item.IsFloat {
		return item.ValueFloat
	}
	return float64(item.Value)
}
