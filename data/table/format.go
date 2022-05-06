package table

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"

	"github.com/grokify/mogo/math/mathutil"
)

// Pivot takes a "straight table" where the columnn names
// and values are in a single column and lays it out as a standard tabular data.
func (tbl *Table) Pivot(colCount uint, haveColumns bool) (Table, error) {
	newTbl := NewTable(tbl.Name)
	if len(tbl.Columns) != 0 {
		return newTbl, fmt.Errorf("has defined columns count [%d]", len(tbl.Columns))
	}
	isWellFormed, colCountActual, _ := tbl.IsWellFormed()
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

/*
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
*/

// FormatRows formats row cells using a start and ending column index and a convert function.
// The `format.ConvertDecommify()` and `format.FormatStringRemoveControls()` functions are available to use.
func (tbl *Table) FormatRows(colIdxMinInc, colIdxMaxInc int, conv func(cellVal string) (string, error)) error {
	err := tbl.formatRowsTry(colIdxMinInc, colIdxMaxInc, conv, false)
	if err != nil {
		return err
	}
	return tbl.formatRowsTry(colIdxMinInc, colIdxMaxInc, conv, true)
}

func (tbl *Table) formatRowsTry(colIdxMinInc, colIdxMaxInc int, conv func(cellVal string) (string, error), exec bool) error {
	if len(tbl.Rows) == 0 {
		return nil
	}
	if colIdxMinInc < 0 {
		colIdxMinInc = 0
	}
	//testand return errors
	for y, row := range tbl.Rows {
		if int(colIdxMinInc) >= len(row) {
			continue
		}
		rowMaxIdxInc := colIdxMaxInc
		if rowMaxIdxInc < 0 || rowMaxIdxInc >= len(row) {
			rowMaxIdxInc = len(row) - 1
		}
		for x := colIdxMinInc; x <= rowMaxIdxInc; x++ {
			val, err := conv(row[x])
			if err != nil {
				return err
			}
			if exec {
				row[x] = val
			}
		}
		if exec {
			tbl.Rows[y] = row
		}
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

	for i, row := range tbl.Rows {
		if err := w.Write(row); err != nil {
			return "", fmt.Errorf("error writing row to csv: idx [%d] content [%s]",
				i, strings.Join(row, ","))
		}
	}

	w.Flush()
	return b.String(), w.Error()
}

// Transpose creates a new table by transposing the matrix data.
// In the new table, it does not set anything other than than `Name`, `Columns`, and `Rows`.
func (tbl *Table) Transpose() (Table, error) {
	tbl2 := NewTable(tbl.Name)
	isWellFormed, _, _ := tbl.IsWellFormed()
	if !isWellFormed {
		return tbl2, errors.New("can only transpose well formed table")
	}
	for x := 0; x < len(tbl.Columns); x++ {
		newRow := []string{}
		if len(tbl.Columns) > 0 {
			newRow = append(newRow, tbl.Columns[x])
		}
		for y := 0; y < len(tbl.Rows); y++ {
			newRow = append(newRow, tbl.Rows[y][x])
		}
		if x == 0 {
			tbl2.Columns = newRow
		} else {
			tbl2.Rows = append(tbl2.Rows, newRow)
		}
	}
	return tbl2, nil
}
