package main

import (
	"fmt"
	"time"

	"github.com/grokify/gocharts/v2/apps/dynatrace"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
)

func main() {
	f := "path/to/dynatrace_extract.csv"

	t, err := dynatrace.ReadCSVTimeDurations(f, time.Millisecond, 2, true)
	logutil.FatalErr(err)
	fmtutil.MustPrintJSON(t.Rows)
	fmtutil.MustPrintJSON(t.Columns)

	if 1 == 0 {
		t, err := dynatrace.ReadCSVTimeDurations(f, time.Millisecond, 2, true)
		logutil.FatalErr(err)
		fmtutil.MustPrintJSON(t.Rows)
		fmtutil.MustPrintJSON(t.Columns)
	}

	fmt.Println("DONE")
}
