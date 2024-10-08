// Code generated by qtc from "linechart_material.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line linechart_material.qtpl:1
package linechart

//line linechart_material.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line linechart_material.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line linechart_material.qtpl:1
func StreamLineChartMaterialPage(qw422016 *qt422016.Writer, chart Chart) {
//line linechart_material.qtpl:1
	qw422016.N().S(`<!DOCTYPE html>
<html>
<head>
  <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
</head>
<body>
  <h1>`)
//line linechart_material.qtpl:7
	qw422016.N().S(chart.PageTitle())
//line linechart_material.qtpl:7
	qw422016.N().S(`</h1>
  <div id="`)
//line linechart_material.qtpl:8
	qw422016.N().S(chart.ChartDivOrDefault())
//line linechart_material.qtpl:8
	qw422016.N().S(`"></div>
  <script>
      google.charts.load('current', {'packages':['line']});
      google.charts.setOnLoadCallback(drawChart);

    function drawChart() {
      var data = new google.visualization.DataTable();
      `)
//line linechart_material.qtpl:15
	for _, col := range chart.Columns {
//line linechart_material.qtpl:15
		qw422016.N().S(`
      data.addColumn('`)
//line linechart_material.qtpl:16
		qw422016.E().S(col.Type)
//line linechart_material.qtpl:16
		qw422016.N().S(`', '`)
//line linechart_material.qtpl:16
		qw422016.E().S(col.Name)
//line linechart_material.qtpl:16
		qw422016.N().S(`');
      `)
//line linechart_material.qtpl:17
	}
//line linechart_material.qtpl:17
	qw422016.N().S(`

      data.addRows(`)
//line linechart_material.qtpl:19
	qw422016.N().Z(chart.DataTableJSON())
//line linechart_material.qtpl:19
	qw422016.N().S(`);

      var options = `)
//line linechart_material.qtpl:21
	qw422016.N().Z(chart.OptionsJSON())
//line linechart_material.qtpl:21
	qw422016.N().S(`

      var chart = new google.charts.Line(document.getElementById('`)
//line linechart_material.qtpl:23
	qw422016.N().S(chart.ChartDivOrDefault())
//line linechart_material.qtpl:23
	qw422016.N().S(`'));

      chart.draw(data, google.charts.Line.convertOptions(options));
    }
    </script>
  </body>
</html>
`)
//line linechart_material.qtpl:30
}

//line linechart_material.qtpl:30
func WriteLineChartMaterialPage(qq422016 qtio422016.Writer, chart Chart) {
//line linechart_material.qtpl:30
	qw422016 := qt422016.AcquireWriter(qq422016)
//line linechart_material.qtpl:30
	StreamLineChartMaterialPage(qw422016, chart)
//line linechart_material.qtpl:30
	qt422016.ReleaseWriter(qw422016)
//line linechart_material.qtpl:30
}

//line linechart_material.qtpl:30
func LineChartMaterialPage(chart Chart) string {
//line linechart_material.qtpl:30
	qb422016 := qt422016.AcquireByteBuffer()
//line linechart_material.qtpl:30
	WriteLineChartMaterialPage(qb422016, chart)
//line linechart_material.qtpl:30
	qs422016 := string(qb422016.B)
//line linechart_material.qtpl:30
	qt422016.ReleaseByteBuffer(qb422016)
//line linechart_material.qtpl:30
	return qs422016
//line linechart_material.qtpl:30
}
