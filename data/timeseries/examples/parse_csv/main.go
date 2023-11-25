package main

import (
	"fmt"
	"log"
	"time"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/jessevdk/go-flags"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gocharts/v2/data/timeseries"
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

	tbl, err := table.ReadFile(nil, opts.File)
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

	counts, err := timeseries.ParseRecordsTimeItems(tbl.Rows, cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(counts)

	if cfg.SeriesSetNameColIdx >= 0 {
		sets := timeseries.NewTimeSeriesSets("Time Series Sets Counts")
		sets.Interval = timeutil.IntervalMonth
		sets.AddItems(counts...)
		fmtutil.PrintJSON(sets)
	} else {
		set := timeseries.NewTimeSeriesSet("Time Series Set Counts")
		set.Interval = timeutil.IntervalMonth
		set.AddItems(counts...)
		set.Inflate()
		fmtutil.PrintJSON(set)
	}

	fmt.Println("DONE")
}
