package table

import (
	"sort"
	"strings"
)

type TableSet struct {
	Columns      []string
	RowFormatter func(val string, col uint) (interface{}, error)
	TableMap     map[string]*TableFormatter
}

func (ts *TableSet) TableNames() []string {
	names := []string{}
	for name := range ts.TableMap {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (ts *TableSet) TablesFormattedSorted() []*TableFormatter {
	tfs := []*TableFormatter{}
	names := ts.TableNames()
	for _, name := range names {
		if tf, ok := ts.TableMap[name]; ok {
			tfs = append(tfs, tf)
		}
	}
	return tfs
}

func (ts *TableSet) AddRecord(tableName string, row []string) {
	tableName = strings.TrimSpace(tableName)
	tf, ok := ts.TableMap[tableName]
	if !ok {
		tbl := NewTableData()
		tbl.Name = tableName
		tbl.Columns = ts.Columns
		tf = &TableFormatter{
			Table:     &tbl,
			Formatter: ts.RowFormatter}
	}
	tf.Table.Records = append(tf.Table.Records, row)
	ts.TableMap[tableName] = tf
}
