package excelizeutil

import (
	"errors"
	"strings"

	"github.com/grokify/gocharts/v2/data/table/sheet"
	excelize "github.com/xuri/excelize/v2"
)

var ErrExcelizeFileCannotBeNil = errors.New("excelize.File cannot be nil")

type ExcelizeMore struct {
	File *excelize.File
}

func NewExcelizeMore(filename string) (*ExcelizeMore, error) {
	x := &ExcelizeMore{}
	if filename != "" {
		xf, err := excelize.OpenFile(filename)
		if err != nil {
			return nil, err
		}
		x.File = xf
	} else {
		x.File = excelize.NewFile()
	}
	return x, nil
}

func (x *ExcelizeMore) Close() error {
	if x.File == nil {
		return nil
	}
	return x.File.Close()
}

func (x *ExcelizeMore) SheetList() []string {
	if x.File == nil {
		return []string{}
	} else {
		return x.File.GetSheetList()
	}
}

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

func (x *ExcelizeMore) TableData(sheetName string, headerRowCount uint, trimSpace bool) ([]string, [][]string, error) {
	cols := []string{}
	rows := [][]string{}

	if x.File == nil {
		return cols, rows, ErrExcelizeFileCannotBeNil
	}
	exCols, err := x.File.GetCols(sheetName)
	if err != nil {
		return cols, rows, err
	}
	exRows, err := x.File.GetRows(sheetName)
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
