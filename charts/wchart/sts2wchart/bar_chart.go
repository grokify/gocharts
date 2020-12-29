package sts2wchart

import (
	"fmt"
	"strconv"

	"github.com/grokify/gocharts/charts/wchart"
	"github.com/grokify/gocharts/data/statictimeseries"

	"github.com/grokify/simplego/math/mathutil"
	"github.com/grokify/simplego/strconv/strconvutil"
	"github.com/grokify/simplego/type/maputil"
	"github.com/wcharczuk/go-chart"
)

func DataSeriesToBarChart(ds statictimeseries.DataSeries) chart.BarChart {
	graph := chart.BarChart{
		Title: ds.SeriesName,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		YAxis: chart.YAxis{
			ValueFormatter: func(v interface{}) string {
				if vf, isFloat := v.(float64); isFloat {
					return strconvutil.Commify(int64(vf))
				}
				return ""
			},
			Ticks: []chart.Tick{},
		},
		ColorPalette: wchart.ColorsDefault(),
		Height:       512,
		BarWidth:     20,
		Bars:         []chart.Value{},
	}
	highValue := int64(0)
	lowValue := int64(0)

	items := ds.ItemsSorted()
	i := 0
	for _, item := range items {
		graph.Bars = append(
			graph.Bars,
			chart.Value{
				Value: float64(item.Value),
				Label: fmt.Sprintf("%s %s",
					item.Time.Format("Jan '06"),
					strconvutil.Int64Abbreviation(item.Value)),
			})
		if i == 0 {
			highValue = item.Value
			lowValue = item.Value
		} else {
			if item.Value > highValue {
				highValue = item.Value
			}
			if item.Value < lowValue {
				lowValue = item.Value
			}
		}
		i++
	}

	tickValues := mathutil.PrettyTicks(10.0, lowValue, highValue)
	graph.YAxis.Ticks = wchart.TicksInt64(tickValues, strconvutil.Int64Abbreviation)
	return graph
}

func MsiToValues(msi maputil.MapStringInt, inclValueInKey bool) []chart.Value {
	values := []chart.Value{}
	keys := msi.Keys(true)
	for _, key := range keys {
		val := msi.MustGet(key, 0)
		if inclValueInKey {
			key += " (" + strconv.Itoa(val) + ")"
		}
		values = append(values,
			chart.Value{
				Label: key,
				Value: float64(val)})
	}
	return values
}
