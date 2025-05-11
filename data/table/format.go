package table

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/grokify/mogo/math/mathutil"
	"github.com/grokify/mogo/text/markdown"
	"github.com/olekukonko/tablewriter"
)

func (tbl *Table) Markdown(newline string, escPipe bool) string {
	var md string
	if len(tbl.Columns) > 0 {
		md += markdown.TableRowToMarkdown(tbl.Columns, escPipe) + newline
		if len(tbl.Rows) > 0 {
			md += markdown.TableSeparator(len(tbl.Columns)) + newline
		}
	}
	return md + markdown.TableRowsToMarkdown(tbl.Rows, newline, escPipe, false)
}

func (tbl *Table) WriteMarkdown(filename string, perm os.FileMode, newline string, escPipe bool) error {
	md := tbl.Markdown(newline, escPipe)
	return os.WriteFile(filename, []byte(md), perm)
}

// Pivot takes a "straight table" where the columnn names
// and values are in a single column and lays it out as a standard tabular data.
func (tbl *Table) Pivot(colCount uint32, haveColumns bool) (Table, error) {
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
	_, remainder := mathutil.Divide(int64(rowCount), int64(colCount))
	if remainder != 0 {
		return newTbl, fmt.Errorf("row count [%d] is not a multiple of col count [%d]", rowCount, colCount)
	}
	addedColumns := false
	newRow := []string{}
	for i, row := range tbl.Rows {
		_, remainder := mathutil.Divide(int64(i), int64(colCount))
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
func (tbl *Table) FormatColumn(colIdx uint32, conv func(cellVal string) (string, error), skipRowLengthMismatch bool) error {
	colInt := int(colIdx)
	for i, row := range tbl.Rows {
		if colInt >= len(row) {
			if skipRowLengthMismatch {
				continue
			} else {
				return fmt.Errorf("row [%d] is len [%d] without col index [%d]", i, len(row), colInt)
			}
		} else if newVal, err := conv(row[colInt]); err != nil {
			return err
		} else {
			tbl.Rows[i][colInt] = newVal
		}
	}
	return nil
}

func (tbl *Table) FormatColumns(colIdxMin uint32, colIdxMax int, conv func(cellVal string) (string, error), skipRowLengthMismatch bool) error {
	colIdxMinInt := int(colIdxMin)
	if colIdxMax >= 0 && colIdxMax < colIdxMinInt {
		return errors.New("colIdxMax cannot be less than colIdxMin")
	}
	for i, row := range tbl.Rows {
		if colIdxMinInt >= len(row) {
			if skipRowLengthMismatch {
				continue
			} else {
				return fmt.Errorf("row [%d] is len [%d] without min col index [%d]", i, len(row), colIdxMinInt)
			}
		} else if !skipRowLengthMismatch && colIdxMax >= 0 && colIdxMax >= len(row) {
			return fmt.Errorf("row [%d] is len [%d] without max col index [%d]", i, len(row), colIdxMax)
		}

		colIdxMaxRow := colIdxMax
		if colIdxMaxRow < 0 {
			colIdxMaxRow = len(row) - 1
		}
		for j := colIdxMinInt; j <= colIdxMaxRow; j++ {
			if j >= len(row) {
				if skipRowLengthMismatch {
					continue
				} else {
					return fmt.Errorf("row [%d] is len [%d] without col index [%d] on max col index [%d]", i, len(row), j, colIdxMax)
				}
			} else if newVal, err := conv(row[j]); err != nil {
				return err
			} else {
				tbl.Rows[i][j] = newVal
			}
		}
	}
	return nil
}

// FormatRows formats row cells using a start and ending column index and a convert function.
// The `format.ConvertDecommify()` and `format.ConvertRemoveControls()` functions are available to use.
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
	// test and return errors
	for y, row := range tbl.Rows {
		if colIdxMinInc >= len(row) {
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

// Text renders the table in text suitable for console reporting.
func (tbl *Table) Text(w io.Writer) error {
	if w == nil {
		return errors.New("writer must be supplied")
	}
	tw := tablewriter.NewWriter(w)
	tw.Header(slices.Clone(tbl.Columns))
	for _, r := range tbl.Rows {
		if err := tw.Append(slices.Clone(r)); err != nil {
			return err
		}
	}
	return tw.Render()
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
