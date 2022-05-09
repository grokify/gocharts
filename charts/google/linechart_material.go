package google

import "encoding/json"

const (
	DefaultWidth    = 900
	DefaultHeight   = 500
	DefaultChartDiv = "chart_div"
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
	Data     [][]interface{}
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
