package main

import (
	"fmt"
	"log"
	"time"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gocharts/data/timeseries"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/time/timeutil"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	File string `short:"f" long:"file" description:"Input OAS Spec File" required:"true"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	tbl, err := table.ReadFile(opts.File, ',', true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Record Count [%d]\n", len(tbl.Rows))
	fmtutil.PrintJSON(tbl.Rows)

	for _, row := range tbl.Rows {
		fmtutil.PrintJSON(row)
	}

	cfg := timeseries.TableConfig{
		TimeColIdx:          1,
		TimeFormat:          time.RFC3339,
		CountColIdx:         0,
		SeriesSetNameColIdx: 2,
		SeriesNameColIdx:    3}

	counts, err := timeseries.ParseRecordsDataItems(tbl.Rows, cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(counts)

	if cfg.SeriesSetNameColIdx >= 0 {
		tss2 := timeseries.NewTimeSeriesSet2("Time Series Sets Counts")
		tss2.Interval = timeutil.Month
		tss2.AddItems(counts...)
		fmtutil.PrintJSON(tss2)
	} else {
		tss := timeseries.NewTimeSeriesSet("Time Series Set Counts")
		tss.Interval = timeutil.Month
		tss.AddItems(counts...)
		tss.Inflate()
		fmtutil.PrintJSON(tss)
	}

	fmt.Println("DONE")
}
