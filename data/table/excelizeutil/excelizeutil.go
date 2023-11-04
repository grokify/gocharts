package excelizeutil

import (
	"errors"
	"strings"

	"github.com/grokify/gocharts/v2/data/table/sheet"
	excelize "github.com/xuri/excelize/v2"
)

var ErrExcelizeFileCannotBeNil = errors.New("excelize.File cannot be nil")

type File struct {
	*excelize.File
}

func NewFile(filename string) (*File, error) {
	f := &File{}
	if filename == "" {
		f.File = excelize.NewFile()
	} else {
		xf, err := excelize.OpenFile(filename)
		if err != nil {
			return nil, err
		}
		f.File = xf
	}
	return f, nil
}

func (f *File) Close() error {
	if f.File == nil {
		return nil
	}
	return f.File.Close()
}

func (f *File) SheetList() []string {
	if f.File == nil {
		return []string{}
	} else {
		return f.File.GetSheetList()
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

func (f *File) TableData(sheetName string, headerRowCount uint, trimSpace bool) ([]string, [][]string, error) {
	cols := []string{}
	rows := [][]string{}

	if f.File == nil {
		return cols, rows, ErrExcelizeFileCannotBeNil
	}
	exCols, err := f.File.GetCols(sheetName)
	if err != nil {
		return cols, rows, err
	}
	exRows, err := f.File.GetRows(sheetName)
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
