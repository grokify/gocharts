package main

import (
	"fmt"
	"log"
	"time"

	"github.com/grokify/gocharts/data/statictimeseries"
	"github.com/grokify/gocharts/data/table"
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

	cfg := statictimeseries.TableConfig{
		TimeColIdx:          1,
		TimeFormat:          time.RFC3339,
		CountColIdx:         0,
		SeriesSetNameColIdx: 2,
		SeriesNameColIdx:    3}

	counts, err := statictimeseries.ParseRecordsDataItems(tbl.Rows, cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(counts)

	if cfg.SeriesSetNameColIdx >= 0 {
		dss2 := statictimeseries.NewDataSeriesSet2("Data Series Sets Counts")
		dss2.Interval = timeutil.Month
		dss2.AddItems(counts...)
		fmtutil.PrintJSON(dss2)
	} else {
		dss := statictimeseries.NewDataSeriesSet("Data Series Set Counts")
		dss.Interval = timeutil.Month
		dss.AddItems(counts...)
		dss.Inflate()
		fmtutil.PrintJSON(dss)
	}

	fmt.Println("DONE")
}
