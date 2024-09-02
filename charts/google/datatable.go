package google

import (
	"errors"
	"time"

	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gocharts/v2/data/timeseries"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/time/timeutil"
)

type DataTable [][]any

func (dt DataTable) MustJSON() []byte {
	return jsonutil.MustMarshalOrDefault(dt, []byte(jsonutil.EmptyArray))
}

// DataTableFromHistogram is tested with barchart.
func DataTableFromHistogram(h *histogram.Histogram, inclUnordered, inclZeroCount, inclZeroCountTail bool) (DataTable, error) {
	dt := DataTable{}
	if h == nil {
		return dt, errors.New("histogram must be supplied")
	}
	cols := []any{h.Name, "Count"}
	dt = append(dt, cols)

	bins := h.OrderOrDefault(inclUnordered)
	idxDtLastNonZero := -1
	for _, binName := range bins {
		cnt := h.GetOrDefault(binName, 0)
		if cnt != 0 {
			dt = append(dt, []any{binName, cnt})
			idxDtLastNonZero = len(dt) - 1
		} else if inclZeroCount || inclZeroCountTail {
			dt = append(dt, []any{binName, cnt})
		}
	}
	if !inclZeroCountTail {
		return dt[:idxDtLastNonZero+1], nil
	} else {
		return dt, nil
	}
}

func DataTableFromTimeSeriesSet(name string, sets []string, set timeseries.TimeSeriesSet) (DataTable, error) {
	dt := DataTable{}
	if len(sets) == 0 {
		sets = set.SeriesNames()
	}
	row1 := []any{name}
	for _, set := range sets {
		row1 = append(row1, set)
	}
	row1 = append(row1, map[string]string{"role": "annotation"})
	dt = append(dt, row1)
	if set.Interval == timeutil.IntervalMonth {
		timeStrings := set.TimeStrings()
		for _, ts := range timeStrings {
			t, err := time.Parse(time.RFC3339, ts)
			if err != nil {
				return dt, err
			}
			mDisplay := t.Format("Jan 2006")
			row := []any{mDisplay}
			for _, sname := range sets {
				val := set.GetInt64WithDefault(sname, ts, 0)
				row = append(row, val)
			}
			row = append(row, "")
			dt = append(dt, row)
		}
	}
	return dt, nil
}
