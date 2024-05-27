package tableau

import (
	"os"
	"time"

	"github.com/grokify/gocharts/v2/charts/google"
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
	tss, err := ReadFileCrosstabXLSX(infile, interval)
	if err != nil {
		return err
	}
	lcm := google.NewLineChartMaterial()
	lcm.Title = title
	err = lcm.LoadTimeSeriesSetMonth(tss, func(t time.Time) string {
		return t.Format("Jan 2006")
	})
	if err != nil {
		return err
	}
	return lcm.WriteFilePage(outfile, perm)
}
