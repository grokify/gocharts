// tabulator provides helper methods for rendering HTML with Tabulator (http://tabulator.info/)
package tabulator

import (
	"encoding/json"

	"github.com/grokify/mogo/text/stringcase"
)

type Column struct {
	Display string
	Slug    string
	Width   float64
}

type Columns []Column

type ColumnSet struct {
	Columns Columns
}

func NewColumnSetSimple(displaynames []string, totalWidth float64) *ColumnSet {
	cs := &ColumnSet{Columns: []Column{}}
	cw := float64(0)
	if len(displaynames) > 0 {
		if totalWidth < 0 {
			totalWidth *= -1
		}
		cw = totalWidth / float64(len(displaynames))
	}
	for _, name := range displaynames {
		col := Column{
			Display: name,
			Slug:    stringcase.ToCamelCase(name),
			Width:   cw,
		}
		cs.Columns = append(cs.Columns, col)
	}
	return cs
}

func (set ColumnSet) DisplayTexts() []string {
	displays := []string{}
	for _, col := range set.Columns {
		displays = append(displays, col.Display)
	}
	return displays
}

type TabulatorColumn struct {
	Title           string           `json:"title,omitempty"`
	Field           string           `json:"field,omitempty"`
	Formatter       string           `json:"formatter,omitempty"`
	FormatterParams *FormatterParams `json:"formatterParams,omitempty"`
	Width           float64          `json:"width,omitempty"`
	HeaderFilter    string           `json:"headerFilter,omitempty"`
}

type TabulatorColumnSet struct {
	Columns []TabulatorColumn
}

func (columns Columns) TabulatorColumnSet() TabulatorColumnSet {
	return BuildColumnsTabulator(columns)
}

func BuildColumnsTabulator(columns []Column) TabulatorColumnSet {
	colsTabulator := []TabulatorColumn{}
	for _, col := range columns {
		colT := TabulatorColumn{
			Title:        col.Display,
			Field:        col.Display,
			HeaderFilter: "input"}
		if col.Width > 0 {
			colT.Width = col.Width
		}
		colsTabulator = append(colsTabulator, colT)
	}
	return TabulatorColumnSet{Columns: colsTabulator}
}

func (tColSet *TabulatorColumnSet) ColumnsJSON() ([]byte, error) {
	return json.Marshal(tColSet.Columns)
}

func (tColSet *TabulatorColumnSet) MustColumnsJSON() []byte {
	bytes, err := tColSet.ColumnsJSON()
	if err != nil {
		return []byte("[]")
	}
	return bytes
}
