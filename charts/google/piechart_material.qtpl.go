// Code generated by qtc from "piechart_material.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line piechart_material.qtpl:1
package google

//line piechart_material.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line piechart_material.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line piechart_material.qtpl:1
func StreamPieChartMaterialPage(qw422016 *qt422016.Writer, data PieChartMaterial) {
//line piechart_material.qtpl:1
	qw422016.N().S(`<!DOCTYPE html>
<html>
<head>
  <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
</head>
<body>
  <h1>`)
//line piechart_material.qtpl:7
	qw422016.E().S(data.Title)
//line piechart_material.qtpl:7
	qw422016.N().S(`</h1>
  <div id="`)
//line piechart_material.qtpl:8
	qw422016.N().S(data.ChartDivOrDefault())
//line piechart_material.qtpl:8
	qw422016.N().S(`"></div>
  <script>
      google.charts.load('current', {'packages':['corechart']});
      google.charts.setOnLoadCallback(drawChart);

    function drawChart() {
      var data = google.visualization.arrayToDataTable(`)
//line piechart_material.qtpl:14
	qw422016.N().Z(data.DataMatrixJSON())
//line piechart_material.qtpl:14
	qw422016.N().S(`);

      var options = `)
//line piechart_material.qtpl:16
	qw422016.N().Z(data.GoogleOptions.MustJSON())
//line piechart_material.qtpl:16
	qw422016.N().S(`

      var chart = new google.visualization.PieChart(document.getElementById('`)
//line piechart_material.qtpl:18
	qw422016.N().S(data.ChartDivOrDefault())
//line piechart_material.qtpl:18
	qw422016.N().S(`'));

      chart.draw(data, options);
    }
    </script>
  </body>
</html>
`)
//line piechart_material.qtpl:25
}

//line piechart_material.qtpl:25
func WritePieChartMaterialPage(qq422016 qtio422016.Writer, data PieChartMaterial) {
//line piechart_material.qtpl:25
	qw422016 := qt422016.AcquireWriter(qq422016)
//line piechart_material.qtpl:25
	StreamPieChartMaterialPage(qw422016, data)
//line piechart_material.qtpl:25
	qt422016.ReleaseWriter(qw422016)
//line piechart_material.qtpl:25
}

//line piechart_material.qtpl:25
func PieChartMaterialPage(data PieChartMaterial) string {
//line piechart_material.qtpl:25
	qb422016 := qt422016.AcquireByteBuffer()
//line piechart_material.qtpl:25
	WritePieChartMaterialPage(qb422016, data)
//line piechart_material.qtpl:25
	qs422016 := string(qb422016.B)
//line piechart_material.qtpl:25
	qt422016.ReleaseByteBuffer(qb422016)
//line piechart_material.qtpl:25
	return qs422016
//line piechart_material.qtpl:25
}
