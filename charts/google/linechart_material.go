package google

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/v2/data/timeseries"
	"github.com/grokify/mogo/time/timeutil"
)

const (
	DefaultWidth    = 900
	DefaultHeight   = 500
	DefaultChartDiv = "chart_div"
	TypeNumber      = "number"
	TypeString      = "string"
)

// LineChartMaterial provides data for Google Material Line Charts described here:
// https://developers.google.com/chart/interactive/docs/gallery/linechart#examples
type LineChartMaterial struct {
	Title    string
	Subtitle string
	ChartDiv string
	Width    int
	Height   int
	Columns  []Column
	Data     [][]any
}

type Column struct {
	Type string
	Name string
}

func (lcm *LineChartMaterial) DataMatrixJSON() []byte {
	bytes, err := json.Marshal(lcm.Data)
	if err != nil {
		return []byte("[]")
	}
	return bytes
}

func (lcm *LineChartMaterial) ChartDivOrDefault() string {
	if len(lcm.ChartDiv) > 0 {
		return lcm.ChartDiv
	}
	return DefaultChartDiv
}

func (lcm *LineChartMaterial) HeightOrDefault() int {
	if lcm.Height > 0 {
		return lcm.Height
	}
	return DefaultHeight
}

func (lcm *LineChartMaterial) WidthOrDefault() int {
	if lcm.Width > 0 {
		return lcm.Width
	}
	return DefaultWidth
}

func LineChartMaterialFromTimeSeriesSet(tss timeseries.TimeSeriesSet, yearLabel string) (LineChartMaterial, error) {
	lcm := LineChartMaterial{}
	if tss.Interval != timeutil.IntervalYear {
		return lcm, errors.New("interval not supported")
	}

	if len(strings.TrimSpace(yearLabel)) == 0 {
		yearLabel = "Year"
	}
	lcmCols := []Column{
		{Type: TypeString, Name: yearLabel},
	}
	for _, seriesName := range tss.Order {
		lcmCols = append(lcmCols, Column{Type: TypeNumber, Name: seriesName})
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
