package table

// RowCellCounts returns a `map[int]int` where the
// key is the cell count and the value is the number
// of rows.
func (tbl *Table) RowCellCounts() map[int]int {
	mii := map[int]int{}
	for _, row := range tbl.Rows {
		mii[len(row)] += 1
	}
	return mii
}
