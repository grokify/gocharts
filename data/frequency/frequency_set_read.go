package frequency

import (
	"fmt"
	"strings"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gotilla/encoding/csvutil"
	"github.com/grokify/gotilla/encoding/jsonutil"
	"github.com/grokify/gotilla/math/mathutil"
	"github.com/grokify/gotilla/type/stringsutil"
)

// NewFrequencySetsCSVs expects multiple files to have same
// columns.
func NewFrequencySetsCSVs(filenames []string, key1ColIdx, key2ColIdx, uidColIdx uint) (FrequencySets, table.TableData, error) {
	fsets := NewFrequencySets()
	tbl := table.NewTableData()
	for _, filename := range filenames {
		filename = strings.TrimSpace(filename)
		if len(filename) == 0 {
			continue
		}
		tblx, err := csvutil.NewTableDataFileSimple(filename, ",", true, true)
		if err != nil {
			return fsets, tbl, err
		}
		if len(tbl.Columns) == 0 {
			tbl.Columns = tblx.Columns
		}
		if len(tblx.Records) > 0 {
			tbl.Records = append(tbl.Records, tblx.Records...)
		}
	}
	fsets, err := NewFrequencySetsTable(tbl, key1ColIdx, key2ColIdx, uidColIdx)
	return fsets, tbl, err
}

func NewFrequencySetsTable(tbl table.TableData, key1ColIdx, key2ColIdx, uidColIdx uint) (FrequencySets, error) {
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
