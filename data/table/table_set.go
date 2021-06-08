package table

import (
	"sort"
	"strings"
)

type TableSet struct {
	Columns    []string
	FormatMap  map[int]string
	FormatFunc func(val string, colIdx uint) (interface{}, error)
	TableMap   map[string]*Table
	//TableMapOld  map[string]*TableFormatter
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
	tbls := []*Table{}
	names := ts.TableNames()
	for _, name := range names {
		if tbl, ok := ts.TableMap[name]; ok {
			tbls = append(tbls, tbl)
		}
	}
	return tbls
}

/*
func (ts *TableSet) TablesFormattedSorted() []*TableFormatter {
	tfs := []*TableFormatter{}
	names := ts.TableNames()
	for _, name := range names {
		if tf, ok := ts.TableMap[name]; ok {
			tfs = append(tfs, tf)
		}
	}
	return tfs
}*/
/*

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
*/

func (ts *TableSet) AddRecord(tableName string, row []string) {
	tableName = strings.TrimSpace(tableName)
	tbl, ok := ts.TableMap[tableName]
	if !ok {
		tbl := NewTable()
		tbl.Name = tableName
		tbl.Columns = ts.Columns
		tbl.FormatMap = ts.FormatMap
		tbl.FormatFunc = ts.FormatFunc
	}
	tbl.Records = append(tbl.Records, row)
	ts.TableMap[tableName] = tbl
}
