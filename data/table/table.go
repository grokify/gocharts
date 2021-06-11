package table

const StyleSimple = "border:1px solid #000;border-collapse:collapse"

// Table is useful for working on CSV data
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

func NewTable() Table {
	return Table{
		Columns:   []string{},
		Rows:      [][]string{},
		FormatMap: map[int]string{}}
}

// LoadMergedRows is used to load data from `[][]string` sources
// like csv.ReadAll()
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

func (tbl *Table) IsWellFormed() (isWellFormed bool, columnCount uint) {
	columnCount = uint(len(tbl.Columns))
	if len(tbl.Rows) == 0 {
		isWellFormed = true
		return
	}
	for i, row := range tbl.Rows {
		if i == 0 && len(tbl.Columns) == 0 {
			columnCount = uint(len(row))
			continue
		}
		if uint(len(row)) != columnCount {
			isWellFormed = false
			return
		}
	}
	isWellFormed = true
	return
}

func (tbl *Table) WriteXLSX(path, sheetname string) error {
	tbl.Name = sheetname
	return WriteXLSX(path, tbl)
}

func (tbl *Table) WriteCSV(path string) error {
	return WriteCSV(path, tbl)
}

func (tbl *Table) ToSliceMSS() []map[string]string {
	slice := []map[string]string{}
	for _, row := range tbl.Rows {
		slice = append(slice, tbl.Columns.RowMap(row, false))
	}
	return slice
}
