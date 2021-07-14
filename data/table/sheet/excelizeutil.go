package sheet

import (
	"github.com/360EntSecGroup-Skylar/excelize"
)

func SetRowValues(f *excelize.File, sheetName string, rowIndex uint, rowValues []interface{}) {
	for colIdx, cellValue := range rowValues {
		cellLocation := CoordinatesToSheetLocation(uint32(colIdx), uint32(rowIndex))
		f.SetCellValue(sheetName, cellLocation, cellValue)
	}
}
