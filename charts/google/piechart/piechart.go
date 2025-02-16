package piechart

import (
	"io"
	"os"
	"strings"

	"github.com/grokify/gocharts/v2/charts/google"
	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gocharts/v2/data/piechart"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/encoding/jsonutil"
)

// Chart provides data for Google Material Pie Charts described here:
// https://developers.google.com/chart/interactive/docs/gallery/piechart
type Chart struct {
	Title               string
	Subtitle            string
	ChartDiv            string
	AddCountToName      bool
	DefaultCategoryName string
	DataTable           *google.DataTable
	Columns             google.Columns
	Data                piechart.PieChartData
	GoogleOptions       *Options
}

func NewPieChartMaterialInts(chartName, sliceName, sliceValueName string, vals map[string]int) *Chart {
	chart := Chart{
		Title: chartName,
		Columns: google.Columns{
			{Name: sliceName, Type: table.FormatString},
			{Name: sliceValueName, Type: table.FormatInt},
		},
		Data: piechart.PieChartData{IsFloat: false},
		GoogleOptions: &Options{
			Title: chartName,
		},
	}
	chart.Data.AddInts(vals)
	return &chart
}

func (chart *Chart) LoadDataTableHistogram(h *histogram.Histogram, cols google.Columns) {
	if len(cols) >= 0 {
		chart.Columns = cols
	}
	colNamesAny := chart.Columns.NamesAny()
	if len(colNamesAny) < 2 {
		colNamesAny = []any{"Categories", "Value"}
	}
	var dt = google.DataTable{colNamesAny}
	binNames := h.BinNamesMore(true, true, true)
	for _, binName := range binNames {
		binVal := h.GetOrDefault(binName, 0)
		dt = append(dt, []any{binName, binVal})
	}
	chart.DataTable = &dt
}

func (chart *Chart) BuildDataTable() google.DataTable {
	if chart.DataTable != nil {
		return *chart.DataTable
	}
	colNamesAny := chart.Columns.NamesAny()
	if len(colNamesAny) < 2 {
		colNamesAny = []any{"Categories", "Value"}
	}
	var dt = google.DataTable{colNamesAny}
	chart.Data.Sort()
	for _, d := range chart.Data.Data {
		var row []any
		name := d.Name
		if chart.AddCountToName {
			name = d.NameWithCount(chart.DefaultCategoryName)
		}
		if d.IsFloat {
			row = append(row, name, d.ValFloat)
		} else {
			row = append(row, name, d.ValInt)
		}
		dt = append(dt, row)
	}
	return dt
}

func (chart *Chart) DataTableJSON() []byte {
	dt := chart.BuildDataTable()
	return dt.MustJSON()
}

func (chart *Chart) ChartDivOrDefault() string {
	if strings.TrimSpace(chart.ChartDiv) != "" {
		return chart.ChartDiv
	}
	return google.DefaultChartDiv
}

func (chart *Chart) OptionsJSON() []byte {
	if chart.GoogleOptions == nil {
		return []byte(jsonutil.EmptyObject)
	} else {
		return chart.GoogleOptions.MustJSON()
	}
}

func (chart *Chart) HTML() string          { return PieChartMaterialHTML(chart) }
func (chart *Chart) PageHTML() string      { return PieChartMaterialHTMLPage(chart) }
func (chart *Chart) PageTitle() string     { return chart.Title }
func (chart *Chart) WritePage(w io.Writer) { WritePieChartMaterialHTMLPage(w, chart) }

func (chart *Chart) WriteFilePageHTML(filename string, perm os.FileMode) error {
	return os.WriteFile(filename, []byte(chart.PageHTML()), perm)
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

func (opts *Options) Inflate() {
	if opts.Height == 0 {
		opts.Height = google.DefaultHeight
	}
	if opts.Width == 0 {
		opts.Width = google.DefaultWidth
	}
}

// MustJSON represents the Google Charts JSON options map as defined here:
// https://developers.google.com/chart/interactive/docs/gallery/piechart .
// The output is intended to be used directly with the client-side JS library call.
func (opts *Options) MustJSON() []byte {
	opts.Inflate()
	return jsonutil.MustMarshalOrDefault(opts, []byte(jsonutil.EmptyObject))
}

type TextStyle struct {
	// https://developers.google.com/chart/interactive/docs/gallery/piechart
	Color    string  `json:"color,omitempty"`
	FontName string  `json:"fontName,omitempty"`
	FontSize float64 `json:"fontSize,omitempty"`
	Bold     bool    `json:"bold,omitempty"`
	Italic   bool    `json:"italic,omitempty"`
}
