package histogram

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/mogo/type/stringsutil"
)

// ParseFileCSV reads a CSV using default settings of
// `,` separator, header row and BOM to be stripped. If you
// have other configurations, use `table.ReadFile()` directly
// and call `HistogramFromTable()`.
func ParseFileCSV(file string, name string, binNameColIdx, binFrequencyColIdx uint) (*Histogram, error) {
	tbl, err := table.ReadFile(nil, file)
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
	for _, row := range tbl.Rows {
		if stringsutil.SliceIsEmpty(row, true) {
			continue
		}
		if int(binNameColIdx) >= len(row) {
			return hist, fmt.Errorf("error row length smaller than binNameColIdx: recordLen[%d] binNameColIdx [%d]",
				len(row), binNameColIdx)
		} else if int(binFrequencyColIdx) >= len(row) {
			return hist, fmt.Errorf("error row length smaller than binFrequencyColIdx: recordLen[%d] binFrequencyColIdx [%d]",
				len(row), binFrequencyColIdx)
		}
		binName := strings.TrimSpace(row[binNameColIdx])
		binFreq := strings.TrimSpace(row[binFrequencyColIdx])
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
