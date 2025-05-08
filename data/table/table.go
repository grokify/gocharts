// table provides a struct to handle tabular data.
package table

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"slices"
	"strconv"

	"github.com/grokify/mogo/encoding/csvutil"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/errors/errorsutil"
)

var ErrTableCannotBeNil = errors.New("table cannot be nil")

// Table is useful for working on CSV data. It stores records as `[]string` with typed
// formatting information per-column to facilitate transformations.
type Table struct {
	Name                string
	Columns             Columns
	Rows                [][]string
	RowsFloat64         [][]float64
	IsFloat64           bool
	FormatMap           map[int]string
	FormatFunc          func(val string, colIdx uint32) (any, error) `json:"-"`
	FormatAutoLink      bool
	BackgroundColorFunc func(colIdx, rowIdx uint) string `json:"-"`
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

func ReadFileCSV(filepath string, sep rune) (*Table, error) {
	cr, f, err := csvutil.NewReaderFile(filepath, sep)
	if err != nil {
		return nil, err
	} else {
		defer f.Close()
	}
	t := NewTable("")
	i := 0
	for {
		row, err := cr.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				return nil, err
			}
		}
		if i == 0 {
			t.Columns = row
		} else {
			t.Rows = append(t.Rows, row)
		}
		i++
	}
	return &t, nil
}

// ReadFileJSON loads a JSON marshal of `Table{}`, such as produced by `Table.WriteJSON()`.
func ReadFileJSON(filepath string) (*Table, error) {
	if b, err := os.ReadFile(filepath); err != nil {
		return nil, err
	} else {
		t := Table{}
		err = json.Unmarshal(b, &t)
		return &t, err
	}
}

func (tbl *Table) Clone(inclRows bool) *Table {
	out := Table{
		Name:        tbl.Name,
		Columns:     slices.Clone(tbl.Columns),
		IsFloat64:   tbl.IsFloat64,
		FormatMap:   map[int]string{},
		ID:          tbl.ID,
		Class:       tbl.Class,
		Style:       tbl.Style,
		Rows:        [][]string{},
		RowsFloat64: [][]float64{},
	}
	for k, v := range tbl.FormatMap {
		out.FormatMap[k] = v
	}
	if inclRows {
		for _, r := range tbl.Rows {
			out.Rows = append(out.Rows, slices.Clone(r))
		}
		for _, r := range tbl.RowsFloat64 {
			out.RowsFloat64 = append(out.RowsFloat64, slices.Clone(r))
		}
	}
	return &out
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

func (tbl *Table) UpsertRowColumnValue(rowIdx, colIdx uint32, value string) {
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

/*
func (tbl *Table) ColumnInsert(name string, index uint, format string) error {


}
*/

func (tbl *Table) ToSliceMSS() []map[string]string {
	slice := []map[string]string{}
	for _, row := range tbl.Rows {
		slice = append(slice, tbl.Columns.RowMap(row, false))
	}
	return slice
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

// WriteJSON writes the raw structure of the Table to JSON so it can be read back with
// metadata intact. It does not cover functions as those cannot be marshaled.
func (tbl *Table) WriteJSON(path string, perm os.FileMode, jsonPrefix, jsonIndent string) error {
	if b, err := jsonutil.MarshalSimple(tbl, jsonPrefix, jsonIndent); err != nil {
		return err
	} else {
		return os.WriteFile(path, b, perm)
	}
}

type jsonRecords struct {
	Records []map[string]string `json:"records,omitempty"`
}

func (tbl *Table) WriteJSONObjects(path string, perm os.FileMode, jsonPrefix, jsonIndent string) error {
	out := jsonRecords{Records: tbl.ToSliceMSS()}
	if b, err := jsonutil.MarshalSimple(out, jsonPrefix, jsonIndent); err != nil {
		return err
	} else {
		return os.WriteFile(path, b, perm)
	}
}

func (tbl *Table) WriteXLSX(path, sheetname string) error {
	tbl.Name = sheetname
	return WriteXLSX(path, []*Table{tbl})
}
