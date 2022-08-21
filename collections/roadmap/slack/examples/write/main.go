package main

import (
	"fmt"
	"log"

	"github.com/grokify/gocharts/v2/collections/roadmap/slack"
	"github.com/grokify/mogo/fmt/fmtutil"
)

func main() {
	rmap := slack.GetRoadmapExample()

	streams, err := rmap.Streams(true, true)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(streams)

	tbl, err := rmap.Table(true, true)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(tbl)

	fmt.Printf("TBLNAME [%s]\n", tbl.Name)
	fmtutil.PrintJSON(tbl.Columns)
	fmtutil.PrintJSON(tbl.Rows)

	err = tbl.WriteXLSX("roadmap.xlsx", tbl.Name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE")
}
