package histogram

import (
	"fmt"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/simplego/encoding/jsonutil"
	"github.com/grokify/simplego/math/mathutil"
	"github.com/grokify/simplego/type/stringsutil"
)

// NewHistogramSetsCSVs expects multiple files to have same columns.
func NewHistogramSetsCSVs(filenames []string, key1ColIdx, key2ColIdx, uidColIdx uint) (*HistogramSets, table.Table, error) {
	hsets := NewHistogramSets("")
	tbl, err := table.ReadFiles(nil, filenames...)
	if err != nil {
		return hsets, tbl, err
	}
	hsets, err = NewHistogramSetsTable(tbl, key1ColIdx, key2ColIdx, uidColIdx)
	return hsets, tbl, err
}

func NewHistogramSetsTable(tbl table.Table, key1ColIdx, key2ColIdx, uidColIdx uint) (*HistogramSets, error) {
	hsets := NewHistogramSets(tbl.Name)
	_, maxIdx := mathutil.MinMaxUint(key1ColIdx, key2ColIdx, uidColIdx)
	for _, row := range tbl.Rows {
		if len(stringsutil.SliceCondenseSpace(row, true, false)) == 0 {
			continue
		} else if len(row) <= int(maxIdx) {
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
