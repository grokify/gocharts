package table

import (
	"strconv"
	"strings"
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

func (tbl *Table) SumColFloat64(colIdx int) (float64, error) {
	sum := 0.0
	for _, row := range tbl.Records {
		if colIdx < 0 || colIdx >= len(row) {
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
