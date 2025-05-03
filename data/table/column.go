package table

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/v2/data/table/sheet"
	"github.com/grokify/mogo/type/slicesutil"
)

// Columns represents a slice of string with tabular functions.
type Columns []string

// Index returns the column index of the requested column name. A value of `-1` is returned if the coliumn name is not found.
func (cols Columns) Index(colName string) int {
	for idx, colNameTry := range cols {
		if colNameTry == colName {
			return idx
		}
	}
	return -1
}

// Equal returns true if the number of elements or the element values of the Columns do not match.
func (cols Columns) Equal(c Columns) bool {
	if len(cols) != len(c) {
		return false
	}
	for i, v := range cols {
		if v != c[i] {
			return false
		}
	}
	return true
}

// Unique returns true if the Column names are unqiue, false if there are duplicates.
func (cols Columns) Unique() bool {
	return slicesutil.Unique(cols)
}

// MustCellString returns a single row value or empty string if the column name doesn't exist.
func (cols Columns) MustCellString(colName string, row []string) string {
	if val, err := cols.CellString(colName, row, true, ""); err != nil {
		return ""
	} else {
		return val
	}
}

// CellString returns a single row value.
func (cols Columns) CellString(colName string, row []string, defaultIfEmpty bool, def string) (string, error) {
	for colIdx, colNameTry := range cols {
		if colNameTry == colName {
			if colIdx < len(row) {
				v := row[colIdx]
				if v == "" && defaultIfEmpty {
					return def, nil
				}
				return row[colIdx], nil
			} else {
				return "", errors.New("column index not present in row")
			}
		}
	}
	return "", fmt.Errorf("columnName (%s) not found", colName)
}

// CellFloat64 returns a single row value.
func (cols Columns) CellFloat64(colName string, row []string, defaultIfEmpty bool, def float64) (float64, error) {
	if val, err := cols.CellString(colName, row, false, ""); err != nil {
		return 0, err
	} else if strings.TrimSpace(val) == "" && defaultIfEmpty {
		return def, nil
	} else {
		return strconv.ParseFloat(val, 64)
	}
}

// CellInt returns a single row value.
func (cols Columns) CellInt(colName string, row []string, defaultIfEmpty bool, def int) (int, error) {
	if val, err := cols.CellString(colName, row, false, ""); err != nil {
		return 0, err
	} else if strings.TrimSpace(val) == "" && defaultIfEmpty {
		return def, nil
	} else {
		return strconv.Atoi(val)
	}
}

// CellUint returns a single row value.
func (cols Columns) CellUint(colName string, row []string, defaultIfEmpty bool, def uint) (uint, error) {
	if val, err := cols.CellString(colName, row, false, ""); err != nil {
		return 0, err
	} else if strings.TrimSpace(val) == "" && defaultIfEmpty {
		return def, nil
	} else if ival, err := strconv.Atoi(val); err != nil {
		return 0, err
	} else if ival < 0 {
		return 0, errors.New("cannot convert to `uint` as `int` value is less than zero")
	} else {
		return uint(ival), nil
	}
}

// CellTime returns a single row value. If no `timeFormat` is provided `time.RFC3339` is used.
func (cols Columns) CellTime(colName, timeFormat string, row []string, defaultIfEmpty bool, def time.Time) (time.Time, error) {
	if val, err := cols.CellString(colName, row, false, ""); err != nil {
		return time.Now(), err
	} else if strings.TrimSpace(val) == "" && defaultIfEmpty {
		return time.Time{}, nil
	} else {
		if strings.TrimSpace(timeFormat) == "" {
			timeFormat = time.RFC3339
		}
		return time.Parse(timeFormat, val)
	}
}

func (cols Columns) MustCellTime(colName, timeFormat string, row []string) *time.Time {
	if val, err := cols.CellString(colName, row, false, ""); err != nil {
		return nil
	} else if strings.TrimSpace(val) == "" {
		return nil
	} else if strings.TrimSpace(timeFormat) == "" {
		return nil
	} else if t, err := time.Parse(timeFormat, val); err != nil {
		return nil
	} else {
		return &t
	}
}

func (cols Columns) MustCellIntOrDefault(colName string, row []string, def int) int {
	if val, err := cols.CellString(colName, row, true, ""); err != nil {
		return def
	} else if valI, err := strconv.Atoi(val); err != nil {
		return def
	} else {
		return valI
	}
}

// MustCellsString returns a slice of values.
func (cols Columns) MustCellsString(colNames []string, row []string) []string {
	vals := []string{}
	for _, colName := range colNames {
		vals = append(vals, cols.MustCellString(colName, row))
	}
	return vals
}

// CellsString returns a slice of values.
func (cols Columns) CellsString(colNames []string, row []string, useDefault bool, def string) ([]string, error) {
	missingColumnNames := []string{}
	vals := []string{}
	for _, colName := range colNames {
		if val, err := cols.CellString(colName, row, useDefault, def); err != nil {
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

func (cols Columns) ModifyNamesColLettersMap(m map[string]string, unique bool) (Columns, error) {
	if out, err := sheet.SliceReplaceValueAtLettersMap(cols, m); err != nil {
		return Columns{}, err
	} else if outcols := Columns(out); unique && !outcols.Unique() {
		return Columns{}, errors.New("non unique column names")
	} else {
		return outcols, nil
	}
}
