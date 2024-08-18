package piechart

import (
	"encoding/json"
	"io"
	"os"

	"github.com/grokify/gocharts/v2/charts/google"
	"github.com/grokify/gocharts/v2/data/piechart"
	"github.com/grokify/gocharts/v2/data/table"
)

// Chart provides data for Google Material Pie Charts described here:
// https://developers.google.com/chart/interactive/docs/gallery/piechart
type Chart struct {
	Title               string
	Subtitle            string
	ChartDiv            string
	Width               int
	Height              int
	AddCountToName      bool
	DefaultCategoryName string
	DataTable           *google.DataTable
	Columns             google.Columns
	Data                piechart.PieChartData
	GoogleOptions       Options
}

func NewPieChartMaterialInts(chartName, sliceName, sliceValueName string, vals map[string]int) Chart {
	c := Chart{
		Title: chartName,
		Columns: google.Columns{
			{Name: sliceName, Type: table.FormatString},
			{Name: sliceValueName, Type: table.FormatInt},
		},
		Data: piechart.PieChartData{IsFloat: false},
		GoogleOptions: Options{
			Title: chartName,
		},
	}
	c.Data.AddInts(vals)
	return c
}

func (cm *Chart) BuildDataTable() google.DataTable {
	if cm.DataTable != nil {
		return *cm.DataTable
	}
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

func (cm *Chart) DataTableJSON() []byte {
	matrix := cm.BuildDataTable()
	return matrix.MustJSON()
}

func (cm *Chart) OptionsJSON() []byte {
	return cm.GoogleOptions.MustJSON()
}

func (cm *Chart) ChartDivOrDefault() string {
	if len(cm.ChartDiv) > 0 {
		return cm.ChartDiv
	}
	return google.DefaultChartDiv
}

func (cm *Chart) HeightOrDefault() int {
	if cm.Height > 0 {
		return cm.Height
	}
	return google.DefaultHeight
}

func (cm *Chart) WidthOrDefault() int {
	if cm.Width > 0 {
		return cm.Width
	}
	return google.DefaultWidth
}

func (cm *Chart) PageHTML() string {
	return PieChartMaterialPage(*cm)
}

func (cm *Chart) WritePage(w io.Writer) {
	WritePieChartMaterialPage(w, *cm)
}

func (cm *Chart) WriteFilePage(filename string, perm os.FileMode) error {
	return os.WriteFile(filename, []byte(cm.PageHTML()), perm)
}

const (
	PieSliceTextLabel = "label"
)

// Options represents the Google Charts JSON options map as defined here:
// https://developers.google.com/chart/interactive/docs/gallery/piechart .
type Options struct {
	Title             string    `json:"title,omitempty"`
	Legend            string    `json:"legend,omitempty"`
	Height            uint      `json:"height,omitempty"`
	Width             uint      `json:"width,omitempty"`
	PieHole           float64   `json:"pieHole,omitempty"`
	PieSliceText      string    `json:"pieSliceText,omitempty"`
	PieSliceTextStyle TextStyle `json:"pieSliceTextStyle,omitempty"`
	PieStartAngle     float64   `json:"pieStartAngle,omitempty"`
}

// MustJSON represents the Google Charts JSON options map as defined here:
// https://developers.google.com/chart/interactive/docs/gallery/piechart .
// The output is intended to be used directly with the client-side JS library call.
func (opts Options) MustJSON() []byte {
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
