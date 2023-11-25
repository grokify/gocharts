package main

import (
	"fmt"
	"time"

	"github.com/grokify/gocharts/v2/collections/cryptocurrency"
	"github.com/grokify/gocharts/v2/data/timeseries"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/time/timeutil"
)

func main() {
	hdBTC := cryptocurrency.HistoricalDataBTCUSDMonthly()
	adjClose, err := hdBTC.AdjCloseTimeSeries(timeutil.IntervalMonth)
	logutil.FatalErr(err, "hdBTC.AdjCloseTimeSeries")

	decsOnly := adjClose.ToMonth(false, time.December)
	decsOnly.SetSeriesName("actual")
	decsOnly.Interval = timeutil.IntervalYear
	err = decsOnly.TimeUpdateIntervalStart()
	logutil.FatalErr(err, "decsOnly.TimeUpdateIntervalStart")

	trainingSeries := decsOnly.Clone()
	trainingSeries.SetSeriesName("training")

	_, err = trainingSeries.Pop()
	logutil.FatalErr(err, "trainingSeries.Pop")
	_, err = trainingSeries.Pop()
	logutil.FatalErr(err, "trainingSeries.Pop")

	fmtutil.PrintJSON(trainingSeries)

	alpha, beta, err := trainingSeries.LinearRegression()
	logutil.FatalErr(err, "trainingSeries.LinearRegression")

	fmt.Printf("ALPHA [%v] BETA [%v]\n", alpha, beta)

	last, err := trainingSeries.Last()
	logutil.FatalErr(err, "trainingSeries.Last")

	forecastSeries := timeseries.NewTimeSeries("target")
	forecastSeries.Interval = decsOnly.Interval

	for i := 0; i < 2; i++ {
		x := last.Time.AddDate(i+1, 0, 0)
		forecastSeries.AddFloat64(x, alpha+beta*float64(x.Year()))
	}

	fmtutil.PrintJSON(forecastSeries)

	tss := timeseries.NewTimeSeriesSet("BTC")
	tss.Interval = timeutil.IntervalYear
	err = tss.AddSeries(decsOnly, trainingSeries, forecastSeries)
	logutil.FatalErr(err, "tss.AddSeries")

	tss.ActualTargetPairs = []timeseries.ActualTargetPair{{
		ActualSeriesName: "actual",
		TargetSeriesName: "target"}}
	tss.Inflate()
	fmtutil.MustPrintJSON(tss.Times)

	tbl, err := tss.TableActualTarget()
	logutil.FatalErr(err, "tss.TableActualTarget")

	err = tbl.WriteXLSX("example.xlsx", "BTC actual v. target")
	logutil.FatalErr(err)

	err = tbl.WriteCSV("example.csv")
	logutil.FatalErr(err)

	fmt.Println("DONE")
}
