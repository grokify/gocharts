package histogram

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gocharts/v2/data/table/excelizeutil"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/type/slicesutil"
	excelize "github.com/xuri/excelize/v2"
)

func (hset *HistogramSet) Table(colNameHist, colNameBin, colNameCount string) table.Table {
	tbl := table.NewTable(hset.Name)
	tbl.Columns = []string{colNameHist, colNameBin, colNameCount}
	tbl.FormatMap = map[int]string{2: table.FormatInt}
	for hName, h := range hset.HistogramMap {
		if h == nil {
			continue
		}
		for bName, count := range h.Bins {
			tbl.Rows = append(tbl.Rows, []string{hName, bName, strconv.Itoa(count)})
		}
	}
	return tbl
}

// WriteXLSXPivot creates an XLSX file where the first column is the
// histogram name and the other columns are the bin names. This is
// useful for easy visualization of a table and also creating
// charts such as grouped bar charts.
func (hset *HistogramSet) WriteXLSXPivot(filename, sheetName, histColName string, addColumnTotalLeft, addColumnTotalRight, addRowTotalTop, addRowTotalBottom bool) error {
	if tbl, err := hset.TablePivot(sheetName, histColName, addColumnTotalLeft, addColumnTotalRight, addRowTotalTop, addRowTotalBottom); err != nil {
		return err
	} else {
		return tbl.WriteXLSX(filename, sheetName)
	}
}

// TablePivot returns a `*table.Table` where the first column is the
// histogram name and the other columns are the bin names. This is
// useful for easy visualization of a table and also creating
// charts such as grouped bar charts.
func (hset *HistogramSet) TablePivot(tableName, histColName string, addColumnTotalLeft, addColumnTotalRight, addRowTotalTop, addRowTotalBottom bool) (*table.Table, error) {
	if len(strings.TrimSpace(tableName)) == 0 {
		tableName = strings.TrimSpace(hset.Name)
	}
	tbl := table.NewTable(tableName)

	if len(strings.TrimSpace(histColName)) == 0 {
		histColName = "Histogram Name"
	}

	binNames := hset.BinNames()
	tbl.Columns = append(tbl.Columns, histColName)
	if addColumnTotalLeft {
		tbl.Columns = append(tbl.Columns, "Total")
	}
	tbl.Columns = append(tbl.Columns, binNames...)
	if addColumnTotalRight {
		tbl.Columns = append(tbl.Columns, "Total")
	}
	tbl.FormatMap = map[int]string{
		-1: table.FormatInt}
	if hset.KeyIsTime {
		tbl.FormatMap[0] = table.FormatTime
	} else {
		tbl.FormatMap[0] = table.FormatString
	}

	hnames := hset.ItemNames()
	colSumRows := [][]int{}
	for _, hname := range hnames {
		row := []string{hname}
		hist, ok := hset.HistogramMap[hname]
		if !ok {
			return nil, fmt.Errorf("histogram name present without histogram [%s]", hname)
		}
		rowTotal := 0
		rowBinCounts := []int{}
		if addColumnTotalLeft {
			rowBinCounts = append(rowBinCounts, 0)
		}
		// rowBinCountVals := []string{}
		for _, binName := range binNames {
			if binVal, ok := hist.Bins[binName]; ok {
				rowTotal += binVal
				rowBinCounts = append(rowBinCounts, binVal)
				// rowBinCountVals = append(rowBinCountVals, strconv.Itoa(binVal))
			} else {
				rowBinCounts = append(rowBinCounts, 0)
				// rowBinCountVals = append(rowBinCountVals, "0")
			}
		}
		rowBinCountTotal := slicesutil.SliceIntSum(rowBinCounts)
		if addColumnTotalLeft {
			//row = append(row, strconv.Itoa(rowTotal))
			rowBinCounts[0] = rowBinCountTotal
		}
		rowBinCountStrs := strconvutil.SliceItoa(rowBinCounts)
		row = append(row, rowBinCountStrs...)
		if addColumnTotalRight {
			row = append(row, strconv.Itoa(rowTotal))
			rowBinCounts = append(rowBinCounts, rowBinCountTotal)
		}
		tbl.Rows = append(tbl.Rows, row)
		colSumRows = append(colSumRows, rowBinCounts)
	}
	if addRowTotalBottom {
		row := []string{"Total"}
		sums := slicesutil.MatrixIntColSums(colSumRows)
		for _, sum := range sums {
			row = append(row, strconv.Itoa(sum))
		}
		tbl.Rows = append(tbl.Rows, row)
	}

	return &tbl, nil
}

// WriteXLSX creates an XLSX file where the first column is the
// histogram name, the second column is the bin name and the
// third column is the bin count.
func (hset *HistogramSet) WriteXLSX(filename, sheetName, colName1, colName2, colNameCount string) error {
	// WriteXLSX writes a table as an Excel XLSX file with
	// row formatter option.
	f := excelize.NewFile()
	// Create a new sheet.

	if len(strings.TrimSpace(sheetName)) == 0 {
		sheetName = strings.TrimSpace(hset.Name)
	}
	if len(sheetName) == 0 {
		sheetName = "Sheet0"
	}
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return errorsutil.Wrap(err, "excelize.File.NewSheet()")
	}

	colName1 = strings.TrimSpace(colName1)
	if len(colName1) == 0 {
		colName1 = hset.Name
	}
	if len(colName1) == 0 {
		colName1 = "Column1"
	}
	colName2 = strings.TrimSpace(colName2)
	if len(colName1) == 0 {
		for _, fstats := range hset.HistogramMap {
			fstats.Name = strings.TrimSpace(fstats.Name)
			if len(fstats.Name) > 0 {
				colName2 = fstats.Name
				break
			}
		}
	}
	colNameCount = strings.TrimSpace(colNameCount)
	if len(colNameCount) == 0 {
		colNameCount = "Count"
	}
	header := []any{colName1, colName2, colNameCount}

	err = excelizeutil.SetRowValues(f, sheetName, 0, header)
	if err != nil {
		return err
	}
	rowIdx := uint(1)
	for fstatsName, fstats := range hset.HistogramMap {
		fstatsNameDt := time.Now()
		if hset.KeyIsTime {
			fstatsNameDt, err = time.Parse(time.RFC3339, fstatsName)
			if err != nil {
				return err
			}
		}
		for binName, binCount := range fstats.Bins {
			var rowVals []any
			if hset.KeyIsTime {
				rowVals = []any{fstatsNameDt, binName, binCount}
			} else {
				rowVals = []any{fstatsName, binName, binCount}
			}
			err := excelizeutil.SetRowValues(f, sheetName, rowIdx, rowVals)
			if err != nil {
				return err
			}
			rowIdx++
		}
	}
	f.SetActiveSheet(index)
	// Delete Original Sheet
	err = f.DeleteSheet(f.GetSheetName(0))
	if err != nil {
		return errorsutil.Wrap(err, "excelize.File.DeleteSheet()")
	}

	// Save xlsx file by the given path.
	return f.SaveAs(filename)
}
