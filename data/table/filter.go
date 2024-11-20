package table

import (
	"fmt"
)

func (tbl *Table) FilterColumnDistinctFirstTable(colIdx int) *Table {
	out := tbl.Clone(false)

	seen := map[string]int{}
	for _, row := range tbl.Rows {
		if colIdx >= 0 && colIdx < len(row) {
			val := row[colIdx]
			if _, ok := seen[val]; !ok {
				out.Rows = append(out.Rows, row)
				seen[val] = 1
			}
		}
	}
	return out
}

// FilterColumnValuesTable returns a Table filtered by column names and column values.
func (tbl *Table) FilterColumnValuesTable(wantColNameValues map[string]string) (*Table, error) {
	out := tbl.Clone(false)
	rows, err := tbl.FilterColumnValuesRows(wantColNameValues)
	if err != nil {
		return out, err
	}
	out.Rows = rows
	return out, nil
}

// FilterRecordsColumnValues returns a set of records filtered by column names and column values.
func (tbl *Table) FilterColumnValuesRows(wantColNameValues map[string]string) ([][]string, error) {
	data := [][]string{}
	wantColIndexes := map[string]int{}
	maxIdx := -1
	for wantColName := range wantColNameValues {
		wantColIdx := tbl.Columns.Index(wantColName)
		if wantColIdx < 0 {
			return data, fmt.Errorf("column not found [%v]", wantColName)
		}
		if wantColIdx > maxIdx {
			maxIdx = wantColIdx
		}
		wantColIndexes[wantColName] = wantColIdx
	}
ROWS:
	for _, row := range tbl.Rows {
		if len(row) > maxIdx {
			for wantColName, wantColIdx := range wantColIndexes {
				colValue := row[wantColIdx]
				wantColValue, ok := wantColNameValues[wantColName]
				if !ok {
					return data, fmt.Errorf("column name (%s) has no desired value", wantColName)
				}
				if colValue != wantColValue {
					continue ROWS
				}
			}
			data = append(data, row)
		}
	}
	return data, nil
}
