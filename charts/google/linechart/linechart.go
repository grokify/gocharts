package linechart

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/v2/charts/google"
	"github.com/grokify/gocharts/v2/data/timeseries"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/time/timeutil"
)

// Chart provides data for Google Material Line Charts described here:
// https://developers.google.com/chart/interactive/docs/gallery/linechart#examples
type Chart struct {
	Title         string
	Subtitle      string
	ChartDiv      string
	Width         int
	Height        int
	Columns       []google.Column
	Data          google.DataTable
	GoogleOptions *Options
}

func NewChart() Chart {
	return Chart{
		Columns: []google.Column{},
		Data:    google.DataTable{}}
}

func (chart *Chart) LoadTimeSeriesSetMonth(tss *timeseries.TimeSeriesSet, fn func(t time.Time) string) error {
	if cols, rows, err := TimeSeriesSetToLineChartMaterial(tss, fn); err != nil {
		return err
	} else {
		chart.Columns = cols
		chart.Data = rows
		return nil
	}
}

func (chart *Chart) ChartDivOrDefault() string {
	if len(chart.ChartDiv) > 0 {
		return chart.ChartDiv
	} else {
		return google.DefaultChartDiv
	}
}

func (chart *Chart) DataMatrixJSON() []byte {
	return jsonutil.MustMarshalOrDefault(chart.Data, []byte(jsonutil.EmptyArray))
}

func (chart *Chart) DataTableJSON() []byte {
	if bytes, err := json.Marshal(chart.Data); err != nil {
		return []byte(jsonutil.EmptyArray)
	} else {
		return bytes
	}
}

func (chart *Chart) OptionsJSON() []byte {
	if chart.GoogleOptions == nil {
		return []byte(jsonutil.EmptyObject)
	} else {
		return chart.GoogleOptions.MustJSON()
	}
}

func (chart *Chart) PageTitle() string { return chart.Title }

/*
func (chart *Chart) HeightOrDefault() int {
	if chart.Height > 0 {
		return chart.Height
	}
	return google.DefaultHeight
}

func (chart *Chart) WidthOrDefault() int {
	if chart.Width > 0 {
		return chart.Width
	}
	return google.DefaultWidth
}
*/

func (chart *Chart) PageHTML() string { return LineChartMaterialPage(*chart) }

func (chart *Chart) WriteFilePage(filename string, perm os.FileMode) error {
	return os.WriteFile(filename, []byte(chart.PageHTML()), perm)
}

func ChartFromTimeSeriesSet(tss timeseries.TimeSeriesSet, yearLabel string) (Chart, error) {
	lcm := Chart{}
	if tss.Interval != timeutil.IntervalYear {
		return lcm, errors.New("interval not supported")
	}

	if len(strings.TrimSpace(yearLabel)) == 0 {
		yearLabel = "Year"
	}
	lcmCols := []google.Column{
		{Type: google.TypeString, Name: yearLabel},
	}
	for _, seriesName := range tss.Order {
		lcmCols = append(lcmCols, google.Column{Type: google.TypeNumber, Name: seriesName})
	}
	lcm.Columns = lcmCols

	rows := [][]any{}

	for _, dt := range tss.Times {
		row := []any{strconv.Itoa(dt.Year())}

		for _, seriesName := range tss.Order {
			item, err := tss.Item(seriesName, dt.Format(time.RFC3339))
			if err != nil {
				row = append(row, 0)
			} else {
				row = append(row, item.Float64())
			}
		}

		rows = append(rows, row)
	}

	lcm.Data = rows

	return lcm, nil
}

type Options struct {
	Chart  OptionsChart `json:"chart,omitempty"`
	Height uint         `json:"height,omitempty"`
	Width  uint         `json:"width,omitempty"`
}

func (opts *Options) Inflate() {
	if opts.Height == 0 {
		opts.Height = google.DefaultHeight
	}
	if opts.Width == 0 {
		opts.Width = google.DefaultWidth
	}
}

func (opts *Options) MustJSON() []byte {
	return jsonutil.MustMarshalOrDefault(opts, []byte(jsonutil.EmptyObject))
}

type OptionsChart struct {
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
}
