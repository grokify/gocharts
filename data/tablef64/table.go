// table provides a struct to handle tabular data.
package tablef64

import (
	"errors"

	"github.com/grokify/gocharts/v2/data/point"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/type/stringsutil"
)

// Table is useful for working on CSV data. It stores records as `[]string` with typed
// formatting information per-column to facilitate transformations.
type Table struct {
	Name    string
	Columns table.Columns
	Rows    [][]float64
}

type Columns []string

// NewTable returns a new empty `Table` struct with slices and maps set to empty (non-nil) values.
func NewTable(name string) *Table {
	return &Table{
		Name:    name,
		Columns: []string{},
		Rows:    [][]float64{},
	}
}

func ReadFile(filename string, skipEmptyRows bool) (*Table, error) {
	tbl, err := table.ReadFile(nil, filename)
	if err != nil {
		return nil, err
	}
	return FromTableString(&tbl, skipEmptyRows, []uint{})
}

func FromTableString(tbl *table.Table, skipEmptyRows bool, colIndexes []uint) (*Table, error) {
	if tbl == nil {
		return nil, errors.New("input table cannot be nil")
	} else if len(colIndexes) == 0 {
		return fromTableStringAll(tbl, skipEmptyRows)
	}
	n := NewTable("")
	for _, colIdx := range colIndexes {
		if colIdx >= uint(len(tbl.Columns)) {
			return nil, errors.New("colIdx out of range: >= len")
		} else {
			n.Columns = append(n.Columns, tbl.Columns[colIdx])
		}
	}
	for _, r := range tbl.Rows {
		if skipEmptyRows && len(r) == 0 {
			continue
		}
		r2 := stringsutil.Strings(r)
		r3, err := r2.FilterIndexes(colIndexes)
		if err != nil {
			return nil, err
		}
		rf, err := strconvutil.SliceAtof(r3, 64)
		if err != nil {
			return nil, err
		}
		if skipEmptyRows && len(rf) == 0 {
			continue
		}
		n.Rows = append(n.Rows, rf)
	}
	return n, nil
}

func fromTableStringAll(tbl *table.Table, skipEmptyRows bool) (*Table, error) {
	n := NewTable("")
	n.Columns = append(n.Columns, tbl.Columns...)
	for _, r := range tbl.Rows {
		if skipEmptyRows && len(r) == 0 {
			continue
		}
		tryEmpty := stringsutil.SliceCondenseSpace(r, false, false)
		if len(tryEmpty) == 0 {
			continue
		}
		nr, err := strconvutil.SliceAtof(r, 64)
		if err != nil {
			return nil, err
		}
		n.Rows = append(n.Rows, nr)
	}
	return n, nil
}

func (tbl *Table) PointXYsColumnIndexes(xColIdx, yColIdx int) (point.PointXYs, error) {
	xys := point.PointXYs{}
	if len(tbl.Rows) == 0 {
		return xys, nil
	} else if xColIdx < 0 {
		return xys, errors.New("xColIdx cannot be < 0")
	} else if yColIdx < 0 {
		return xys, errors.New("yColIdx cannot be < 0")
	}
	for _, r := range tbl.Rows {
		if len(r) == 0 {
			continue
		} else if xColIdx >= len(r) {
			return xys, errors.New("cannot get x col val")
		} else if yColIdx >= len(r) {
			return xys, errors.New("cannot get y col val")
		}
		xys = append(xys, point.PointXY{
			X: r[xColIdx],
			Y: r[yColIdx]})
	}
	return xys, nil
}

func (tbl *Table) PointXYsColumnNames(xColName, yColName string) (point.PointXYs, error) {
	xys := point.PointXYs{}
	xColIdx := tbl.Columns.Index(xColName)
	if xColIdx < 0 {
		return xys, errors.New("xColName not found")
	}
	yColIdx := tbl.Columns.Index(yColName)
	if yColIdx < 0 {
		return xys, errors.New("yColName not found")
	}
	return tbl.PointXYsColumnIndexes(xColIdx, yColIdx)
}

func (tbl *Table) ValuesColumnIndex(colIndex uint) ([]float64, error) {
	var vals []float64
	colIndexInt := int(colIndex)
	for _, r := range tbl.Rows {
		if colIndexInt >= len(r) {
			return vals, errors.New("index out of bounds")
		}
		vals = append(vals, r[colIndex])
	}
	return vals, nil
}

func (tbl *Table) ValuesColumnName(colName string) ([]float64, error) {
	if colIdx := tbl.Columns.Index(colName); colIdx < 0 {
		return []float64{}, errors.New("column name not found")
	} else {
		return tbl.ValuesColumnIndex(uint(colIdx))
	}
}
