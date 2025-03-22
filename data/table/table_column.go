package table

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/grokify/mogo/type/slicesutil"
)

func (tbl *Table) RowsModify(fn func(i int, row []string) ([]string, error)) error {
	for i, row := range tbl.Rows {
		if try, err := fn(i, row); err != nil {
			return err
		} else {
			tbl.Rows[i] = try
		}
	}
	return nil
}

/*
func (tbl *Table) addColumnLineNumberClone(colName string, startNumber int) (*Table, error) {
	out := tbl.Clone(true)
	// Update Columns
	if colName == "" {
		colName = "Number"
	}
	out.Columns = []string{colName}
	out.Columns = append(out.Columns, tbl.Columns...)
	// Update Format Map
	out.FormatMap = map[int]string{0: FormatInt}
	for k, v := range tbl.FormatMap {
		if k >= 0 {
			out.FormatMap[k+1] = v
		} else {
			out.FormatMap[k] = v
		}
	}
	// Update Rows
	err := out.RowsModify(func(i int, row []string) ([]string, error) {
		outRow := []string{strconv.Itoa(i + startNumber)}
		outRow = append(outRow, row...)
		return slices.Clone(outRow), nil
	})
	return out, err
}
*/

func (tbl *Table) AddColumnLineNumber(colName string, startNumber int) {
	// Update Columns
	if colName == "" {
		colName = "Number"
	}
	newCols := []string{colName}
	newCols = append(newCols, tbl.Columns...)
	tbl.Columns = slices.Clone(newCols)
	// Update Format Map
	newFmtMap := map[int]string{0: FormatInt}
	for k, v := range tbl.FormatMap {
		if k >= 0 {
			newFmtMap[k+1] = v
		} else {
			newFmtMap[k] = v
		}
	}
	tbl.FormatMap = newFmtMap
	// Update Rows
	err := tbl.RowsModify(func(i int, row []string) ([]string, error) {
		outRow := []string{strconv.Itoa(i + startNumber)}
		outRow = append(outRow, row...)
		return slices.Clone(outRow), nil
	})
	if err != nil {
		panic(err)
	}
}

func (tbl *Table) ColumnsValuesDistinct(wantColNames []string, stripSpace bool) (map[string]int, error) {
	data := map[string]int{}
	if len(wantColNames) == 0 {
		return data, nil
	}
	wantIdxs := []int{}
	maxIdx := -1
	for _, wantCol := range wantColNames {
		wantIdx := tbl.Columns.Index(wantCol)
		if wantIdx < 0 {
			return data, fmt.Errorf("column not found [%v]", wantCol)
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

func (tbl *Table) ColumnValuesCountsByName(colName string, trimSpace, includeEmpty, lowerCase bool) (map[string]int, error) {
	colIdx := tbl.Columns.Index(colName)
	if colIdx <= 0 {
		return map[string]int{}, errors.New("column name not found")
	}
	return tbl.ColumnValuesCounts(colIdx, trimSpace, includeEmpty, lowerCase), nil
}

func (tbl *Table) ColumnValuesCounts(colIdx int, trimSpace, includeEmpty, lowerCase bool) map[string]int {
	m := map[string]int{}
	if colIdx < 0 {
		return m
	}
	for _, row := range tbl.Rows {
		if colIdx >= len(row) {
			continue
		}
		v := row[colIdx]
		if trimSpace {
			v = strings.TrimSpace(v)
		}
		if !includeEmpty && v == "" {
			continue
		}
		if lowerCase {
			v = strings.ToLower(v)
		}
		m[v]++
	}
	return m
}

func (tbl *Table) ColumnValuesSplit(colIdx uint32, split bool, sep string, unique, sortResults bool) ([]string, map[string]int, error) {
	msi := map[string]int{}
	vals := []string{}
	for _, row := range tbl.Rows {
		if int(colIdx) < len(row) {
			v := row[colIdx]
			if split {
				valsi := strings.Split(v, sep)
				vals = append(vals, valsi...)
				for _, v := range valsi {
					msi[v] += 1
				}
			} else {
				vals = append(vals, v)
				msi[v] += 1
			}
		} else {
			return vals, msi, fmt.Errorf("column index not found for index [%d] row length [%d]", colIdx, len(row))
		}
	}
	if unique {
		vals = slicesutil.Dedupe(vals)
	}
	if sortResults {
		sort.Strings(vals)
	}
	return vals, msi, nil
}

func (tbl *Table) ColumnValues(colIdx int, unique, sortResults bool) ([]string, error) {
	vals := []string{}
	if colIdx < 0 {
		return vals, fmt.Errorf("colIdx cannot be negative (%d)", colIdx)
	}
	seen := map[string]int{}
	for _, row := range tbl.Rows {
		if colIdx < len(row) {
			if unique {
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

func (tbl *Table) ColumnValuesName(colName string, unique, sortResults bool) ([]string, error) {
	idx := tbl.Columns.Index(colName)
	if idx < 0 {
		return []string{}, fmt.Errorf("column name not found (%s)", colName)
	}
	return tbl.ColumnValues(idx, unique, sortResults)
}

func (tbl *Table) ColumnValuesForColumnName(colName string, dedupeValues, sortValues bool) ([]string, error) {
	colIdx := tbl.Columns.Index(colName)
	if colIdx <= 0 {
		return []string{}, fmt.Errorf("column [%s] not found", colName)
	}
	return tbl.ColumnValues(colIdx, dedupeValues, sortValues)
}

func (tbl *Table) columnValuesDistinct(colIdx uint32) map[string]int {
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
	return data
}

func (tbl *Table) ColumnValuesMinMax(colIdx uint32) (string, string, error) {
	vals := tbl.columnValuesDistinct(colIdx)
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

func (tbl *Table) ColumnSumFloat64(colIdx uint32) (float64, error) {
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
