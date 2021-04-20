// statictimeseriesdata provides tools for adding and formatting
// static time series data for reporting purposes.
package statictimeseries

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/simplego/time/timeutil"
)

type TableConfig struct {
	TimeColIdx          uint
	TimeFormat          string
	CountColIdx         uint
	SeriesSetNameColIdx int // optional. Set < 0 to discard.
	SeriesNameColIdx    int
	Interval            timeutil.Interval
}

func (cfg *TableConfig) GetTimeFormat() string {
	cfg.TimeFormat = strings.TrimSpace(cfg.TimeFormat)
	if len(cfg.TimeFormat) == 0 {
		return time.RFC3339
	}
	return cfg.TimeFormat
}

func ParseRecordsDataItems(records [][]string, cfg TableConfig) ([]DataItem, error) {
	items := []DataItem{}
	for i, rec := range records {
		item := DataItem{}
		if cfg.TimeColIdx >= uint(len(rec)) {
			return items, fmt.Errorf("row [%d] missing time index col [%d]", i, cfg.TimeColIdx)
		}
		dtRaw := rec[int(cfg.TimeColIdx)]
		dt, err := time.Parse(cfg.GetTimeFormat(), dtRaw)
		if err != nil {
			return items, fmt.Errorf("row [%d] col [%d] time error raw [%s] error [%s]", i, cfg.TimeColIdx, dtRaw, err.Error())
		}
		item.Time = dt

		if cfg.CountColIdx >= uint(len(rec)) {
			return items, fmt.Errorf("row [%d] missing count index [%d]", i, cfg.TimeColIdx)
		}
		countRaw := rec[int(cfg.CountColIdx)]
		count, err := strconv.Atoi(countRaw)
		if err != nil {
			return items, fmt.Errorf("row [%d] col [%d] count error raw [%s] error [%s]", i, cfg.TimeColIdx, countRaw, err.Error())
		}
		item.Value = int64(count)

		if cfg.SeriesSetNameColIdx >= 0 {
			if cfg.SeriesSetNameColIdx >= len(rec) {
				return items, fmt.Errorf("row [%d] missing group1 index [%d]", i, cfg.SeriesSetNameColIdx)
			}
			item.SeriesSetName = rec[int(cfg.SeriesSetNameColIdx)]
		}

		if cfg.SeriesNameColIdx >= 0 {
			if cfg.SeriesNameColIdx >= len(rec) {
				return items, fmt.Errorf("row [%d] missing group2 index [%d]", i, cfg.SeriesNameColIdx)
			}
			item.SeriesName = rec[int(cfg.SeriesNameColIdx)]
		}

		items = append(items, item)
	}
	return items, nil
}
