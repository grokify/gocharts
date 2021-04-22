package table

import (
	"fmt"
)

func (tbl *Table) NewTableFilterColDistinctFirst(colIdx int) *Table {
	newTbl := NewTable()
	newTbl.Columns = tbl.Columns

	seen := map[string]int{}
	for _, row := range tbl.Records {
		if colIdx >= 0 && colIdx < len(row) {
			val := row[colIdx]
			if _, ok := seen[val]; !ok {
				newTbl.Records = append(newTbl.Records, row)
				seen[val] = 1
			}
		}
	}
	return &newTbl
}

// NewTableFilterColumnValues returns a Table filtered
// by column names and column values.
func (tbl *Table) NewTableFilterColumnValues(wantColNameValues map[string]string) (Table, error) {
	t2 := Table{Columns: tbl.Columns}
	records, err := tbl.FilterRecordsColumnValues(wantColNameValues)
	if err != nil {
		return t2, err
	}
	t2.Records = records
	return t2, nil
}

// FilterRecordsColumnValues returns a set of records filtered
// by column names and column values.
func (tbl *Table) FilterRecordsColumnValues(wantColNameValues map[string]string) ([][]string, error) {
	data := [][]string{}
	wantColIndexes := map[string]int{}
	maxIdx := -1
	for wantColName := range wantColNameValues {
		wantColIdx := tbl.ColumnIndex(wantColName)
		if wantColIdx < 0 {
			return data, fmt.Errorf("Column Not Found [%v]", wantColName)
		}
		if wantColIdx > maxIdx {
			maxIdx = wantColIdx
		}
		wantColIndexes[wantColName] = wantColIdx
	}
RECORDS:
	for _, rec := range tbl.Records {
		if len(rec) > maxIdx {
			for wantColName, wantColIdx := range wantColIndexes {
				colValue := rec[wantColIdx]
				wantColValue, ok := wantColNameValues[wantColName]
				if !ok {
					return data, fmt.Errorf("Column Name [%v] has no desired value", wantColName)
				}
				if colValue != wantColValue {
					continue RECORDS
				}
			}
			data = append(data, rec)
		}
	}
	return data, nil
}
