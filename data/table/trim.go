package table

func (tbl *Table) TrimColumnsRight() {
	for {
		if !tbl.trimColumnRight() {
			break
		}
	}
}

func (tbl *Table) trimColumnRight() bool {
	if len(tbl.Columns) == 0 {
		return false
	}
	colIdxMax := len(tbl.Columns) - 1
	if tbl.Columns[colIdxMax] != "" {
		return false
	}
	for _, row := range tbl.Rows {
		rowIdxMax := len(row) - 1
		if colIdxMax != rowIdxMax {
			return false
		}
		if row[rowIdxMax] != "" {
			return false
		}
	}
	tbl.Columns = tbl.Columns[:len(tbl.Columns)-1]
	for i, row := range tbl.Rows {
		tbl.Rows[i] = row[:len(row)-1]
	}
	return true
}
