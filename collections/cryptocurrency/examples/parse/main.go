package main

import (
	"fmt"

	"github.com/grokify/gocharts/v2/collections/cryptocurrency"
	"github.com/grokify/gocharts/v2/data/timeseries"
	"github.com/grokify/gocharts/v2/data/yahoohistorical"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/math/accounting"
	"github.com/grokify/mogo/time/timeutil"
)

func main() {
	hd := cryptocurrency.HistoricalDataBTCUSDMonthly()
	err := procHistoricalData(hd)
	logutil.FatalErr(err)

	hd2 := cryptocurrency.HistoricalDataETHUSDMonthly()
	err2 := procHistoricalData(hd2)
	logutil.FatalErr(err2)

	fmt.Println("DONE")
}

func procHistoricalData(hd *yahoohistorical.HistoricalDataYahoo) error {
	tbl := hd.Table
	fmtutil.MustPrintJSON(tbl.Rows)
	fmtutil.MustPrintJSON(tbl.Columns)
	ts, err := hd.CloseData(timeutil.Month)
	if err != nil {
		return err
	}
	fmtutil.MustPrintJSON(ts)

	tblXox := ts.TableMonthXOX("Jan 2006", "", "", "", "", "",
		&timeseries.TableMonthXOXOpts{
			AddMOMGrowth: true,
			MOMGrowthPct: accounting.AnnualToMonthly(.3),
			MOMBaseMonth: timeutil.MustParse(timeutil.RFC3339FullDate, "2021-12-01"),
		})
	fmtutil.MustPrintJSON(tblXox.Rows)
	fmtutil.MustPrintJSON(tblXox.Columns)

	err = tblXox.WriteXLSX("_bitcoin.xlsx", "bitcoin data")
	if err != nil {
		return err
	}

	err = tblXox.WriteCSV("_bitcoin.csv")
	return err
}
