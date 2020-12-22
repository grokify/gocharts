package table

import (
	"errors"
	"fmt"

	"github.com/grokify/gotilla/math/mathutil"
)

// FormatStraightToTabular takes a "straight table" wheere the columnn names
// and values are in a single column and lays it out as a standard tabular data.
func FormatStraightToTabular(tbl Table, colCount uint) (Table, error) {
	newTbl := NewTable()
	if len(tbl.Columns) != 0 {
		return newTbl, fmt.Errorf("Has Defined Columns Count [%d]", len(tbl.Columns))
	}
	isWellFormed, colCountActual := tbl.IsWellFormed()
	if !isWellFormed {
		return newTbl, errors.New("table is not well-defined")
	} else if colCountActual != 1 {
		return newTbl, fmt.Errorf("has non-1 column count [%d]", colCountActual)
	}
	rowCount := len(tbl.Records)
	_, remainder := mathutil.DivideInt64(int64(rowCount), int64(colCount))
	if remainder != 0 {
		return newTbl, fmt.Errorf("row count [%d] is not a multiple of col count [%d]", rowCount, colCount)
	}
	newRow := []string{}
	for i, rec := range tbl.Records {
		_, remainder := mathutil.DivideInt64(int64(i), int64(colCount))
		if remainder == 0 {
			if len(newRow) > 0 {
				newTbl.Records = append(newTbl.Records, newRow)
				newRow = []string{}
			}
		}
		newRow = append(newRow, rec[0])
	}
	if len(newRow) > 0 {
		newTbl.Records = append(newTbl.Records, newRow)
	}
	return newTbl, nil
}
