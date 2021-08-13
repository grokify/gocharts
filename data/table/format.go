package table

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"

	"github.com/grokify/simplego/math/mathutil"
)

// Pivot takes a "straight table" where the columnn names
// and values are in a single column and lays it out as a standard tabular data.
func (tbl *Table) Pivot(colCount uint, haveColumns bool) (Table, error) {
	newTbl := NewTable(tbl.Name)
	if len(tbl.Columns) != 0 {
		return newTbl, fmt.Errorf("has defined columns count [%d]", len(tbl.Columns))
	}
	isWellFormed, colCountActual := tbl.IsWellFormed()
	if !isWellFormed {
		return newTbl, errors.New("table is not well-defined")
	} else if colCountActual != 1 {
		return newTbl, fmt.Errorf("has non-1 column count [%d]", colCountActual)
	}
	rowCount := len(tbl.Rows)
	_, remainder := mathutil.DivideInt64(int64(rowCount), int64(colCount))
	if remainder != 0 {
		return newTbl, fmt.Errorf("row count [%d] is not a multiple of col count [%d]", rowCount, colCount)
	}
	addedColumns := false
	newRow := []string{}
	for i, row := range tbl.Rows {
		_, remainder := mathutil.DivideInt64(int64(i), int64(colCount))
		if remainder == 0 {
			if len(newRow) > 0 {
				if haveColumns && !addedColumns {
					newTbl.Columns = newRow
					addedColumns = true
				} else {
					newTbl.Rows = append(newTbl.Rows, newRow)
				}
				newRow = []string{}
			}
		}
		newRow = append(newRow, row[0])
	}
	if len(newRow) > 0 {
		if haveColumns && !addedColumns {
			newTbl.Columns = newRow
		} else {
			newTbl.Rows = append(newTbl.Rows, newRow)
		}
	}
	return newTbl, nil
}

// FormatColumn takes a function to format all cell values.
func (tbl *Table) FormatColumn(colIdx uint, conv func(cellVal string) (string, error)) error {
	colInt := int(colIdx)
	for i, row := range tbl.Rows {
		if colInt >= len(row) {
			return fmt.Errorf("row [%d] is len [%d] without col index [%d]", i, len(row), colInt)
		}
		newVal, err := conv(row[colInt])
		if err != nil {
			return err
		}
		tbl.Rows[i][colInt] = newVal
	}
	return nil
}

// String writes the table out to a CSV string.
func (tbl *Table) String(comma rune, useCRLF bool) (string, error) {
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	w.Comma = comma
	w.UseCRLF = useCRLF

	if len(tbl.Columns) > 0 {
		if err := w.Write(tbl.Columns); err != nil {
			return "", fmt.Errorf("error writing columns to csv [%s]",
				strings.Join(tbl.Columns, ","))
		}
	}

	for _, row := range tbl.Rows {
		if err := w.Write(row); err != nil {
			return "", fmt.Errorf("error writing row to csv [%s]",
				strings.Join(row, ","))
		}
	}

	w.Flush()
	return b.String(), w.Error()
}
