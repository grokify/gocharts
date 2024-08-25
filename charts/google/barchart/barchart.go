package barchart

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/grokify/gocharts/v2/charts/google"
	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gocharts/v2/data/timeseries"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/time/timeutil"
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

func (chart *Chart) WriteFilePage(filename string, perm os.FileMode) error {
	pg := BarChartMaterialPage(chart)
	return os.WriteFile(filename, []byte(pg), perm)
}

func DataTableFromHistogram(h *histogram.Histogram, inclUnordered, inclZeroCount, inclZeroCountTail bool) (google.DataTable, error) {
	dt := google.DataTable{}
	if h == nil {
		return dt, errors.New("histogram must be supplied")
	}
	cols := []any{h.Name, "Count"}
	dt = append(dt, cols)

	bins := h.OrderOrDefault(inclUnordered)
	idxDtLastNonZero := -1
	for _, binName := range bins {
		cnt := h.GetOrDefault(binName, 0)
		if cnt != 0 {
			dt = append(dt, []any{binName, cnt})
			idxDtLastNonZero = len(dt) - 1
		} else if inclZeroCount || inclZeroCountTail {
			dt = append(dt, []any{binName, cnt})
		}
	}
	if !inclZeroCountTail {
		return dt[:idxDtLastNonZero+1], nil
	} else {
		return dt, nil
	}
}

// func DataTableFromTimeSeriesSet(name string, sets []string, set timeseries.TimeSeriesSet) (google.DataTable, error) {
func DataTableFromTimeSeriesSet(name string, sets []string, set timeseries.TimeSeriesSet) (google.DataTable, error) {
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
	IsStackedDefault  = IsStackedAbsolute
)

type Options struct {
	Width          uint          `json:"width"`
	Height         uint          `json:"height"`
	Legend         OptionsLegend `json:"legend,omitempty"`
	Bar            OptionsBar    `json:"bar"`
	IsStacked      string        `json:"isStacked"`
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

type OptionsLegend struct {
	Position string `json:"position,omitempty"`
	MaxLines int    `json:"maxLines,omitempty"`
}

type OptionsBar struct {
	GroupWidth string `json:"groupWidth,omitempty"`
}

type OptionsHAxis struct {
	MinValue int   `json:"minValue"`
	Ticks    []int `json:"ticks"`
}
