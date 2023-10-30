package table

import (
	"strings"

	"github.com/grokify/gocharts/v2/data/table/excelizeutil"
	"github.com/grokify/mogo/type/maputil"
	excelize "github.com/xuri/excelize/v2"
)

type TableSet struct {
	Name       string
	Columns    []string
	FormatMap  map[int]string
	FormatFunc func(val string, colIdx uint) (any, error)
	TableMap   map[string]*Table
	Order      []string
}

func NewTableSet(name string) *TableSet {
	return &TableSet{
		Name:      name,
		Columns:   []string{},
		FormatMap: map[int]string{},
		TableMap:  map[string]*Table{}}
}

func (ts *TableSet) TableNames() []string {
	return maputil.Keys(ts.TableMap)
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
	return WriteXLSX(filename, tbls)
}

func ReadFileXLSX(filename string, headerRowCount uint, trimSpace bool) (*TableSet, error) {
	ts := NewTableSet("")
	xm, err := excelizeutil.NewExcelizeMore(filename)
	if err != nil {
		return nil, err
	}
	sheetNames := xm.SheetList()
	for _, sheetName := range sheetNames {
		cols, rows, err := xm.TableData(sheetName, headerRowCount, trimSpace)
		if err != nil {
			return nil, err
		}
		tbl := NewTable(sheetName)
		tbl.Columns = cols
		tbl.Rows = rows
		ts.TableMap[sheetName] = &tbl
	}

	return ts, xm.Close()
}

func XSLXGetSheetTable(f *excelize.File, sheetName string, headerRowCount uint, trimSpace bool) (*Table, error) {
	xm := excelizeutil.ExcelizeMore{File: f}
	cols, rows, err := xm.TableData(sheetName, headerRowCount, trimSpace)
	if err != nil {
		return nil, err
	}
	tbl := NewTable(sheetName)
	tbl.Columns = excelizeutil.ColumnsCollapse([][]string{cols}, trimSpace)
	tbl.Rows = rows
	return &tbl, nil
}
