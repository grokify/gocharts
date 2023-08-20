package table

import (
	"fmt"
	"strings"
)

// ColumnDefinitions provides a way to load definitions with their format types using
// the `Table.LoadColumnDefinitions()` method.
type ColumnDefinitions struct {
	DefaultFormat string
	Definitions   []ColumnDefinition
}

type ColumnDefinition struct {
	Name   string
	Format string
}

// LoadColumnDefinitions loads a set of column definitions with names and formats. It
// can be used to add columns to a table without any existing columns or to add to a
// table with existing columns.
func (tbl *Table) LoadColumnDefinitions(colDefs *ColumnDefinitions) {
	if colDefs == nil {
		return
	}
	if tbl.Columns == nil {
		tbl.Columns = Columns{}
	}
	if tbl.FormatMap == nil {
		tbl.FormatMap = map[int]string{}
	}
	for _, def := range colDefs.Definitions {
		colIdx := len(tbl.Columns)
		colName := strings.TrimSpace(def.Name)
		if colName == "" {
			colName = fmt.Sprintf("Column %d", (colIdx + 1))
		}
		tbl.Columns = append(tbl.Columns, colName)
		defFormat := strings.TrimSpace(def.Format)
		if defFormat != "" {
			tbl.FormatMap[colIdx] = defFormat
		}
	}
	defaultType := strings.TrimSpace(colDefs.DefaultFormat)
	if defaultType != "" {
		tbl.FormatMap[-1] = defaultType
	}
}
