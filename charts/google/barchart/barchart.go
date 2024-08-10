package barchart

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/grokify/gocharts/v2/charts/google"
	"github.com/grokify/gocharts/v2/data/timeseries"
	"github.com/grokify/mogo/time/timeutil"
)

// BarChartMaterial represents the chart at:
// https://developers-dot-devsite-v2-prod.appspot.com/chart/interactive/docs/gallery/barchart
type BarChartMaterial struct {
	Title     string
	ChartDiv  string
	DataTable google.DataTable
	Options   BarChartOptions
}

func (chart BarChartMaterial) ChartDivOrDefault() string {
	div := strings.TrimSpace(chart.ChartDiv)
	if div != "" {
		return div
	} else {
		return google.DefaultChartDiv
	}
}

func (chart BarChartMaterial) PageTitle() string {
	return chart.Title
}

func (chart BarChartMaterial) DataTableJSON() []byte {
	return chart.DataTable.MustJSON()
}

func (chart BarChartMaterial) OptionsJSON() []byte {
	return chart.Options.MustJSON()
}

func (chart BarChartMaterial) WritePageHTML(filename string, perm os.FileMode) error {
	pg := BarChartMaterialPage(chart)
	return os.WriteFile(filename, []byte(pg), perm)
}

func TimeSeriesSetToDataTable(name string, sets []string, set timeseries.TimeSeriesSet) (google.DataTable, error) {
	dt := google.DataTable{}
	if len(sets) == 0 {
		sets = set.SeriesNames()
	}
	row1 := []any{name}
	for _, set := range sets {
		row1 = append(row1, set)
	}
	row1 = append(row1, map[string]string{"role": "annotation"})
	dt = append(dt, row1)
	if set.Interval == timeutil.IntervalMonth {
		timeStrings := set.TimeStrings()
		for _, ts := range timeStrings {
			t, err := time.Parse(time.RFC3339, ts)
			if err != nil {
				return dt, err
			}
			mDisplay := t.Format("Jan 2006")
			row := []any{mDisplay}
			for _, sname := range sets {
				val := set.GetInt64WithDefault(sname, ts, 0)
				row = append(row, val)
			}
			row = append(row, "")
			dt = append(dt, row)
		}
	}
	return dt, nil
}

const (
	IsStackedAbsolute = "absolute"
	IsStackedPercent  = "percent"
	IsStackedRelative = "relative"
)

type BarChartOptions struct {
	Width          uint                  `json:"width"`
	Height         uint                  `json:"height"`
	Legend         BarChartOptionsLegend `json:"legend,omitempty"`
	Bar            BarChartOptionsBar    `json:"bar"`
	IsStacked      string                `json:"isStacked"`
	HorizontalAxis BarChartOptionsHAxis  `json:"hAxis"`
}

func BarChartOptionsDefault() BarChartOptions {
	return BarChartOptions{
		Width:     600,
		Height:    400,
		Legend:    BarChartOptionsLegend{Position: "top", MaxLines: 3},
		Bar:       BarChartOptionsBar{GroupWidth: "75%"},
		IsStacked: IsStackedAbsolute,
	}
}

func (opts BarChartOptions) MustJSON() []byte {
	if b, err := json.Marshal(opts); err != nil {
		return []byte("[]")
	} else {
		return b
	}
}

type BarChartOptionsLegend struct {
	Position string `json:"position,omitempty"`
	MaxLines int    `json:"maxLines,omitempty"`
}

type BarChartOptionsBar struct {
	GroupWidth string `json:"groupWidth,omitempty"`
}

type BarChartOptionsHAxis struct {
	MinValue int   `json:"minValue"`
	Ticks    []int `json:"ticks"`
}
