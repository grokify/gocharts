package collections

import (
	"strings"

	"github.com/grokify/gocharts/v2/charts/wchart/sts2wchart"
	"github.com/grokify/gocharts/v2/data/timeseries"
	"github.com/grokify/gocharts/v2/data/yahoohistorical"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/math/accounting"
	"github.com/grokify/mogo/time/timeutil"
)

func WriteFilesHistoricalData(filePrefix string, hd *yahoohistorical.HistoricalData, verbose bool) error {
	if len(strings.TrimSpace(filePrefix)) == 0 {
		filePrefix = "_data"
	}
	tbl := hd.Table
	if verbose {
		fmtutil.MustPrintJSON(tbl.Rows)
		fmtutil.MustPrintJSON(tbl.Columns)
	}
	ts, err := hd.CloseTimeSeries(timeutil.IntervalMonth)
	if err != nil {
		return err
	}
	if verbose {
		fmtutil.MustPrintJSON(ts)
	}
	return WriteFilesTimeSeries(filePrefix, ts, verbose)
}

func WriteFilesTimeSeries(filePrefix string, ts timeseries.TimeSeries, verbose bool) error {
	if len(strings.TrimSpace(filePrefix)) == 0 {
		filePrefix = "_data"
	}

	tblXox := ts.TableMonthXOX("Jan 2006", "", "", "", "", "",
		&timeseries.TableMonthXOXOpts{
			AddMOMGrowth: true,
			MOMGrowthPct: accounting.AnnualToMonthly(.3),
			MOMBaseMonth: timeutil.MustParse(timeutil.RFC3339FullDate, "2021-12-01"),
		})
	if verbose {
		fmtutil.MustPrintJSON(tblXox.Rows)
		fmtutil.MustPrintJSON(tblXox.Columns)
	}
	err := tblXox.WriteXLSX(filePrefix+".xlsx", "data")
	if err != nil {
		return err
	}
	err = tblXox.WriteCSV(filePrefix + ".csv")
	if err != nil {
		return err
	}

	opts := sts2wchart.DefaultLineChartOpts()
	opts.XAxisTickInterval = timeutil.IntervalYear
	err = sts2wchart.WriteLineChartTimeSeries(filePrefix+".png", ts, opts)
	if err != nil {
		return err
	}
	return nil
}
