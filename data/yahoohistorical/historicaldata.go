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

func (hd *HistoricalData) TimeSeriesSet(interval timeutil.Interval) (timeseries.TimeSeriesSet, error) {
	tss := timeseries.NewTimeSeriesSet(hd.Table.Name)
	tsOpen, err := hd.OpenTimeSeries(interval)
	if err != nil {
		return tss, err
	}
	tsHigh, err := hd.HighTimeSeries(interval)
	if err != nil {
		return tss, err
	}
	tsClose, err := hd.CloseTimeSeries(interval)
	if err != nil {
		return tss, err
	}
	tsAdjClose, err := hd.AdjCloseTimeSeries(interval)
	if err != nil {
		return tss, err
	}
	tsVolume, err := hd.VolumeTimeSeries(interval)
	if err != nil {
		return tss, err
	}
	tsVolume.ConvertFloat64()
	tss.Interval = interval
	tss.IsFloat = true
	err = tss.AddSeries(tsOpen, tsHigh, tsClose, tsAdjClose, tsVolume)
	return tss, err
}

func (hd *HistoricalData) OpenTimeSeries(interval timeutil.Interval) (timeseries.TimeSeries, error) {
	return hd.columnData(interval, ColumnOpen)
}

func (hd *HistoricalData) HighTimeSeries(interval timeutil.Interval) (timeseries.TimeSeries, error) {
	return hd.columnData(interval, ColumnHigh)
}

func (hd *HistoricalData) CloseTimeSeries(interval timeutil.Interval) (timeseries.TimeSeries, error) {
	return hd.columnData(interval, ColumnClose)
}

func (hd *HistoricalData) AdjCloseTimeSeries(interval timeutil.Interval) (timeseries.TimeSeries, error) {
	return hd.columnData(interval, ColumnAdjClose)
}

func (hd *HistoricalData) VolumeTimeSeries(interval timeutil.Interval) (timeseries.TimeSeries, error) {
	return hd.columnData(interval, ColumnVolume)
}

func (hd *HistoricalData) columnData(interval timeutil.Interval, columnName string) (timeseries.TimeSeries, error) {
	ts := timeseries.NewTimeSeries(columnName)
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
		if interval == timeutil.IntervalYear && stringsutil.ReverseIndex(row[0], "-01-01") != 0 {
			continue
		} else if interval == timeutil.IntervalQuarter && stringsutil.ReverseIndex(row[0], "-01-01") != 0 &&
			stringsutil.ReverseIndex(row[0], "-04-01") != 0 &&
			stringsutil.ReverseIndex(row[0], "-07-01") != 0 &&
			stringsutil.ReverseIndex(row[0], "-10-01") != 0 {
			continue
		} else if interval == timeutil.IntervalMonth && stringsutil.ReverseIndex(row[0], "-01") != 0 {
			continue
		}
		dt, err := time.Parse(timeutil.RFC3339FullDate, row[0])
		if err != nil {
			return ts, err
		}
		if columnName == ColumnVolume {
			val, err := hd.Table.Columns.CellInt(columnName, row, false, 0)
			if err != nil {
				return ts, err
			}
			ts.AddInt64(dt, int64(val))
		} else {
			val, err := hd.Table.Columns.CellFloat64(columnName, row, false, 0)
			if err != nil {
				return ts, err
			}
			ts.AddFloat64(dt, val)
		}
	}
	return ts, nil
}
