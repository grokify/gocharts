package table

import (
	"sort"
	"strings"
)

type TableSet struct {
	Name       string
	Columns    []string
	FormatMap  map[int]string
	FormatFunc func(val string, colIdx uint) (interface{}, error)
	TableMap   map[string]*Table
	Order      []string
}

func NewTableSet(name string) TableSet {
	return TableSet{
		Name:      name,
		Columns:   []string{},
		FormatMap: map[int]string{},
		TableMap:  map[string]*Table{}}
}

func (ts *TableSet) TableNames() []string {
	names := []string{}
	for name := range ts.TableMap {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (ts *TableSet) TablesSorted() []*Table {
	return ts.Tables(ts.TableNames())
}

func (ts *TableSet) Tables(names []string) []*Table {
	tbls := []*Table{}
	for _, name := range names {
		if tbl, ok := ts.TableMap[name]; ok {
			tbls = append(tbls, tbl)
		}
	}
	return tbls
}

func (ts *TableSet) AddRow(tableName string, row []string) {
	tableName = strings.TrimSpace(tableName)
	tbl, ok := ts.TableMap[tableName]
	if !ok {
		tbl := NewTable(tableName)
		tbl.Columns = ts.Columns
		tbl.FormatMap = ts.FormatMap
		tbl.FormatFunc = ts.FormatFunc
	}
	tbl.Rows = append(tbl.Rows, row)
	ts.TableMap[tableName] = tbl
}

func (ts *TableSet) WriteXLSX(filename string) error {
	names := ts.Order
	if len(names) == 0 {
		names = ts.TableNames()
	}
	tbls := ts.Tables(names)
	return WriteXLSX(filename, tbls...)
}
