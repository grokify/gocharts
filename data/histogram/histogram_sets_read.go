package histogram

import (
	"fmt"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/math/mathutil"
	"github.com/grokify/mogo/type/stringsutil"

	"github.com/grokify/gocharts/v2/data/table"
)

// NewHistogramSetsCSVs expects multiple files to have same columns.
func NewHistogramSetsCSVs(filenames []string, key1ColIdx, key2ColIdx, uidColIdx uint32) (*HistogramSets, table.Table, error) {
	hsets := NewHistogramSets("")
	tbl, err := table.ReadFile(nil, filenames...)
	if err != nil {
		return hsets, tbl, err
	}
	hsets, err = NewHistogramSetsTable(tbl, key1ColIdx, key2ColIdx, uidColIdx)
	return hsets, tbl, err
}

func NewHistogramSetsTable(tbl table.Table, key1ColIdx, key2ColIdx, uidColIdx uint32) (*HistogramSets, error) {
	hsets := NewHistogramSets(tbl.Name)
	_, maxIdx := mathutil.MinMaxUint(uint(key1ColIdx), uint(key2ColIdx), uint(uidColIdx))
	for _, row := range tbl.Rows {
		if len(stringsutil.SliceCondenseSpace(row, true, false)) == 0 {
			continue
		} else if uint(len(row)) <= maxIdx {
			return hsets, fmt.Errorf(
				"NewHistogramSetsTable.E_ROW_LEN_ERROR NEED_ROW_LEN [%v] HAVE_ROW_LEN [%v] ROW_DATA [%s]",
				maxIdx+1, len(row),
				jsonutil.MustMarshalSimple(row, "", ""))
		}
		hsets.Add(
			row[key1ColIdx],
			row[key2ColIdx],
			row[uidColIdx], 1, true)
	}
	return hsets, nil
}
