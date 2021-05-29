package frequency

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/simplego/type/stringsutil"
)

// ParseFileCSV reads a CSV using default settings of
// `,` separator, header row and BOM to be stripped. If you
// have other configurations, use `table.ReadFile()` directly
// and call `HistogramFromTable()`.
func ParseFileCSV(file string, name string, binNameColIdx, binFrequencyColIdx uint) (*Histogram, error) {
	tbl, err := table.ReadFile(file, ',', true, true)
	if err != nil {
		return nil, err
	}
	tbl.Name = name
	return ParseTable(tbl, binNameColIdx, binFrequencyColIdx)
}

// ParseTable parses a `table.Table` to a `Histogram` given a table,
// binName column index and binFrequency column index. Empty rows are
// skipped.
func ParseTable(tbl table.Table, binNameColIdx, binFrequencyColIdx uint) (*Histogram, error) {
	hist := NewHistogram(tbl.Name)
	for _, rec := range tbl.Records {
		if stringsutil.SliceIsEmpty(rec, true) {
			continue
		}
		if int(binNameColIdx) >= len(rec) {
			return hist, fmt.Errorf("error row length smaller than binNameColIdx: recordLen[%d] binNameColIdx [%d]",
				len(rec), binNameColIdx)
		} else if int(binFrequencyColIdx) >= len(rec) {
			return hist, fmt.Errorf("error row length smaller than binFrequencyColIdx: recordLen[%d] binFrequencyColIdx [%d]",
				len(rec), binFrequencyColIdx)
		}
		binName := strings.TrimSpace(rec[binNameColIdx])
		binFreq := strings.TrimSpace(rec[binFrequencyColIdx])
		if len(binName) == 0 && len(binFreq) == 0 {
			continue
		}
		if len(binFreq) == 0 {
			hist.Add(binName, 0)
		} else {
			binFreqInt, err := strconv.Atoi(binFreq)
			if err != nil {
				return hist, fmt.Errorf("error strconv frequency string[%s] err[%s]", binFreq, err.Error())
			}
			hist.Add(binName, binFreqInt)
		}
	}
	hist.Inflate()
	return hist, nil
}
