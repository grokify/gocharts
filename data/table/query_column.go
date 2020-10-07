package table

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func (t *Table) ColumnsValuesDistinct(wantCols []string, stripSpace bool) (map[string]int, error) {
	data := map[string]int{}
	wantIdxs := []int{}
	maxIdx := -1
	for _, wantCol := range wantCols {
		wantIdx := t.ColumnIndex(wantCol)
		if wantIdx < 0 {
			return data, fmt.Errorf("Column Not Found [%v]", wantCol)
		}
		wantIdxs = append(wantIdxs, wantIdx)
		if wantIdx > maxIdx {
			maxIdx = wantIdx
		}
	}
	for _, rec := range t.Records {
		if len(rec) > maxIdx {
			vals := []string{}
			for _, wantIdx := range wantIdxs {
				val := rec[wantIdx]
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

func (tbl *Table) ColumnValues(colIdx uint) ([]string, error) {
	idx := int(colIdx)

	vals := []string{}
	for _, row := range tbl.Records {
		if idx < len(row) {
			vals = append(vals, row[colIdx])
		} else {
			return vals, fmt.Errorf("E_COL_IDX [%d] ROW_LEN [%d]", colIdx, len(row))
		}
	}
	return vals, nil
}

func (tbl *Table) ColumnValuesDistinct(colIdx uint) (map[string]int, error) {
	data := map[string]int{}
	idx := int(colIdx)

	for _, rec := range tbl.Records {
		if len(rec) > idx {
			val := rec[idx]
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
		return "", "", errors.New("No Values Found")
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
	for _, row := range tbl.Records {
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

func (tbl *Table) columnIndexMore(colIdx int, colName string) (int, error) {
	if colIdx >= 0 {
		return colIdx, nil
	}
	if len(colName) == 0 {
		return colIdx, errors.New("Must supply colIndex or colName")
	}
	colIdx = tbl.ColumnIndex(colName)
	if colIdx < 0 {
		return colIdx, fmt.Errorf("Column Not Found [%v]", colName)
	}
	return colIdx, nil
}
