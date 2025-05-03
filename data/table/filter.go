package table

import (
	"errors"
	"fmt"
	"slices"
)

func (tbl *Table) FilterColumnDistinctFirstTable(colIdx int) (*Table, error) {
	if tbl.IsFloat64 {
		return nil, errors.New("cannot filter float table on string values")
	}
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
	return out, nil
}

// FilterColumnValuesTable returns a Table filtered by column names and column values.
func (tbl *Table) FilterColumnValuesTable(wantColNameValues map[string][]string) (*Table, error) {
	if tbl.IsFloat64 {
		return nil, errors.New("cannot filter float table on string values")
	}
	out := tbl.Clone(false)
	if rows, err := tbl.FilterColumnValuesRows(wantColNameValues); err != nil {
		return out, err
	} else {
		out.Rows = rows
		return out, nil
	}
}

// FilterRecordsColumnValues returns a set of records filtered by column names and column values.
// The supplied `wantColNameValues` provides a list of column names and a set of values,
// any of which can match the desired rows.
func (tbl *Table) FilterColumnValuesRows(wantColNameValues map[string][]string) ([][]string, error) {
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
				wantColValues, ok := wantColNameValues[wantColName]
				if !ok {
					return data, fmt.Errorf("column name (%s) has no desired value", wantColName)
				}
				if !slices.Contains(wantColValues, colValue) {
					continue ROWS
				}
			}
			data = append(data, row)
		}
	}
	return data, nil
}
