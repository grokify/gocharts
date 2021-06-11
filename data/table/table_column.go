package table

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func (tbl *Table) ColumnsValuesDistinct(wantCols []string, stripSpace bool) (map[string]int, error) {
	data := map[string]int{}
	wantIdxs := []int{}
	maxIdx := -1
	for _, wantCol := range wantCols {
		wantIdx := tbl.Columns.Index(wantCol)
		if wantIdx < 0 {
			return data, fmt.Errorf("Column Not Found [%v]", wantCol)
		}
		wantIdxs = append(wantIdxs, wantIdx)
		if wantIdx > maxIdx {
			maxIdx = wantIdx
		}
	}
	for _, row := range tbl.Rows {
		if len(row) > maxIdx {
			vals := []string{}
			for _, wantIdx := range wantIdxs {
				val := row[wantIdx]
				if stripSpace {
					val = strings.TrimSpace(val)
				}
				vals = append(vals, val)
			}
			valsStr := strings.Join(vals, " ")
			_, ok := data[valsStr]
			if !ok {
				data[valsStr] = 0
			}
			data[valsStr]++
		}
	}
	return data, nil
}

/*
func (tbl *Table) ColumnIndex(colName string) int {
	for i, tryColName := range tbl.Columns {
		if tryColName == colName {
			return i
		}
	}
	return -1
}
*/

func (tbl *Table) ColumnValues(colIdx uint, dedupeValues, sortResults bool) ([]string, error) {
	idx := int(colIdx)

	seen := map[string]int{}
	vals := []string{}
	for _, row := range tbl.Rows {
		if idx < len(row) {
			if dedupeValues {
				if _, ok := seen[row[colIdx]]; ok {
					continue
				}
				seen[row[colIdx]] = 1
			}
			vals = append(vals, row[colIdx])
		} else {
			return vals, fmt.Errorf("column index not found for index [%d] row length [%d]", colIdx, len(row))
		}
	}
	if sortResults {
		sort.Strings(vals)
	}
	return vals, nil
}

func (tbl *Table) ColumnValuesForColumnName(colName string, dedupeValues, sortValues bool) ([]string, error) {
	colIdx := tbl.Columns.Index(colName)
	if colIdx <= 0 {
		return []string{}, fmt.Errorf("column [%s] not found", colName)
	}
	return tbl.ColumnValues(uint(colIdx), dedupeValues, sortValues)
}

func (tbl *Table) ColumnValuesDistinct(colIdx uint) (map[string]int, error) {
	data := map[string]int{}
	idx := int(colIdx)

	for _, row := range tbl.Rows {
		if len(row) > idx {
			val := row[idx]
			_, ok := data[val]
			if !ok {
				data[val] = 0
			}
			data[val]++
		}
	}
	return data, nil
}

func (tbl *Table) ColumnValuesMinMax(colIdx uint) (string, string, error) {
	vals, err := tbl.ColumnValuesDistinct(colIdx)
	if err != nil {
		return "", "", err
	}
	if len(vals) == 0 {
		return "", "", errors.New("no values found")
	}

	arr := []string{}
	for val := range vals {
		arr = append(arr, val)
	}

	sort.Strings(arr)
	return arr[0], arr[len(arr)-1], nil
}

func (tbl *Table) ColumnSumFloat64(colIdx uint) (float64, error) {
	sum := 0.0
	idx := int(colIdx)
	for _, row := range tbl.Rows {
		if idx >= len(row) {
			continue
		}
		vstr := strings.TrimSpace(row[colIdx])
		if len(vstr) == 0 {
			continue
		}
		vnum, err := strconv.ParseFloat(vstr, 64)
		if err != nil {
			return sum, err
		}
		sum += vnum
	}
	return sum, nil
}

/*
func (tbl *Table) columnIndexMore(colIdx int, colName string) (int, error) {
	if colIdx >= 0 {
		return colIdx, nil
	}
	if len(colName) == 0 {
		return colIdx, errors.New("must supply `colIndex` or `colName`")
	}
	colIdx = tbl.Columns.Index(colName)
	if colIdx < 0 {
		return colIdx, fmt.Errorf("columnName not found [%v]", colName)
	}
	return colIdx, nil
}
*/
