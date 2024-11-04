package table

import (
	"fmt"
	"slices"
	"strings"
)

// ColumnDefinitions provides a way to load definitions with their format types using
// the `Table.LoadColumnDefinitions()` method.
type ColumnDefinitions struct {
	DefaultFormat string
	Definitions   []ColumnDefinition
}

// BuildColumnDefinitions returns a `ColumnDefinitions{}` struct given a set of column
// names and a format map, similar to stored in `Table{}`.`
func BuildColumnDefinitions(names []string, formatMap map[int]string) ColumnDefinitions {
	cds := ColumnDefinitions{}
	if f, ok := formatMap[-1]; ok {
		cds.DefaultFormat = f
	}
	for i, name := range names {
		cd := ColumnDefinition{
			Name: name,
		}
		if f, ok := formatMap[i]; ok {
			cd.Format = f
		}
		cds.Definitions = append(cds.Definitions, cd)
	}
	return cds
}

// ColumnDefinition represents one column including its name and format.
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

// ColumnDefinitions returns a `ColumnDefinitions{}` struct for the `Table{}`.
func (tbl *Table) ColumnDefinitions() ColumnDefinitions {
	return BuildColumnDefinitions(
		slices.Clone(tbl.Columns),
		tbl.FormatMap)
}
