package d3bullet

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grokify/gocharts/v2/data/bullet"
)

const (
	JSPath = "github.com/grokify/gocharts/charts/d3/d3bullet/d3bullet.js"
)

// D3BulletChartBuilder is used to create a default D3Bullet chart
// based on minimal data. Instantiate D3BulletChartBuilder and then
// call D3Bullet(). Then add to DataInt64 and then call GetBulletDataJSON.
type D3BulletChartBuilder struct {
	Title    string
	YStart   int64 // Default to 0
	YTarget  int64 // Mandatory
	YCurrent int64 // Mandatory
	XStart   int64 // Mandatory, e.g. time
	XCurrent int64 // Mandatory, e.g. time
	XEnd     int64 // Mandatory, e.g. time
}

func (d3builder *D3BulletChartBuilder) D3Bullet() BulletInt64 {
	return BulletChartToD3Bullet(d3builder.Bullet())
}

func (d3builder *D3BulletChartBuilder) Bullet() bullet.BulletChart {
	thisBulletChart := bullet.BulletChart{
		Title: d3builder.Title,
		ProjectionData: bullet.ProjectionDataInt64{
			Start:   d3builder.YStart,
			Target:  d3builder.YTarget,
			Current: d3builder.YCurrent,
		},
	}
	thisBulletChart.ProjectionData.CreateProjection(
		d3builder.XStart, d3builder.XCurrent, d3builder.XEnd,
	)
	thisBulletChart.Subtitle =
		thisBulletChart.ProjectionData.ToString([]string{"T", "C", "P", "D"}, true)
	return thisBulletChart
}

type DataInt64 struct {
	Bullets []BulletInt64
}

// GetBulletDataJSON returns JSON structure for using in JavaScript.
func (data *DataInt64) GetBulletDataJSON() []byte {
	bytes, _ := json.Marshal(data.Bullets)
	return bytes
}

type BulletInt64 struct {
	Title    string  `json:"title,omitempty"`
	Subtitle string  `json:"subtitle,omitempty"`
	Ranges   []int64 `json:"ranges,omitempty"`
	Measures []int64 `json:"measures,omitempty"`
	Markers  []int64 `json:"markers,omitempty"`
}

func BulletChartToD3Bullet(bullet bullet.BulletChart) BulletInt64 {
	return ProjectionToBulletInt64(bullet.ProjectionData, bullet.Title, bullet.Subtitle)
}

func ProjectionToBulletInt64(prjData bullet.ProjectionDataInt64, title string, subtitle string) BulletInt64 {
	rangeMax := int64(float64(prjData.Target) * 1.2)
	if prjData.Projection > rangeMax {
		rangeMax = int64(float64(prjData.Projection) * 1.2)
	}
	return BulletInt64{
		Title:    title,
		Subtitle: subtitle,
		Ranges:   []int64{0, prjData.Target, rangeMax},
		Measures: []int64{prjData.Current, prjData.Projection},
		Markers:  []int64{prjData.Target},
	}
}

type Data struct {
	Bullets []Bullet
}

func (data *Data) GetBulletDataJSON() []byte {
	bytes, _ := json.Marshal(data.Bullets)
	return bytes
}

type Bullet struct {
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	Ranges   []int  `json:"ranges,omitempty"`
	Measures []int  `json:"measures,omitempty"`
	Markers  []int  `json:"markers,omitempty"`
}

func ProjectionToBullet(prjData bullet.ProjectionDataInt, title string, subtitle string) Bullet {
	rangeMax := int(float64(prjData.Target) * 1.2)
	if prjData.Projection > rangeMax {
		rangeMax = int(float64(prjData.Projection) * 1.2)
	}
	return Bullet{
		Title:    title,
		Subtitle: subtitle,
		Ranges:   []int{0, prjData.Target, rangeMax},
		Measures: []int{prjData.Current, prjData.Projection},
		Markers:  []int{prjData.Target},
	}
}

func GetJS() []byte {
	jsFilePath := filepath.Clean(filepath.Join(os.Getenv("GOPATH"), "src", JSPath))
	bytes, err := os.ReadFile(jsFilePath) //nolint:gosec // G703: path is constructed from known constant JSPath
	if err != nil {
		return []byte("")
	}
	return bytes
}

func GetExampleCSS(includeHTML bool) string {
	css := `.bullet { font: 10px sans-serif; }
.bullet .marker { stroke: #000; stroke-width: 2px; }
.bullet .tick line { stroke: #666; stroke-width: .5px; }
.bullet .range.s0 { fill: #eee; }
.bullet .range.s1 { fill: #ddd; }
.bullet .range.s2 { fill: #ccc; }
.bullet .measure.s0 { fill: lightsteelblue; }
.bullet .measure.s1 { fill: steelblue; }
.bullet .title { font-size: 14px; font-weight: bold; }
.bullet .subtitle { fill: #999; }`
	if includeHTML {
		return fmt.Sprintf("<style>\n%v\n</style>", css)
	}
	return css
}

func GetExampleCSSGreenRed(includeHTML bool) string {
	// red ffb9b9 fdd1d1
	// grn bbffb9 cfffce
	css := `.bullet { font: 10px sans-serif; }
.bullet .marker { stroke: #000; stroke-width: 2px; }
.bullet .tick line { stroke: #666; stroke-width: .5px; }
.bullet .range.s0 { fill: #cfffce; }
.bullet .range.s1 { fill: #fdd1d1; }
.bullet .range.s2 { fill: #fdd1d1; }
.bullet .measure.s0 { fill: lightsteelblue; }
.bullet .measure.s1 { fill: steelblue; }
.bullet .title { font-size: 14px; font-weight: bold; }
.bullet .subtitle { fill: #999; }`
	if includeHTML {
		return fmt.Sprintf("<style>\n%v\n</style>", css)
	}
	return css
}

func GetExampleJSData() string {
	return `var data = [
  {"title":"Revenue","subtitle":"US$, in thousands","ranges":[150,225,300],"measures":[220,270],"markers":[250]},
  {"title":"Profit","subtitle":"%","ranges":[20,25,30],"measures":[21,23],"markers":[26]},
  {"title":"Order Size","subtitle":"US$, average","ranges":[350,500,600],"measures":[100,320],"markers":[550]},
  {"title":"New Customers","subtitle":"count","ranges":[1400,2000,2500],"measures":[1000,1650],"markers":[2100]},
  {"title":"Satisfaction","subtitle":"out of 5","ranges":[3.5,4.25,5],"measures":[3.2,4.7],"markers":[4.4]}
];`
}

func GetExampleJSVars() string {
	return `var margin = {top: 5, right: 40, bottom: 20, left: 120},
    width = 960 - margin.left - margin.right,
    height = 50 - margin.top - margin.bottom;

var chart = d3.bullet()
    .width(width)
    .height(height);`
}

func GetExampleJS() string {
	return `var svg = d3.select("body").selectAll("svg")
      .data(data)
    .enter().append("svg")
      .attr("class", "bullet")
      .attr("width", width + margin.left + margin.right)
      .attr("height", height + margin.top + margin.bottom)
    .append("g")
      .attr("transform", "translate(" + margin.left + "," + margin.top + ")")
      .call(chart);

  var title = svg.append("g")
      .style("text-anchor", "end")
      .attr("transform", "translate(-6," + height / 2 + ")");

  title.append("text")
      .attr("class", "title")
      .text(function(d) { return d.title; });

  title.append("text")
      .attr("class", "subtitle")
      .attr("dy", "1em")
      .text(function(d) { return d.subtitle; });`
}
