package table

import (
	"errors"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
)

// ColumnExpand adds columns to the table representing
// each value in the provided column.
func (tbl *Table) ColumnExpand(colIdx uint, split bool, sep, existString, notExistString string) error {
	if int(colIdx) >= len(tbl.Columns) {
		return errors.New("colIdx is too large")
	}
	isWellFormed, _, _ := tbl.IsWellFormed()
	if !isWellFormed {
		return errors.New("table is not well formed. Cannot expand")
	}
	newColVals, _, err := tbl.ColumnValuesSplit(colIdx, split, sep, true, true)
	if err != nil {
		return err
	}
	colName := tbl.Columns[colIdx]
	for _, v := range newColVals {
		newColName := colName + ": " + v
		tbl.Columns = append(tbl.Columns, newColName)
	}
	for i, row := range tbl.Rows {
		for _, newColVal := range newColVals {
			exist := false
			if !split {
				if newColVal == row[colIdx] {
					exist = true
				}
			} else {
				rowVals := strings.Split(row[colIdx], sep)
				if stringsutil.SliceIndex(rowVals, newColVal) > -1 {
					exist = true
				}
			}
			if exist {
				row = append(row, existString)
			} else {
				row = append(row, notExistString)
			}
		}
		tbl.Rows[i] = row
	}
	return nil
}
