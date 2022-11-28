package table

import (
	"errors"
	"strconv"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
)

const (
	ColNameNone  = "(none)"
	ColNameCount = "(count)"
)

// ColumnExpandPivot adds columns to the table representing each value in the provided column.
func (tbl *Table) ColumnExpandPivot(colIdx uint, split bool, sep, colNamePrefix, colNameNone, existString, notExistString string, addCounts bool, colNameCounts string) (map[int]string, error) {
	colFormats := map[int]string{}
	if int(colIdx) >= len(tbl.Columns) {
		return colFormats, errors.New("colIdx is too large")
	}
	isWellFormed, _, _ := tbl.IsWellFormed()
	if !isWellFormed {
		return colFormats, errors.New("table is not well formed. Cannot expand")
	}
	newColVals, _, err := tbl.ColumnValuesSplit(colIdx, split, sep, true, true)
	if err != nil {
		return colFormats, err
	}
	colNamePrefix = strings.TrimSpace(colNamePrefix)
	if len(colNamePrefix) == 0 {
		colNamePrefix = tbl.Columns[colIdx]
	}
	for _, v := range newColVals {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			if len(strings.TrimSpace(colNameNone)) > 0 {
				v = colNameNone
			} else {
				v = ColNameNone
			}
		}
		newColName := colNamePrefix + ": " + v
		tbl.Columns = append(tbl.Columns, newColName)
	}
	if addCounts {
		colNameCountsActual := colNameCounts
		if len(strings.TrimSpace(colNameCountsActual)) == 0 {
			colNameCountsActual = ColNameCount
		}
		tbl.Columns = append(tbl.Columns, colNamePrefix+": "+colNameCountsActual)
	}
	rowColVals := 0
	for i, row := range tbl.Rows {
		for _, newColVal := range newColVals {
			exist := false
			if !split {
				row[colIdx] = strings.TrimSpace(row[colIdx])
				if newColVal == row[colIdx] {
					exist = true
				}
				if len(row[colIdx]) > 0 {
					rowColVals = 1
				}
			} else {
				rowVals := stringsutil.SliceCondenseSpace(strings.Split(row[colIdx], sep), true, true)
				if stringsutil.SliceIndex(rowVals, newColVal, false, nil) > -1 {
					exist = true
				}
				rowColVals = len(rowVals)
			}
			if exist {
				row = append(row, existString)
			} else {
				row = append(row, notExistString)
			}
		}
		if addCounts {
			row = append(row, strconv.Itoa(rowColVals))
			colFormats[len(row)-1] = FormatInt
			tbl.FormatMap[len(row)-1] = FormatInt
		}
		tbl.Rows[i] = row
	}
	return colFormats, nil
}
