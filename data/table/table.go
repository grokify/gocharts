// table provides a struct to handle tabular data.
package table

const (
	FormatFloat  = "float"
	FormatInt    = "int"
	FormatString = "string"
	FormatTime   = "time"
	StyleSimple  = "border:1px solid #000;border-collapse:collapse"
)

// Table is useful for working on CSV data. It stores
// records as `[]string` with typed formatting information
// per-column to facilitate transformations.
type Table struct {
	Name           string
	Columns        Columns
	Rows           [][]string
	FormatMap      map[int]string
	FormatFunc     func(val string, colIdx uint) (interface{}, error)
	FormatAutoLink bool
	ID             string
	Class          string
	Style          string
}

// NewTable returns a new empty `Table` struct with
// slices and maps set to empty (non-nil) values.
func NewTable(name string) Table {
	return Table{
		Name:      name,
		Columns:   []string{},
		Rows:      [][]string{},
		FormatMap: map[int]string{}}
}

// LoadMergedRows is used to load data including both
// column names and rows from `[][]string` sources
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
	if len(tbl.Rows) == 0 {
		return isWellFormed, columnCount, []int{}
	}
	for i, row := range tbl.Rows {
		if i == 0 && columnCount == 0 {
			columnCount = len(row)
			continue
		}
		if len(row) != columnCount {
			isWellFormed = false
			mismatchRows = append(mismatchRows, i)
		}
	}
	return
}

func (tbl *Table) WriteXLSX(path, sheetname string) error {
	tbl.Name = sheetname
	return WriteXLSX(path, tbl)
}

func (tbl *Table) WriteCSV(path string) error {
	return writeCSV(path, tbl)
}

func (tbl *Table) ToSliceMSS() []map[string]string {
	slice := []map[string]string{}
	for _, row := range tbl.Rows {
		slice = append(slice, tbl.Columns.RowMap(row, false))
	}
	return slice
}
