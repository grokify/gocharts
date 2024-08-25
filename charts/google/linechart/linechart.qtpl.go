// Code generated by qtc from "linechart.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line linechart.qtpl:1
package linechart

//line linechart.qtpl:1
import "github.com/grokify/gocharts/v2/charts/google"

//line linechart.qtpl:2
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line linechart.qtpl:2
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line linechart.qtpl:2
func StreamLineChartPage(qw422016 *qt422016.Writer, chart google.Chart) {
//line linechart.qtpl:2
	qw422016.N().S(`<!DOCTYPE html>
<html>
<head>
  <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
</head>
<body>
  <h1>`)
//line linechart.qtpl:8
	qw422016.N().S(chart.PageTitle())
//line linechart.qtpl:8
	qw422016.N().S(`</h1>
  <div id="`)
//line linechart.qtpl:9
	qw422016.N().S(chart.ChartDivOrDefault())
//line linechart.qtpl:9
	qw422016.N().S(`"></div>
  <script>
      google.charts.load('current', {'packages':['line']});
      google.charts.setOnLoadCallback(drawChart);

    function drawChart() {
      var data = google.visualization.arrayToDataTable(`)
//line linechart.qtpl:15
	qw422016.N().Z(chart.DataTableJSON())
//line linechart.qtpl:15
	qw422016.N().S(`);

      var options = `)
//line linechart.qtpl:17
	qw422016.N().Z(chart.OptionsJSON())
//line linechart.qtpl:17
	qw422016.N().S(`

      var chart = new google.visualization.LineChart(document.getElementById('`)
//line linechart.qtpl:19
	qw422016.N().S(chart.ChartDivOrDefault())
//line linechart.qtpl:19
	qw422016.N().S(`'));

      chart.draw(data, options);
    }
    </script>
  </body>
</html>
`)
//line linechart.qtpl:26
}

//line linechart.qtpl:26
func WriteLineChartPage(qq422016 qtio422016.Writer, chart google.Chart) {
//line linechart.qtpl:26
	qw422016 := qt422016.AcquireWriter(qq422016)
//line linechart.qtpl:26
	StreamLineChartPage(qw422016, chart)
//line linechart.qtpl:26
	qt422016.ReleaseWriter(qw422016)
//line linechart.qtpl:26
}

//line linechart.qtpl:26
func LineChartPage(chart google.Chart) string {
//line linechart.qtpl:26
	qb422016 := qt422016.AcquireByteBuffer()
//line linechart.qtpl:26
	WriteLineChartPage(qb422016, chart)
//line linechart.qtpl:26
	qs422016 := string(qb422016.B)
//line linechart.qtpl:26
	qt422016.ReleaseByteBuffer(qb422016)
//line linechart.qtpl:26
	return qs422016
//line linechart.qtpl:26
}
