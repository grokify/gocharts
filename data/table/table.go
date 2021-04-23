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

func (tbl *Table) UpsertRowColumnValue(rowIdx, colIdx uint, value string) {
	rowIdxInt := int(rowIdx)
	colIdxInt := int(colIdx)
	for rowIdxInt < len(tbl.Records)-1 {
		tbl.Records = append(tbl.Records, []string{})
	}
	row := tbl.Records[rowIdxInt]
	for colIdxInt < len(row)-1 {
		row = append(row, "")
	}
	row[colIdxInt] = value
	tbl.Records[rowIdxInt] = row
}

func (tbl *Table) RecordValue(wantCol string, record []string) (string, error) {
	idx := tbl.ColumnIndex(wantCol)
	if idx < 0 {
		return "", fmt.Errorf("Column Not Found [%v]", wantCol)
	}
	if idx >= len(record) {
		return "", fmt.Errorf("Record does not have enough columns [%v]", idx+1)
	}
	return record[idx], nil
}

func (tbl *Table) RecordValueOrEmpty(wantCol string, record []string) string {
	val, err := tbl.RecordValue(wantCol, record)
	if err != nil {
		return ""
	}
	return val
}

func (tbl *Table) IsWellFormed() (isWellFormed bool, columnCount uint) {
	columnCount = uint(len(tbl.Columns))
	if len(tbl.Records) == 0 {
		isWellFormed = true
		return
	}
	for i, rec := range tbl.Records {
		if i == 0 && len(tbl.Columns) == 0 {
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

func (tbl *Table) WriteXLSX(path, sheetname string) error {
	tbl.Name = sheetname
	return WriteXLSX(path, tbl)
}

func (tbl *Table) WriteCSV(path string) error {
	return WriteCSV(path, tbl)
}

func (tbl *Table) RecordToMSS(record []string) map[string]string {
	mss := map[string]string{}
	for i, key := range tbl.Columns {
		if i < len(tbl.Columns) {
			mss[key] = record[i]
		}
	}
	return mss
}

func (tbl *Table) ToSliceMSS() []map[string]string {
	slice := []map[string]string{}
	for _, rec := range tbl.Records {
		slice = append(slice, tbl.RecordToMSS(rec))
	}
	return slice
}
