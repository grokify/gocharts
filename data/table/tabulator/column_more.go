// tabulator provides helper methods for rendering HTML
// with Tabulator (http://tabulator.info/)
package tabulator

import "encoding/json"

type Column struct {
	Display string
	Slug    string
	Width   float64
}

type ColumnSet struct {
	Columns []Column
}

func (set *ColumnSet) DisplayTexts() []string {
	displays := []string{}
	for _, col := range set.Columns {
		displays = append(displays, col.Display)
	}
	return displays
}

type TabulatorColumn struct {
	Title        string  `json:"title,omitempty"`
	Field        string  `json:"field,omitempty"`
	Width        float64 `json:"width,omitempty"`
	HeaderFilter string  `json:"headerFilter,omitempty"`
}

type TabulatorColumnSet struct {
	Columns []TabulatorColumn
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
