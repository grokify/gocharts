package main

import (
	"fmt"
	"log"
	"time"

	"github.com/grokify/gocharts/data/statictimeseries"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/time/timeutil"
)

func main() {
	tbl, err := table.ReadFileSimple("data.csv", ",", true, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Record Count [%d]\n", len(tbl.Records))
	fmtutil.PrintJSON(tbl.Records)

	for _, rec := range tbl.Records {
		fmtutil.PrintJSON(rec)
	}

	cfg := statictimeseries.TableConfig{
		TimeColIdx:          1,
		TimeFormat:          time.RFC3339,
		CountColIdx:         0,
		SeriesSetNameColIdx: 2,
		SeriesNameColIdx:    3}

	counts, err := statictimeseries.ParseRecordsDataItems(tbl.Records, cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(counts)

	dss2 := statictimeseries.NewDataSeriesSet2("Data Series Sets Counts")
	dss2.Interval = timeutil.Month
	dss2.AddItems(counts...)

	fmtutil.PrintJSON(dss2)

	fmt.Println("DONE")
}
