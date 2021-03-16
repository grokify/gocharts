package table

import (
	"fmt"
)

const StyleSimple = "border:1px solid #000;border-collapse:collapse"

var DebugReadCSV = false // should not need to use this.

// Table is useful for working on CSV data
type Table struct {
	Name           string
	Columns        []string
	Records        [][]string
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
		Records:   [][]string{},
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
		tbl.Records = data[1:]
	}
}

func (t *Table) UpsertRowColumnValue(rowIdx, colIdx uint, value string) {
	rowIdxInt := int(rowIdx)
	colIdxInt := int(colIdx)
	for rowIdxInt < len(t.Records)-1 {
		t.Records = append(t.Records, []string{})
	}
	row := t.Records[rowIdxInt]
	for colIdxInt < len(row)-1 {
		row = append(row, "")
	}
	row[colIdxInt] = value
	t.Records[rowIdxInt] = row
}

func (t *Table) RecordValue(wantCol string, record []string) (string, error) {
	idx := t.ColumnIndex(wantCol)
	if idx < 0 {
		return "", fmt.Errorf("Column Not Found [%v]", wantCol)
	}
	if idx >= len(record) {
		return "", fmt.Errorf("Record does not have enough columns [%v]", idx+1)
	}
	return record[idx], nil
}

func (t *Table) RecordValueOrEmpty(wantCol string, record []string) string {
	val, err := t.RecordValue(wantCol, record)
	if err != nil {
		return ""
	}
	return val
}

func (t *Table) IsWellFormed() (isWellFormed bool, columnCount uint) {
	columnCount = uint(len(t.Columns))
	if len(t.Records) == 0 {
		isWellFormed = true
		return
	}
	for i, rec := range t.Records {
		if i == 0 && len(t.Columns) == 0 {
			columnCount = uint(len(rec))
			continue
		}
		if uint(len(rec)) != columnCount {
			isWellFormed = false
			return
		}
	}
	isWellFormed = true
	return
}

func (t *Table) WriteXLSX(path, sheetname string) error {
	t.Name = sheetname
	return WriteXLSX(path, t)
}

func (t *Table) WriteCSV(path string) error {
	return WriteCSV(path, t)
}

func (t *Table) RecordToMSS(record []string) map[string]string {
	mss := map[string]string{}
	for i, key := range t.Columns {
		if i < len(t.Columns) {
			mss[key] = record[i]
		}
	}
	return mss
}

func (t *Table) ToSliceMSS() []map[string]string {
	slice := []map[string]string{}
	for _, rec := range t.Records {
		slice = append(slice, t.RecordToMSS(rec))
	}
	return slice
}
