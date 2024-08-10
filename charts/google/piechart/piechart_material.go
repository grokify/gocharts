package piechart

import (
	"encoding/json"
	"io"
	"os"

	"github.com/grokify/gocharts/v2/charts/google"
	"github.com/grokify/gocharts/v2/data/piechart"
	"github.com/grokify/gocharts/v2/data/table"
)

// PieChartMaterial provides data for Google Material Pie Charts described here:
// https://developers.google.com/chart/interactive/docs/gallery/piechart
type PieChartMaterial struct {
	Title               string
	Subtitle            string
	ChartDiv            string
	Width               int
	Height              int
	AddCountToName      bool
	DefaultCategoryName string
	Columns             google.Columns
	Data                piechart.PieChartData
	GoogleOptions       PieChartOptionsGoogle
}

func NewPieChartMaterialInts(chartName, sliceName, sliceValueName string, vals map[string]int) PieChartMaterial {
	c := PieChartMaterial{
		Title: chartName,
		Columns: google.Columns{
			{Name: sliceName, Type: table.FormatString},
			{Name: sliceValueName, Type: table.FormatInt},
		},
		Data: piechart.PieChartData{IsFloat: false},
		GoogleOptions: PieChartOptionsGoogle{
			Title: chartName,
		},
	}
	c.Data.AddInts(vals)
	return c
}

func (cm *PieChartMaterial) DataTable() google.DataTable {
	colNamesAny := cm.Columns.NamesAny()
	if len(colNamesAny) < 2 {
		colNamesAny = []any{"Categories", "Value"}
	}
	var matrix = [][]any{colNamesAny}
	cm.Data.Sort()
	for _, d := range cm.Data.Data {
		var row []any
		name := d.Name
		if cm.AddCountToName {
			name = d.NameWithCount(cm.DefaultCategoryName)
		}
		if d.IsFloat {
			row = append(row, name, d.ValFloat)
		} else {
			row = append(row, name, d.ValInt)
		}
		matrix = append(matrix, row)
	}
	return matrix
}

func (cm *PieChartMaterial) DataTableJSON() []byte {
	matrix := cm.DataTable()
	return matrix.MustJSON()
}

func (cm *PieChartMaterial) OptionsJSON() []byte {
	return cm.GoogleOptions.MustJSON()
}

func (cm *PieChartMaterial) ChartDivOrDefault() string {
	if len(cm.ChartDiv) > 0 {
		return cm.ChartDiv
	}
	return google.DefaultChartDiv
}

func (cm *PieChartMaterial) HeightOrDefault() int {
	if cm.Height > 0 {
		return cm.Height
	}
	return google.DefaultHeight
}

func (cm *PieChartMaterial) WidthOrDefault() int {
	if cm.Width > 0 {
		return cm.Width
	}
	return google.DefaultWidth
}

func (cm *PieChartMaterial) PageHTML() string {
	return PieChartMaterialPage(*cm)
}

func (cm *PieChartMaterial) WritePage(w io.Writer) {
	WritePieChartMaterialPage(w, *cm)
}

func (cm *PieChartMaterial) WriteFilePage(filename string, perm os.FileMode) error {
	return os.WriteFile(filename, []byte(cm.PageHTML()), perm)
}

// PieChartOptionsGoogle represents the Google Charts JSON options map as defined here:
// https://developers.google.com/chart/interactive/docs/gallery/piechart .
type PieChartOptionsGoogle struct {
	Title             string    `json:"title,omitempty"`
	Legend            string    `json:"legend,omitempty"`
	Height            string    `json:"height,omitempty"`
	Width             string    `json:"width,omitempty"`
	PieHole           float64   `json:"pieHole,omitempty"`
	PieSliceText      string    `json:"pieSliceText,omitempty"`
	PieSliceTextStyle TextStyle `json:"pieSliceTextStyle,omitempty"`
	PieStartAngle     float64   `json:"pieStartAngle,omitempty"`
}

// MustJSON represents the Google Charts JSON options map as defined here:
// https://developers.google.com/chart/interactive/docs/gallery/piechart .
// The output is intended to be used directly with the client-side JS library call.
func (opts PieChartOptionsGoogle) MustJSON() []byte {
	if b, err := json.Marshal(opts); err != nil {
		return []byte("{}")
	} else {
		return b
	}
}

type TextStyle struct {
	// https://developers.google.com/chart/interactive/docs/gallery/piechart
	Color    string  `json:"color,omitempty"`
	FontName string  `json:"fontName,omitempty"`
	FontSize float64 `json:"fontSize,omitempty"`
	Bold     bool    `json:"bold,omitempty"`
	Italic   bool    `json:"italic,omitempty"`
}
