package sheet

import (
	"github.com/xuri/excelize/v2"
)

func SetRowValues(f *excelize.File, sheetName string, rowIndex uint, rowValues []interface{}) error {
	for colIdx, cellValue := range rowValues {
		cellLocation := CoordinatesToSheetLocation(uint32(colIdx), uint32(rowIndex))
		err := f.SetCellValue(sheetName, cellLocation, cellValue)
		if err != nil {
			return err
		}
	}
	return nil
}
