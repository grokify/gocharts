package tableau

import (
	"os"
	"time"

	"github.com/grokify/gocharts/v2/charts/google/linechart"
	"github.com/grokify/gocharts/v2/charts/wchart"
	"github.com/grokify/gocharts/v2/charts/wchart/sts2wchart"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gocharts/v2/data/timeseries"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/time/timeutil"
)

// ReadFileCrosstabXLSX reads a Tableau crosstab XLSX file where there is one sheet
// named "Sheet 1" with Column headings, "empty" followed by "FullMonth Year", e.g.
// "January 2024". The rows are various timeseries with numbers as integers with
// commas for thousands.
func ReadFileCrosstabXLSX(filename string, interval timeutil.Interval) (*timeseries.TimeSeriesSet, error) {
	if tbl, err := table.ReadTableXLSXIndexFile(filename, 0, 1, true); err != nil {
		return nil, err
	} else if tss, err := timeseries.ParseTableTimeSeriesSetMatrixRows(*tbl, timeutil.IntervalMonth, false,
		timeutil.ParseTimeCanonicalFunc("Jan 2006"),
		strconvutil.AtoiMoreFunc(",", "."),
		nil); err != nil {
		return nil, err
	} else {
		tss.Interval = interval
		return tss, nil
	}
}

func WriteFileLineChartCrosstabXLSX(infile, outfile string, perm os.FileMode, interval timeutil.Interval, title string) error {
	if tss, err := ReadFileCrosstabXLSX(infile, interval); err != nil {
		return err
	} else {
		return WriteFileLineChartTimeSeriesSet(tss, outfile, perm, interval, title)
	}
}

func WriteFileLineChartTimeSeriesSet(tss *timeseries.TimeSeriesSet, outfile string, perm os.FileMode, interval timeutil.Interval, title string) error {
	lcm := linechart.NewChart()
	lcm.Title = title
	if err := lcm.LoadTimeSeriesSetMonth(tss, func(t time.Time) string { return t.Format("Jan 2006") }); err != nil {
		return err
	} else {
		return lcm.WriteFilePage(outfile, perm)
	}
}

func WriteFileLineChartWchartXLSX(infile, outfile string, interval timeutil.Interval) error {
	if tss, err := ReadFileCrosstabXLSX(infile, interval); err != nil {
		return err
	} else {
		return WriteFileLineChartWchartTimeSeriesSet(tss, outfile, interval)
	}
}

func WriteFileLineChartWchartTimeSeriesSet(tss *timeseries.TimeSeriesSet, outfile string, interval timeutil.Interval) error {
	opts := sts2wchart.DefaultLineChartOpts()
	opts.Interval = interval
	if chart, err := sts2wchart.TimeSeriesSetToLineChart(*tss, opts); err != nil {
		return err
	} else {
		return wchart.WritePNGFile(outfile, chart)
	}
}
