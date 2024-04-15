package excelizeutil

import (
	"errors"
	"sort"
	"strings"

	"github.com/grokify/gocharts/v2/data/table/sheet"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
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
	} else if xf, err := excelize.OpenFile(filename); err != nil {
		return nil, err
	} else {
		f.File = xf
	}
	return f, nil
}

func ReadFile(filename string) (*File, error) {
	f := &File{}
	if xf, err := excelize.OpenFile(filename); err != nil {
		return nil, err
	} else {
		f.File = xf
		return f, nil
	}
}

func (f *File) Close() error {
	if f.File == nil {
		return nil
	} else {
		return f.File.Close()
	}
}

func (f *File) SheetNames(sortAsc bool) []string {
	if f.File == nil {
		return []string{}
	} else {
		if !sortAsc {
			return f.File.GetSheetList()
		} else {
			colNames := f.File.GetSheetList()
			sort.Strings(colNames)
			return colNames
		}
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

func (f *File) SheetColumnNames(sheetName string, trimSpace bool) ([]string, error) {
	if exCols, err := f.File.GetCols(sheetName); err != nil {
		return []string{}, err
	} else {
		return ColumnsTop(exCols, trimSpace), nil
	}
}

func (f *File) TableDataIndex(sheetIdx uint, headerRowCount uint, trimSpace, umerge bool) ([]string, [][]string, error) {
	if f.File == nil {
		return []string{}, [][]string{}, ErrExcelizeFileCannotBeNil
	} else {
		return f.TableData(
			f.File.GetSheetName(int(sheetIdx)),
			headerRowCount,
			trimSpace,
			umerge)
	}
}

func (f *File) TableData(sheetName string, headerRowCount uint, trimSpace, umerge bool) ([]string, [][]string, error) {
	var cols []string
	var rows [][]string
	var err error
	if f.File == nil {
		return []string{}, [][]string{}, ErrExcelizeFileCannotBeNil
	} else {
		cols, err = f.SheetColumnNames(sheetName, trimSpace)
		if err != nil {
			return []string{}, [][]string{}, ErrExcelizeFileCannotBeNil
		}
		rows, err = f.File.GetRows(sheetName)
		if err != nil {
			return []string{}, [][]string{}, err
		}
		if headerRowCount > 0 && headerRowCount <= uint(len(rows)) {
			rows = rows[headerRowCount:]
		}
		if trimSpace {
			for i, row := range rows {
				for j := range row {
					rows[i][j] = strings.TrimSpace(rows[i][j])
				}
			}
		}
		if !umerge {
			// no need to unmerge
			return cols, rows, nil
		} else if len(rows) == 0 || len(rows) == 1 && len(cols) == len(rows[0]) {
			// no need to unmerge
			return cols, rows, nil
		} else {
			mapLengthCounts := slicesutil.LengthCounts(rows)
			if len(mapLengthCounts) == 1 {
				for l := range mapLengthCounts {
					if l == uint(len(cols)) {
						return cols, rows, nil
					}
					break
				}
			}
		}
	}

	var newRows [][]string
	for rowIdx, row := range rows {
		try := stringsutil.SliceCondenseSpace(row, false, false)
		if len(try) == 0 {
			continue
		}
		l := len(row)
		if l == len(cols) {
			newRows = append(newRows, row)
		} else if l > len(cols) {
			return cols, newRows, errors.New("row longer than cols")
		} else {
			newRow := []string{}
			for colIdx := 0; colIdx < len(cols); colIdx++ {
				if cell, err := f.File.GetCellValue(sheetName,
					sheet.CoordinatesToSheetLocation(
						uint(colIdx),
						uint(rowIdx)+headerRowCount,
					)); err != nil {
					return cols, rows, err
				} else {
					if trimSpace {
						cell = strings.TrimSpace(cell)
					}
					newRow = append(newRow, cell)
				}
			}
			if len(newRow) != len(cols) {
				return cols, rows, errors.New("modified row length mismatch")
			} else {
				newRows = append(newRows, newRow)
			}
		}
	}
	mapLengthCountsNew := slicesutil.LengthCounts(newRows)
	if len(mapLengthCountsNew) != 1 {
		return cols, newRows, errors.New("row mismatch after unmerging")
	}
	return cols, newRows, nil
}

// ColumnsTop converts a response from `excelize.File.GetCols()` to an `[]string` suitable for `Table.Columns`.
func ColumnsTop(cols [][]string, trimSpace bool) []string {
	var collapsed []string
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
