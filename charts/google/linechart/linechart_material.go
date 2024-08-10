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
	"github.com/grokify/mogo/time/timeutil"
)

// LineChartMaterial provides data for Google Material Line Charts described here:
// https://developers.google.com/chart/interactive/docs/gallery/linechart#examples
type LineChartMaterial struct {
	Title    string
	Subtitle string
	ChartDiv string
	Width    int
	Height   int
	Columns  []google.Column
	Data     [][]any
}

func NewLineChartMaterial() LineChartMaterial {
	return LineChartMaterial{
		Columns: []google.Column{},
		Data:    [][]any{}}
}

func (lcm *LineChartMaterial) LoadTimeSeriesSetMonth(tss *timeseries.TimeSeriesSet, fn func(t time.Time) string) error {
	if cols, rows, err := TimeSeriesSetToLineChartMaterial(tss, fn); err != nil {
		return err
	} else {
		lcm.Columns = cols
		lcm.Data = rows
		return nil
	}
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
	return google.DefaultChartDiv
}

func (lcm *LineChartMaterial) HeightOrDefault() int {
	if lcm.Height > 0 {
		return lcm.Height
	}
	return google.DefaultHeight
}

func (lcm *LineChartMaterial) WidthOrDefault() int {
	if lcm.Width > 0 {
		return lcm.Width
	}
	return google.DefaultWidth
}

func (lcm *LineChartMaterial) PageHTML() string {
	return LineChartMaterialPage(*lcm)
}

func (lcm *LineChartMaterial) WriteFilePage(filename string, perm os.FileMode) error {
	return os.WriteFile(filename, []byte(lcm.PageHTML()), perm)
}

func LineChartMaterialFromTimeSeriesSet(tss timeseries.TimeSeriesSet, yearLabel string) (LineChartMaterial, error) {
	lcm := LineChartMaterial{}
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
