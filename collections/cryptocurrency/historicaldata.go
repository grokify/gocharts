package cryptocurrency

import (
	"bytes"
	"embed"
	"time"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gocharts/v2/data/timeseries"
	"github.com/grokify/gocharts/v2/data/yahoohistorical"
	"github.com/grokify/mogo/time/timeutil"
)

//go:embed BTC-USD_monthly_2022-04.csv
//go:embed ETH-USD_monthly_2022-04.csv
var f embed.FS

func TableBTCUSDMonthly() table.Table {
	data, err := f.ReadFile("BTC-USD_monthly_2022-04.csv")
	if err != nil {
		panic(err)
	}
	tbl, err := table.ParseReadSeeker(nil, bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	tbl.Name = "BTC-USD monthly"
	return tbl
}

func HistoricalDataBTCUSDMonthly() *yahoohistorical.HistoricalData {
	return &yahoohistorical.HistoricalData{Table: TableBTCUSDMonthly()}
}

func TableETHUSDMonthly() table.Table {
	data, err := f.ReadFile("ETH-USD_monthly_2022-04.csv")
	if err != nil {
		panic(err)
	}
	tbl, err := table.ParseReadSeeker(nil, bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	tbl.Name = "ETH-USD monthly"
	return tbl
}

func HistoricalDataETHUSDMonthly() *yahoohistorical.HistoricalData {
	return &yahoohistorical.HistoricalData{Table: TableETHUSDMonthly()}
}

func TableToTimeSeriesSet(t table.Table) (timeseries.TimeSeriesSet, error) {
	return timeseries.ParseTableTimeSeriesSetMatrixColumns(t, true, func(s string) (time.Time, error) {
		return time.Parse(timeutil.RFC3339FullDate, s)
	})
}
