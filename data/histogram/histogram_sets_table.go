package histogram

import (
	"strconv"
	"strings"

	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/type/stringsutil"

	"github.com/grokify/gocharts/v2/data/table"
)

func (hsets *HistogramSets) Table(tableName, colNameHSet, colNameHist, colNameBinName, colNameBinCount string) table.Table {
	tbl := table.NewTable(tableName)
	tbl.Columns = []string{
		stringsutil.FirstNonEmpty(colNameHSet, "Histogram Set"),
		stringsutil.FirstNonEmpty(colNameHist, "Histogram"),
		stringsutil.FirstNonEmpty(colNameBinName, "Bin Name"),
		stringsutil.FirstNonEmpty(colNameBinCount, "Bin Count")}
	hsets.Visit(func(hsetName, histName, binName string, binCount int) {
		tbl.Rows = append(tbl.Rows, []string{
			hsetName, histName, binName, strconv.Itoa(binCount)})
	})
	return tbl
}

type TablePivotOpts struct {
	TableName           string
	ColNameHistogramSet string
	ColNameHistogram    string
	ColNameBinPrefix    string
	ColNameBinSuffix    string
	ColNameBinCountsSum string
	BinNamesOrder       []string
	InclBinsUnordered   bool
	InclBinCounts       bool
	InclBinCountsSum    bool
	InclBinPercentages  bool
}

func (opts TablePivotOpts) ColNameHistogramSetOrDefault() string {
	return stringsutil.FirstNonEmpty(opts.ColNameHistogramSet, "Histogram Set")
}

func (opts TablePivotOpts) ColNameHistogramOrDefault() string {
	return stringsutil.FirstNonEmpty(opts.ColNameHistogramSet, "Histogram")
}

func (opts TablePivotOpts) ColNameBinCountsSumOrDefault() string {
	return stringsutil.FirstNonEmpty(opts.ColNameHistogramSet, "Total")
}

func (opts TablePivotOpts) InflateBinName(binName string, binNumber int, isPct bool) string {
	if strings.TrimSpace(binName) == "" {
		binName = "Unnamed Bin"
		if binNumber >= 0 {
			binName += " " + strconv.Itoa(binNumber)
		}
	}
	return opts.ColNameBinPrefix + binName + opts.ColNameBinSuffix
}

func (opts TablePivotOpts) TableColumns(binNames []string) ([]string, map[int]string) {
	cols := []string{
		opts.ColNameHistogramSetOrDefault(),
		opts.ColNameHistogramOrDefault()}
	fmtMap := map[int]string{}
	if opts.InclBinCounts {
		for i, binName := range binNames {
			cols = append(cols, opts.InflateBinName(binName, i+1, false))
		}
	}
	if opts.InclBinCountsSum {
		cols = append(cols, opts.ColNameBinCountsSumOrDefault())
	}
	if opts.InclBinPercentages {
		for i, binName := range binNames {
			cols = append(cols, opts.InflateBinName(binName, i+1, true))
			fmtMap[len(cols)-1] = table.FormatFloat
		}
	}

	return cols, fmtMap
}

// TablePivot returns a `*table.Table` where the first column is the histogram
// set name, the second column is the histogram name and the other columns are
// the bin names.
func (hsets *HistogramSets) TablePivot(opts TablePivotOpts) table.Table {
	// func (hsets *HistogramSets) TablePivot(tableName, colNameHSet, colNameHist, colNameBinNamePrefix, colNameBinNameSuffix string, binNamesOrder []string, binsInclUnordered bool, inclPercentages bool) table.Table {
	tbl := table.NewTable(opts.TableName)
	tbl.FormatMap = map[int]string{
		-1: table.FormatInt,
		0:  table.FormatString,
		1:  table.FormatString}
	binNames := table.Columns(hsets.BinNames())
	if len(opts.BinNamesOrder) > 0 {
		binNamesOrdered, _ := stringsutil.SliceOrderExplicit(binNames, opts.BinNamesOrder, opts.InclBinsUnordered)
		binNames = binNamesOrdered
	}

	tblCols, fmtMap := opts.TableColumns(binNames)
	tbl.Columns = tblCols
	for k, v := range fmtMap {
		tbl.FormatMap[k] = v
	}

	for hsetName, hset := range hsets.HistogramSetMap {
		for histName, hist := range hset.HistogramMap {
			row := []string{hsetName, histName}
			histSum := 0
			for _, binName := range binNames {
				if binCount, ok := hist.Bins[binName]; ok {
					if opts.InclBinCounts {
						row = append(row, strconv.Itoa(binCount))
						histSum += binCount
					}
				} else {
					if opts.InclBinCounts {
						row = append(row, "0")
					}
				}
			}
			if opts.InclBinCountsSum {
				row = append(row, strconv.Itoa(histSum))
			}
			if opts.InclBinPercentages {
				for _, binName := range binNames {
					if histSum == 0 {
						row = append(row, "0")
					} else if binCount, ok := hist.Bins[binName]; ok && binCount != 0 {
						row = append(row, strconvutil.Ftoa(float64(binCount)/float64(histSum), -1))
					} else {
						row = append(row, "0")
					}
				}
			}
			tbl.Rows = append(tbl.Rows, row)
		}
	}
	return tbl
}

func (hsets *HistogramSets) WriteXLSX(filename, sheetname, colNameHSet, colNameHist, colNameBinName, colNameBinCount string) error {
	tbl := hsets.Table(sheetname, colNameHSet, colNameHist, colNameBinName, colNameBinCount)
	return tbl.WriteXLSX(filename, sheetname)
}

/*
func (hsets *HistogramSets) WriteXLSXPivot(filename, sheetname, colNameHSet, colNameHist, colNameBinNamePrefix, colNameBinNameSuffix string, binNamesOrder []string, binsInclUnordered bool) error {
	tbl := hsets.TablePivot(sheetname, colNameHSet, colNameHist, colNameBinNamePrefix, colNameBinNameSuffix, binNamesOrder, binsInclUnordered)
	return tbl.WriteXLSX(filename, sheetname)
}
*/

func (hsets *HistogramSets) WriteXLSXPivot(filename string, opts TablePivotOpts) error {
	// func (hsets *HistogramSets) WriteXLSXPivot(filename, sheetname, colNameHSet, colNameHist, colNameBinNamePrefix, colNameBinNameSuffix string, binNamesOrder []string, binsInclUnordered bool) error {
	//tbl := hsets.TablePivot(sheetname, colNameHSet, colNameHist, colNameBinNamePrefix, colNameBinNameSuffix, binNamesOrder, binsInclUnordered)
	tbl := hsets.TablePivot(opts)
	return tbl.WriteXLSX(filename, opts.TableName)
}
