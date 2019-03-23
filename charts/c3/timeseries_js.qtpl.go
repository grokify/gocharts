// This file is automatically generated by qtc from "timeseries_js.qtpl".
// See https://github.com/valyala/quicktemplate for details.

//line timeseries_js.qtpl:1
package c3

//line timeseries_js.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line timeseries_js.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line timeseries_js.qtpl:1
func StreamTimeseriesHTML(qw422016 *qt422016.Writer, data TimeseriesData) {
	//line timeseries_js.qtpl:1
	qw422016.N().S(`
`)
	//line timeseries_js.qtpl:2
	if data.IncludeTitle {
		//line timeseries_js.qtpl:2
		qw422016.N().S(`
<`)
		//line timeseries_js.qtpl:3
		qw422016.N().S(data.TitleLevel)
		//line timeseries_js.qtpl:3
		qw422016.N().S(`>`)
		//line timeseries_js.qtpl:3
		qw422016.N().S(data.Title)
		//line timeseries_js.qtpl:3
		qw422016.N().S(`</`)
		//line timeseries_js.qtpl:3
		qw422016.N().S(data.TitleLevel)
		//line timeseries_js.qtpl:3
		qw422016.N().S(`>
`)
		//line timeseries_js.qtpl:4
	}
	//line timeseries_js.qtpl:4
	qw422016.N().S(`
<div class='c3-editor' id='`)
	//line timeseries_js.qtpl:5
	qw422016.N().S(data.DivID)
	//line timeseries_js.qtpl:5
	qw422016.N().S(`'></div>
<script>

var `)
	//line timeseries_js.qtpl:8
	qw422016.N().S(data.JSDataVar)
	//line timeseries_js.qtpl:8
	qw422016.N().S(` = `)
	//line timeseries_js.qtpl:8
	qw422016.N().Z(data.JSONData.JSON())
	//line timeseries_js.qtpl:8
	qw422016.N().S(`
var `)
	//line timeseries_js.qtpl:9
	qw422016.N().S(data.JSChartVar)
	//line timeseries_js.qtpl:9
	qw422016.N().S(` = c3.generate({
    bindto: '#`)
	//line timeseries_js.qtpl:10
	qw422016.N().S(data.DivID)
	//line timeseries_js.qtpl:10
	qw422016.N().S(`',
    data: {
        x: 'x',
        columns: `)
	//line timeseries_js.qtpl:13
	qw422016.N().S(data.JSDataVar)
	//line timeseries_js.qtpl:13
	qw422016.N().S(`.columns
    },
    axis: {
        x: {
            type: 'timeseries',
            tick: {
                format: '%Y-%m'
            }
        }
    },
    /*
    tooltip: {
        format: {
            title: function (d) {
                var iso8601ym = JSON.stringify(d).substr(1,7)
                return iso8601ym + ' ' + `)
	//line timeseries_js.qtpl:28
	qw422016.N().S(data.JSDataVar)
	//line timeseries_js.qtpl:28
	qw422016.N().S(`.totalsMap[iso8601ym]
            }
        }
    }*/
});
</script>
`)
//line timeseries_js.qtpl:34
}

//line timeseries_js.qtpl:34
func WriteTimeseriesHTML(qq422016 qtio422016.Writer, data TimeseriesData) {
	//line timeseries_js.qtpl:34
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line timeseries_js.qtpl:34
	StreamTimeseriesHTML(qw422016, data)
	//line timeseries_js.qtpl:34
	qt422016.ReleaseWriter(qw422016)
//line timeseries_js.qtpl:34
}

//line timeseries_js.qtpl:34
func TimeseriesHTML(data TimeseriesData) string {
	//line timeseries_js.qtpl:34
	qb422016 := qt422016.AcquireByteBuffer()
	//line timeseries_js.qtpl:34
	WriteTimeseriesHTML(qb422016, data)
	//line timeseries_js.qtpl:34
	qs422016 := string(qb422016.B)
	//line timeseries_js.qtpl:34
	qt422016.ReleaseByteBuffer(qb422016)
	//line timeseries_js.qtpl:34
	return qs422016
//line timeseries_js.qtpl:34
}
