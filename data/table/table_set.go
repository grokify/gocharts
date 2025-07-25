package table

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/grokify/gocharts/v2/data/table/excelizeutil"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/type/maputil"
	excelize "github.com/xuri/excelize/v2"
)

type TableSet struct {
	Name       string
	Columns    []string
	FormatMap  map[int]string
	FormatFunc func(val string, colIdx uint32) (any, error)
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

func (ts *TableSet) Add(tbl ...*Table) error {
	for _, t := range tbl {
		if t == nil {
			continue
		}
		name := t.Name
		if strings.TrimSpace(name) == "" {
			name = "Table " + strconv.Itoa(len(ts.TableMap)+1)
		}
		if _, ok := ts.TableMap[name]; ok {
			return errors.New("table name collision")
		}
		ts.TableMap[name] = t
		ts.Order = append(ts.Order, name)
	}
	return nil
}

func (ts *TableSet) ColumnSetByTableName() (map[string]Columns, error) {
	out := map[string]Columns{}
	for _, t := range ts.TableMap {
		if _, ok := out[t.Name]; ok {
			return out, fmt.Errorf("column name duplicate (%s)", t.Name)
		} else {
			out[t.Name] = t.Columns
		}
	}
	return out, nil
}

// LensOrdered returns a slice of table row counts.
func (ts *TableSet) LensOrdered() []int {
	var lens []int
	for _, ord := range ts.Order {
		if tbl, ok := ts.TableMap[ord]; ok {
			if tbl != nil {
				lens = append(lens, len(tbl.Rows))
			} else {
				lens = append(lens, 0)
			}
		} else {
			lens = append(lens, 0)
		}
	}
	if len(lens) != len(ts.Order) {
		panic("internal mismatch in `TableSet.LensOrdered(), calculated length vs. `TableSet.Order`")
	}
	return lens
}

func (ts *TableSet) Table(name string) (*Table, error) {
	if t, ok := ts.TableMap[name]; ok {
		return t, nil
	} else {
		return nil, fmt.Errorf("table name not found (%s)", name)
	}
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
	if err := WriteXLSX(filename, tbls); err != nil {
		return errorsutil.Wrapf(err, "error in TableSet.WriteXLSX(%s)", filename)
	} else {
		return nil
	}
}

// ReadTableSetXLSXFile reads in an entire XLSX file as a `TableSet`. Warning: this can be resource
// intensive if there's a lot of data. If you just want one sheet of many, it is better to
// extract an individual sheet or sheets from an `excelize.File`, such as using `ParseTableXLSX()` or `ParseTableXLSXIndex()`.
func ReadTableSetXLSXFile(filename string, headerRowCount uint32, trimSpace bool) (*TableSet, error) {
	ts := NewTableSet("")
	xm, err := excelizeutil.NewFile(filename)
	if err != nil {
		return nil, err
	}
	sheetNames := xm.SheetNames(false)
	for _, sheetName := range sheetNames {
		cols, rows, err := xm.TableData(sheetName, headerRowCount, trimSpace, false)
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

func ReadTableXLSXFile(filename, sheetName string, headerRowCount uint32, trimSpace bool) (*Table, error) {
	if xf, err := excelizeutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return ParseTableXLSX(xf.File, sheetName, headerRowCount, trimSpace)
	}
}

func ReadTableXLSXIndexFile(filename string, sheetIdx, headerRowCount uint32, trimSpace bool) (*Table, error) {
	if xf, err := excelizeutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return ParseTableXLSXIndex(xf.File, sheetIdx, headerRowCount, trimSpace)
	}
}

func ParseTableXLSX(f *excelize.File, sheetName string, headerRowCount uint32, trimSpace bool) (*Table, error) {
	if f == nil {
		return nil, excelizeutil.ErrExcelizeFileCannotBeNil
	}
	xm := excelizeutil.File{File: f}
	if cols, rows, err := xm.TableData(sheetName, headerRowCount, trimSpace, false); err != nil {
		return nil, err
	} else {
		tbl := NewTable(sheetName)
		tbl.Columns = cols
		tbl.Rows = rows
		return &tbl, nil
	}
}

func ParseTableXLSXIndex(f *excelize.File, sheetIdx, headerRowCount uint32, trimSpace bool) (*Table, error) {
	if f == nil {
		return nil, excelizeutil.ErrExcelizeFileCannotBeNil
	}
	xm := excelizeutil.File{File: f}
	if cols, rows, err := xm.TableDataIndex(sheetIdx, headerRowCount, trimSpace, false); err != nil {
		return nil, err
	} else {
		tbl := NewTable(f.GetSheetName(int(sheetIdx)))
		tbl.Columns = cols
		tbl.Rows = rows
		return &tbl, nil
	}
}
