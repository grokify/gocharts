package table

type Columns []string

func (cols Columns) RowVal(colName string, row []string) string {
	for colIdx, colNameTry := range cols {
		if colNameTry == colName {
			if colIdx < len(row) {
				return row[colIdx]
			}
		}
	}
	return ""
}

func (cols Columns) RowVals(colNames []string, row []string) []string {
	vals := []string{}
	for _, colName := range colNames {
		vals = append(vals, cols.RowVal(colName, row))
	}
	return vals
}
