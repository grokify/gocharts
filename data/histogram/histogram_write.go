package histogram

import (
	"strconv"

	"github.com/grokify/gocharts/data/table"
)

func (hist *Histogram) ToTable(colNameBinName, colNameBinCount string) *table.Table {
	tbl := table.NewTable()
	tbl.Columns = []string{colNameBinName, colNameBinCount}
	for k, v := range hist.Items {
		tbl.Records = append(tbl.Records,
			[]string{k, strconv.Itoa(v)})
	}
	tbl.FormatMap = map[int]string{1: "int"}
	return &tbl
}

func (hist *Histogram) WriteXLSX(filename, sheetname, colNameBinName, colNameBinCount string) error {
	tbl := hist.ToTable(colNameBinName, colNameBinCount)
	return tbl.WriteXLSX(filename, sheetname)
}
