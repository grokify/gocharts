package frequency

import (
	"fmt"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gotilla/encoding/csvutil"
	"github.com/grokify/gotilla/encoding/jsonutil"
	"github.com/grokify/gotilla/math/mathutil"
	"github.com/grokify/gotilla/type/stringsutil"
)

func NewFrequencySetsCSV(filename string, key1ColIdx, key2ColIdx, uidColIdx uint) (FrequencySets, error) {
	fsets := NewFrequencySets()
	tbl, err := csvutil.NewTableDataFileSimple(filename, ",", true, true)
	if err != nil {
		return fsets, err
	}
	return NewFrequencySetsTable(tbl, key1ColIdx, key2ColIdx, uidColIdx)
}

func NewFrequencySetsTable(tbl table.TableData, key1ColIdx, key2ColIdx, uidColIdx uint) (FrequencySets, error) {
	fsets := NewFrequencySets()
	_, maxIdx := mathutil.MinMaxUint(key1ColIdx, key2ColIdx, uidColIdx)
	for _, row := range tbl.Records {
		if len(stringsutil.SliceCondenseSpace(row, true, false)) == 0 {
			continue
		} else if len(row) <= int(maxIdx) {
			return fsets, fmt.Errorf(
				"E_ROW_LEN_ERROR NEED [%v] HAVE [%v] ROW [%s]",
				maxIdx+1, len(row),
				jsonutil.MustMarshalSimple(row, "", ""))
		}
		val1 := row[key1ColIdx]
		val2 := row[key2ColIdx]
		vuid := row[uidColIdx]
		fsets.Add(val1, val2, vuid, true)
	}
	return fsets, nil
}
