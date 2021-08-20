package tblconv

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gocharts/data/timeseries"
)

type TableToTimeSeriesOpts struct {
	TimeColIdx   uint
	TimeFormat   string
	CountColIdx  uint
	CountIsFloat bool
}

func DefaultTableToTimeSeriesOpts() *TableToTimeSeriesOpts {
	return &TableToTimeSeriesOpts{
		TimeColIdx:   0,
		TimeFormat:   time.RFC3339,
		CountColIdx:  1,
		CountIsFloat: false}
}

func TableToTimeSeries(tbl table.Table, opts *TableToTimeSeriesOpts) (timeseries.TimeSeries, error) {
	if opts == nil {
		opts = DefaultTableToTimeSeriesOpts()
	}
	ts := timeseries.NewTimeSeries("")
	if opts.CountIsFloat {
		ts.IsFloat = true
	}
	dtFormat := opts.TimeFormat
	if len(dtFormat) == 0 {
		dtFormat = time.RFC3339
	}
	for _, row := range tbl.Rows {
		if len(row) == 0 {
			continue
		}
		if int(opts.TimeColIdx) >= len(row) {
			return ts, fmt.Errorf("time column doesn't exist")
		} else if int(opts.CountColIdx) >= len(row) {
			return ts, fmt.Errorf("count column doesn't exist")
		}
		dtStr := row[int(opts.TimeColIdx)]
		dt, err := time.Parse(dtFormat, dtStr)
		if err != nil {
			return ts, err
		}
		countStr := strings.TrimSpace(row[int(opts.CountColIdx)])
		if len(countStr) > 0 {
			if opts.CountIsFloat {
				countFloat, err := strconv.ParseFloat(countStr, 64)
				if err != nil {
					return ts, err
				}
				ts.AddFloat64(dt, countFloat)
			} else {
				countInt, err := strconv.Atoi(countStr)
				if err != nil {
					return ts, err
				}
				ts.AddInt64(dt, int64(countInt))
			}
		} else {
			ts.AddInt64(dt, 0)
		}
	}
	return ts, nil
}
