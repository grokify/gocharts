package yahoohistorical

import (
	"time"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gocharts/v2/data/timeseries"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	HistoricalDataColumns = "Date,Open,High,Low,Close,Adj Close,Volume"
	ColumnDate            = "Date"
	ColumnOpen            = "Open"
	ColumnHigh            = "High"
	ColumnClose           = "Close"
	ColumnAdjClose        = "Adj Close"
	ColumnVolume          = "Volume"
)

type HistoricalData struct {
	Table table.Table
}

// ReadFileHistoricalData reads a Yahoo! Finance Historical Data CSV file.
func ReadFileHistoricalData(filename string) (*HistoricalData, error) {
	tbl, err := table.ReadFile(nil, filename)
	if err != nil {
		return nil, err
	}
	tbl.Name = filename
	return &HistoricalData{Table: tbl}, nil
}

func (hd *HistoricalData) OpenData(interval timeutil.Interval) (timeseries.TimeSeries, error) {
	return hd.columnData(interval, ColumnOpen)
}

func (hd *HistoricalData) HighData(interval timeutil.Interval) (timeseries.TimeSeries, error) {
	return hd.columnData(interval, ColumnHigh)
}

func (hd *HistoricalData) CloseData(interval timeutil.Interval) (timeseries.TimeSeries, error) {
	return hd.columnData(interval, ColumnClose)
}

func (hd *HistoricalData) AdjCloseData(interval timeutil.Interval) (timeseries.TimeSeries, error) {
	return hd.columnData(interval, ColumnAdjClose)
}

func (hd *HistoricalData) VolumenData(interval timeutil.Interval) (timeseries.TimeSeries, error) {
	return hd.columnData(interval, ColumnVolume)
}

func (hd *HistoricalData) columnData(interval timeutil.Interval, columnName string) (timeseries.TimeSeries, error) {
	ts := timeseries.NewTimeSeries(hd.Table.Name)
	if columnName == ColumnVolume {
		ts.IsFloat = false
	} else {
		ts.IsFloat = true
	}
	ts.Interval = interval
	for _, row := range hd.Table.Rows {
		if len(row) == 0 {
			continue
		}
		if interval == timeutil.Year && stringsutil.ReverseIndex(row[0], "-01-01") != 0 {
			continue
		} else if interval == timeutil.Quarter && stringsutil.ReverseIndex(row[0], "-01-01") != 0 &&
			stringsutil.ReverseIndex(row[0], "-04-01") != 0 &&
			stringsutil.ReverseIndex(row[0], "-07-01") != 0 &&
			stringsutil.ReverseIndex(row[0], "-10-01") != 0 {
			continue
		} else if interval == timeutil.Month && stringsutil.ReverseIndex(row[0], "-01") != 0 {
			continue
		}
		dt, err := time.Parse(timeutil.RFC3339FullDate, row[0])
		if err != nil {
			return ts, err
		}
		if columnName == ColumnVolume {
			val, err := hd.Table.Columns.CellInt(columnName, row)
			if err != nil {
				return ts, err
			}
			ts.AddInt64(dt, int64(val))
		} else {
			val, err := hd.Table.Columns.CellFloat64(columnName, row)
			if err != nil {
				return ts, err
			}
			ts.AddFloat64(dt, val)
		}
	}
	return ts, nil
}
