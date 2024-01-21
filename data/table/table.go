// table provides a struct to handle tabular data.
package table

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/grokify/mogo/errors/errorsutil"
)

// Table is useful for working on CSV data. It stores records as `[]string` with typed
// formatting information per-column to facilitate transformations.
type Table struct {
	Name                string
	Columns             Columns
	Rows                [][]string
	RowsFloat64         [][]float64
	IsFloat64           bool
	FormatMap           map[int]string
	FormatFunc          func(val string, colIdx uint) (any, error)
	FormatAutoLink      bool
	BackgroundColorFunc func(colIdx, rowIdx uint) string
	ID                  string
	Class               string
	Style               string
}

// NewTable returns a new empty `Table` struct with slices and maps set to empty (non-nil) values.
func NewTable(name string) Table {
	return Table{
		Name:      name,
		Columns:   []string{},
		Rows:      [][]string{},
		FormatMap: map[int]string{}}
}

// LoadMergedRows is used to load data including both column names and rows from `[][]string` sources
// like `csv.ReadAll()`.
func (tbl *Table) LoadMergedRows(data [][]string) {
	if len(data) == 0 {
		return
	}
	tbl.Columns = data[0]
	if len(data) > 1 {
		tbl.Rows = data[1:]
	}
}

func (tbl *Table) UpsertRowColumnValue(rowIdx, colIdx uint, value string) {
	rowIdxInt := int(rowIdx)
	colIdxInt := int(colIdx)
	for rowIdxInt < len(tbl.Rows)-1 {
		tbl.Rows = append(tbl.Rows, []string{})
	}
	row := tbl.Rows[rowIdxInt]
	for colIdxInt < len(row)-1 {
		row = append(row, "")
	}
	row[colIdxInt] = value
	tbl.Rows[rowIdxInt] = row
}

// IsWellFormed returns true when the number of columns equals
// the length of each row. If columns is empty, the length of the
// first row is used for comparison.
func (tbl *Table) IsWellFormed() (isWellFormed bool, columnCount int, mismatchRows []int) {
	isWellFormed = true
	columnCount = len(tbl.Columns)
	if !tbl.IsFloat64 {
		if len(tbl.Rows) == 0 {
			return isWellFormed, columnCount, []int{}
		}
		for i, row := range tbl.Rows {
			if i == 0 && columnCount == 0 {
				columnCount = len(row)
			} else if len(row) != columnCount {
				isWellFormed = false
				mismatchRows = append(mismatchRows, i)
			}
		}
	} else {
		if len(tbl.RowsFloat64) == 0 {
			return isWellFormed, columnCount, []int{}
		}
		for i, row := range tbl.RowsFloat64 {
			if i == 0 && columnCount == 0 {
				columnCount = len(row)
			} else if len(row) != columnCount {
				isWellFormed = false
				mismatchRows = append(mismatchRows, i)
			}
		}
	}
	return
}

// BuildFloat64 populates `RowsFloat64` from `Rows`. It is an experimental feature for machine learning.
// Most features use `Rows`.
func (tbl *Table) BuildFloat64(skipEmpty bool) error {
	for i, row := range tbl.Rows {
		if skipEmpty && len(row) == 0 {
			continue
		}
		var rowFloat64 []float64
		for j, v := range row {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return errorsutil.Wrapf(err, "cannot parse float on row [%d] col [%d] with value [%s]", i, j, v)
			}
			rowFloat64 = append(rowFloat64, f)
		}
		tbl.RowsFloat64 = append(tbl.RowsFloat64, rowFloat64)
	}
	tbl.IsFloat64 = true
	return nil
}

func (tbl *Table) WriteXLSX(path, sheetname string) error {
	tbl.Name = sheetname
	return WriteXLSX(path, []*Table{tbl})
}

func (tbl *Table) WriteCSV(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if len(tbl.Columns) > 0 {
		err = writer.Write(tbl.Columns)
		if err != nil {
			return err
		}
	}
	err = writer.WriteAll(tbl.Rows)
	if err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}

func (tbl *Table) ToSliceMSS() []map[string]string {
	slice := []map[string]string{}
	for _, row := range tbl.Rows {
		slice = append(slice, tbl.Columns.RowMap(row, false))
	}
	return slice
}
