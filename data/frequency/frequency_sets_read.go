package frequency

import (
	"fmt"
	"strings"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/simplego/encoding/jsonutil"
	"github.com/grokify/simplego/math/mathutil"
	"github.com/grokify/simplego/type/stringsutil"
)

// NewFrequencySetsCSVs expects multiple files to have same columns.
func NewFrequencySetsCSVs(filenames []string, key1ColIdx, key2ColIdx, uidColIdx uint) (FrequencySets, table.Table, error) {
	fsets := NewFrequencySets()
	tbl, err := table.NewTableFilesSimple(filenames, ",", true, true)
	if err != nil {
		return fsets, tbl, err
	}
	fsets, err = NewFrequencySetsTable(tbl, key1ColIdx, key2ColIdx, uidColIdx)
	return fsets, tbl, err
}

func NewFrequencySetsTable(tbl table.Table, key1ColIdx, key2ColIdx, uidColIdx uint) (FrequencySets, error) {
	fsets := NewFrequencySets()
	_, maxIdx := mathutil.MinMaxUint(key1ColIdx, key2ColIdx, uidColIdx)
	for _, row := range tbl.Records {
		if len(stringsutil.SliceCondenseSpace(row, true, false)) == 0 {
			continue
		} else if len(row) <= int(maxIdx) {
			return fsets, fmt.Errorf(
				"NewFrequencySetsTable.E_ROW_LEN_ERROR NEED_ROW_LEN [%v] HAVE_ROW_LEN [%v] ROW_DATA [%s]",
				maxIdx+1, len(row),
				jsonutil.MustMarshalSimple(row, "", ""))
		}
		fsets.Add(
			strings.TrimSpace(row[key1ColIdx]),
			strings.TrimSpace(row[key2ColIdx]),
			strings.TrimSpace(row[uidColIdx]), true)
	}
	return fsets, nil
}
