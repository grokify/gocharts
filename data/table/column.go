package table

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Columns represents a slice of string with tabular functions.
type Columns []string

// MustRowVal returns a single row value.
func (cols Columns) MustRowVal(colName string, row []string) string {
	val, err := cols.RowVal(colName, row)
	if err != nil {
		return ""
	}
	return val
}

// RowVal returns a single row value.
func (cols Columns) RowVal(colName string, row []string) (string, error) {
	for colIdx, colNameTry := range cols {
		if colNameTry == colName {
			if colIdx < len(row) {
				return row[colIdx], nil
			}
		}
	}
	return "", fmt.Errorf("columnName [%s] not found", colName)
}

// RowValFloat64 returns a single row value.
func (cols Columns) RowValFloat64(colName string, row []string) (float64, error) {
	val, err := cols.RowVal(colName, row)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(val, 64)
}

// RowValInt returns a single row value.
func (cols Columns) RowValInt(colName string, row []string) (int, error) {
	val, err := cols.RowVal(colName, row)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

// RowValTime returns a single row value.
func (cols Columns) RowValTime(colName string, row []string) (time.Time, error) {
	val, err := cols.RowVal(colName, row)
	if err != nil {
		return time.Now(), err
	}
	return time.Parse(time.RFC3339, val)
}

// MustRowVals returns a slice of values.
func (cols Columns) MustRowVals(colNames []string, row []string) []string {
	vals := []string{}
	for _, colName := range colNames {
		vals = append(vals, cols.MustRowVal(colName, row))
	}
	return vals
}

// RowVals returns a slice of values.
func (cols Columns) RowVals(colNames []string, row []string) ([]string, error) {
	missingColumnNames := []string{}
	vals := []string{}
	for _, colName := range colNames {
		val, err := cols.RowVal(colName, row)
		if err != nil {
			missingColumnNames = append(missingColumnNames, colName)
		} else {
			vals = append(vals, val)
		}
	}
	if len(missingColumnNames) > 0 {
		return vals, fmt.Errorf(
			"columnNames missing [%s]", strings.Join(missingColumnNames, ","))
	}
	return vals, nil
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
