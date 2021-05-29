package table

import (
	"strings"
)

// Columns represents a slice of string with tabular functions.
type Columns []string

// RowVal returns a single row value.
func (cols Columns) RowVal(colName string, row []string) string {
	for colIdx, colNameTry := range cols {
		if colNameTry == colName {
			if colIdx < len(row) {
				return row[colIdx]
			}
		}
	}
	return ""
}

// RowVals returns a slice of values.
func (cols Columns) RowVals(colNames []string, row []string) []string {
	vals := []string{}
	for _, colName := range colNames {
		vals = append(vals, cols.RowVal(colName, row))
	}
	return vals
}

// RowMap converts a CSV row to a `map[string]string`.
func (cols Columns) RowMap(row []string, omitEmpty bool) map[string]string {
	mss := map[string]string{}
	for i, key := range cols {
		if i < len(row) {
			val := strings.TrimSpace(row[i])
			if !omitEmpty || len(val) > 0 {
				mss[key] = row[i]
			}
		} else if !omitEmpty {
			mss[key] = ""
		}
	}
	return mss
}
