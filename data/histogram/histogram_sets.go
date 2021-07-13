package histogram

import (
	"sort"
	"strconv"
	"strings"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/simplego/type/stringsutil"
)

type HistogramSets struct {
	Name            string
	HistogramSetMap map[string]*HistogramSet
}

func NewHistogramSets(name string) *HistogramSets {
	return &HistogramSets{
		Name:            name,
		HistogramSetMap: map[string]*HistogramSet{}}
}

func (hsets *HistogramSets) Add(hsetName, histName, binName string, binCount int, trimSpace bool) {
	if trimSpace {
		hsetName = strings.TrimSpace(hsetName)
		histName = strings.TrimSpace(histName)
		binName = strings.TrimSpace(binName)
	}
	hset, ok := hsets.HistogramSetMap[hsetName]
	if !ok {
		hset = NewHistogramSet(hsetName)
	}
	hset.Add(histName, binName, binCount)
	hsets.HistogramSetMap[hsetName] = hset
}

func (hsets *HistogramSets) Flatten(name string) *HistogramSet {
	hsetFlat := NewHistogramSet(name)
	for _, hset := range hsets.HistogramSetMap {
		for histName, hist := range hset.HistogramMap {
			for binName, binCount := range hist.Bins {
				hsetFlat.Add(histName, binName, binCount)
			}
		}
	}
	return hsetFlat
}

func (hsets *HistogramSets) BinNames() []string {
	binNamesMap := map[string]int{}
	hsets.VisitEach(func(hsetName, histName, binName string, binCount int) {
		binNamesMap[binName] = 1
	})
	binNames := []string{}
	for binName := range binNamesMap {
		binNames = append(binNames, binName)
	}
	sort.Strings(binNames)
	return binNames
}

// ItemCount returns the number of histogram sets.
func (hsets *HistogramSets) ItemCount() uint {
	return uint(len(hsets.HistogramSetMap))
}

func (hsets *HistogramSets) Counts() *HistogramSetsCounts {
	return NewHistogramSetsCounts(*hsets)
}

func (hsets *HistogramSets) VisitEach(visit func(hsetName, histName, binName string, binCount int)) {
	for hsetName, hset := range hsets.HistogramSetMap {
		for histName, hist := range hset.HistogramMap {
			for binName, binCount := range hist.Bins {
				visit(hsetName, histName, binName, binCount)
			}
		}
	}
}

func (hsets *HistogramSets) Table(tableName, colNameHSet, colNameHist, colNameBinName, colNameBinCount string) table.Table {
	tbl := table.NewTable(tableName)
	tbl.Columns = []string{
		stringsutil.FirstNonEmpty(colNameHSet, "Histogram Set"),
		stringsutil.FirstNonEmpty(colNameHist, "Histogram"),
		stringsutil.FirstNonEmpty(colNameBinName, "Bin Name"),
		stringsutil.FirstNonEmpty(colNameBinCount, "Bin Count")}
	hsets.VisitEach(func(hsetName, histName, binName string, binCount int) {
		tbl.Rows = append(tbl.Rows, []string{
			hsetName, histName, binName, strconv.Itoa(binCount)})
	})
	return tbl
}

// TableMatrix returns a `*table.Table` where the first column is the histogram
// set name, the second column is the histogram name and the other columns are
// the bin names.
func (hsets *HistogramSets) TableMatrix(tableName, colNameHSet, colNameHist, colNameBinNamePrefix, colNameBinNameSuffix string) table.Table {
	tbl := table.NewTable(tableName)
	tbl.FormatMap = map[int]string{
		-1: table.FormatInt,
		0:  table.FormatString,
		1:  table.FormatString}
	binNames := table.Columns(hsets.BinNames())
	tbl.Columns = []string{
		stringsutil.FirstNonEmpty(colNameHSet, "Histogram Set"),
		stringsutil.FirstNonEmpty(colNameHist, "Histogram")}
	for i, binName := range binNames {
		if len(strings.TrimSpace(binName)) == 0 {
			binName = "Unnamed Bin " + strconv.Itoa(i+1)
		}
		tbl.Columns = append(tbl.Columns, colNameBinNamePrefix+binName+colNameBinNameSuffix)
	}
	for hsetName, hset := range hsets.HistogramSetMap {
		for histName, hist := range hset.HistogramMap {
			row := []string{hsetName, histName}
			for _, binName := range binNames {
				if val, ok := hist.Bins[binName]; ok {
					row = append(row, strconv.Itoa(val))
				} else {
					row = append(row, "0")
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

func (hsets *HistogramSets) WriteXLSXMatrix(filename, sheetname, colNameHSet, colNameHist, colNameBinNamePrefix, colNameBinNameSuffix string) error {
	tbl := hsets.TableMatrix(sheetname, colNameHSet, colNameHist, colNameBinNamePrefix, colNameBinNameSuffix)
	return tbl.WriteXLSX(filename, sheetname)
}
