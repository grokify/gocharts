package excelizeutil

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/grokify/gocharts/data/table"
)

func SetRowValues(f *excelize.File, sheetName string, rowIndex uint, rowValues []interface{}) {
	for colIdx, cellValue := range rowValues {
		cellLocation := table.CoordinatesToSheetLocation(uint32(colIdx), uint32(rowIndex))
		f.SetCellValue(sheetName, cellLocation, cellValue)
	}
}
