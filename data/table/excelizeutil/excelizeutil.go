package excelizeutil

import (
	"errors"
	"strings"

	"github.com/grokify/gocharts/v2/data/table/sheet"
	excelize "github.com/xuri/excelize/v2"
)

var ErrExcelizeFileCannotBeNil = errors.New("excelize.File cannot be nil")

func SetRowValues(f *excelize.File, sheetName string, rowIndex uint, rowValues []any) error {
	if f == nil {
		return ErrExcelizeFileCannotBeNil
	}
	for colIdx, cellValue := range rowValues {
		cellLocation := sheet.CoordinatesToSheetLocation(uint(colIdx), rowIndex)
		err := f.SetCellValue(sheetName, cellLocation, cellValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetCellValue(f *excelize.File, sheetName string, colIdx, rowIdx uint, opts ...excelize.Options) (string, error) {
	if f == nil {
		return "", ErrExcelizeFileCannotBeNil
	}
	cellLoc := sheet.CoordinatesToSheetLocation(colIdx, rowIdx)
	return f.GetCellValue(sheetName, cellLoc, opts...)
}

func GetTableData(f *excelize.File, sheetName string, headerRowCount uint, trimSpace bool) ([]string, [][]string, error) {
	cols := []string{}
	rows := [][]string{}
	if f == nil {
		return cols, rows, ErrExcelizeFileCannotBeNil
	}
	exCols, err := f.GetCols(sheetName)
	if err != nil {
		return cols, rows, err
	}
	exRows, err := f.GetRows(sheetName)
	if err != nil {
		return cols, rows, err
	}
	if headerRowCount > 0 && headerRowCount <= uint(len(exRows)) {
		exRows = exRows[headerRowCount:]
	}
	if trimSpace {
		for i, row := range exRows {
			for j := range row {
				exRows[i][j] = strings.TrimSpace(exRows[i][j])
			}
		}
	}
	cols = ColumnsCollapse(exCols, trimSpace)
	return cols, exRows, nil
}

// ColumnsCollapse converts a response from `excelize.File.GetCols()` to an `[]string` suitable for
// `Table.Columns`.
func ColumnsCollapse(cols [][]string, trimSpace bool) []string {
	collapsed := []string{}
	for _, col := range cols {
		if len(col) > 0 {
			colName := col[0]
			if trimSpace {
				colName = strings.TrimSpace(colName)
			}
			collapsed = append(collapsed, colName)
		} else {
			collapsed = append(collapsed, "")
		}
	}
	return collapsed
}
