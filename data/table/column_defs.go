package table

import (
	"fmt"
	"slices"
	"strings"
)

// ColumnDefinitions provides a way to load definitions with their format types using
// the `Table.LoadColumnDefinitions()` method.
type ColumnDefinitionSet struct {
	DefaultFormat string
	Definitions   ColumnDefinitions
}

// BuildColumnDefinitionSet returns a `ColumnDefinitions{}` struct given a set of column
// names and a format map, similar to stored in `Table{}`.`
func BuildColumnDefinitionSet(names []string, formatMap map[int]string) ColumnDefinitionSet {
	cds := ColumnDefinitionSet{}
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

type ColumnDefinitions []ColumnDefinition

func (cd ColumnDefinitions) SourceNames(defToName, errOnMissing bool) ([]string, error) {
	var out []string
	for i, cdi := range cd {
		if cdi.SourceName != "" {
			out = append(out, cdi.SourceName)
		} else if defToName {
			out = append(out, cdi.Name)
		} else if errOnMissing {
			return out, fmt.Errorf("empty at index %d", i)
		} else {
			out = append(out, "")
		}
	}
	return out, nil
}

// ColumnDefinition represents one column including its name and format.
type ColumnDefinition struct {
	Name         string // internal name
	Format       string
	ExportName   string // for exports only
	SourceName   string // used for imports
	DefaultValue string
}

// LoadColumnDefinitions loads a set of column definitions with names and formats. It
// can be used to add columns to a table without any existing columns or to add to a
// table with existing columns.
func (tbl *Table) LoadColumnDefinitionSet(colDefs ColumnDefinitionSet) {
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
func (tbl *Table) ColumnDefinitionSet() ColumnDefinitionSet {
	return BuildColumnDefinitionSet(slices.Clone(tbl.Columns), tbl.FormatMap)
}
