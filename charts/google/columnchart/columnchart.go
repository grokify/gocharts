package columnchart

import (
	"os"
	"strings"

	"github.com/grokify/mogo/encoding/jsonutil"

	"github.com/grokify/gocharts/v2/charts/google"
)

// Chart represents the chart at:
// https://developers-dot-devsite-v2-prod.appspot.com/chart/interactive/docs/gallery/barchart
type Chart struct {
	Title     string
	ChartDiv  string
	DataTable google.DataTable
	Options   *Options
}

func (chart *Chart) ChartDivOrDefault() string {
	if div := strings.TrimSpace(chart.ChartDiv); div != "" {
		return div
	} else {
		return google.DefaultChartDiv
	}
}

func (chart *Chart) DataTableJSON() []byte {
	return chart.DataTable.MustJSON()
}

func (chart *Chart) OptionsJSON() []byte {
	if chart.Options == nil {
		return []byte(jsonutil.EmptyObject)
	} else {
		return chart.Options.MustJSON()
	}
}

func (chart *Chart) PageTitle() string { return chart.Title }

func (chart *Chart) WriteFilePageHTML(filename string, perm os.FileMode) error {
	pg := ColumnChartMaterialPage(chart)
	return os.WriteFile(filename, []byte(pg), perm)
}

type IsStacked string

const (
	IsStackedAbsolute IsStacked = "absolute"
	IsStackedPercent  IsStacked = "percent"
	IsStackedRelative IsStacked = "relative"
	IsStackedDefault  IsStacked = IsStackedAbsolute
)

type Options struct {
	Chart          OptionsChart  `json:"chart,omitempty"`
	Title          string        `json:"title,omitempty"`
	Subtitle       string        `json:"subtitle,omitempty"`
	Width          uint          `json:"width,omitempty"`
	Height         uint          `json:"height,omitempty"`
	Legend         OptionsLegend `json:"legend,omitempty"`
	Bar            OptionsBar    `json:"bar"`
	IsStacked      IsStacked     `json:"isStacked"`
	HorizontalAxis OptionsHAxis  `json:"hAxis"`
}

func OptionsDefault() Options {
	return Options{
		Width:     600,
		Height:    400,
		Legend:    OptionsLegend{Position: "top", MaxLines: 3},
		Bar:       OptionsBar{GroupWidth: "75%"},
		IsStacked: IsStackedAbsolute,
	}
}

func (opts *Options) MustJSON() []byte {
	return jsonutil.MustMarshalOrDefault(opts, []byte(jsonutil.EmptyObject))
}

type OptionsChart struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type OptionsLegend struct {
	Position string `json:"position,omitempty"` // ["none"]
	MaxLines int    `json:"maxLines,omitempty"`
}

type OptionsBar struct {
	GroupWidth string `json:"groupWidth,omitempty"` // e.g. "95%"
}

type OptionsHAxis struct {
	MinValue int   `json:"minValue"`
	Ticks    []int `json:"ticks"`
}
