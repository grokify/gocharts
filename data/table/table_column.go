package table

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/grokify/mogo/type/slicesutil"
)

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
	return tbl.ColumnValuesCounts(uint(colIdx), trimSpace, includeEmpty, lowerCase), nil
}

func (tbl *Table) ColumnValuesCounts(colIdx uint, trimSpace, includeEmpty, lowerCase bool) map[string]int {
	m := map[string]int{}
	colIdxInt := int(colIdx)
	for _, row := range tbl.Rows {
		if colIdxInt >= len(row) {
			continue
		}
		v := row[colIdxInt]
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

func (tbl *Table) ColumnValuesSplit(colIdx uint, split bool, sep string, unique, sortResults bool) ([]string, map[string]int, error) {
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

func (tbl *Table) ColumnValues(colIdx uint, unique, sortResults bool) ([]string, error) {
	//idx := int(colIdx)
	seen := map[string]int{}
	vals := []string{}
	for _, row := range tbl.Rows {
		if int(colIdx) < len(row) {
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
	return tbl.ColumnValues(uint(idx), unique, sortResults)
}

func (tbl *Table) ColumnValuesForColumnName(colName string, dedupeValues, sortValues bool) ([]string, error) {
	colIdx := tbl.Columns.Index(colName)
	if colIdx <= 0 {
		return []string{}, fmt.Errorf("column [%s] not found", colName)
	}
	return tbl.ColumnValues(uint(colIdx), dedupeValues, sortValues)
}

func (tbl *Table) columnValuesDistinct(colIdx uint) map[string]int {
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

func (tbl *Table) ColumnValuesMinMax(colIdx uint) (string, string, error) {
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
