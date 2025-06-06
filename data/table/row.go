package table

import (
	"errors"
)

// RowCellCounts returns a `map[int]int` where the key is the cell count
// and the value is the number of rows.
func (tbl *Table) RowCellCounts() map[uint]uint {
	mii := map[uint]uint{}
	for _, row := range tbl.Rows {
		mii[uint(len(row))] += 1
	}
	return mii
}

// RowsToMap converts two columns to a map, given a `keyIdx` and `valIdx` for the keys and values.
func (tbl *Table) RowsToMap(keyIdx, valIdx uint32) (map[string]string, error) {
	m := map[string]string{}
	kIdxInt := int(keyIdx)
	vIdxInt := int(valIdx)
	for _, r := range tbl.Rows {
		if kIdxInt >= len(r) || vIdxInt >= len(r) {
			return m, errors.New("index not present in row")
		}
		m[r[keyIdx]] = r[valIdx]
	}
	return m, nil
}
