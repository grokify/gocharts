package main

import (
	"fmt"

	"github.com/grokify/gocharts/v2/apps/tableau"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/time/timeutil"
)

func main() {
	filename := "path/to/monthly-items-trend.xlsx"
	outfileHTML := "chart_apps.html"
	outfilePNG := "chart_apps.png"
	err := tableau.WriteFileLineChartCrosstabXLSX(filename, outfileHTML, 0600, timeutil.IntervalMonth, "Monthly Accounts")
	logutil.FatalErr(err)
	fmt.Printf("Wrote (%s)\n", outfileHTML)

	err = tableau.WriteFileLineChartWchartXLSX(filename, outfilePNG, timeutil.IntervalMonth)
	logutil.FatalErr(err)
	fmt.Printf("Wrote (%s)\n", outfilePNG)

	fmt.Println("DONE")
}
