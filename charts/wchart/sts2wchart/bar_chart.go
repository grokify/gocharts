package sts2wchart

import (
	"fmt"
	"strconv"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/grokify/mogo/math/mathutil"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/type/maputil"

	"github.com/grokify/gocharts/v2/charts/wchart"
	"github.com/grokify/gocharts/v2/data/timeseries"
)

func TimeSeriesToBarChart(ds timeseries.TimeSeries) chartdraw.BarChart {
	graph := chartdraw.BarChart{
		Title: ds.SeriesName,
		Background: chartdraw.Style{
			Padding: chartdraw.Box{
				Top: 40,
			},
		},
		YAxis: chartdraw.YAxis{
			ValueFormatter: func(v any) string {
				if vf, isFloat := v.(float64); isFloat {
					return strconvutil.Commify(int64(vf))
				}
				return ""
			},
			Ticks: []chartdraw.Tick{},
		},
		ColorPalette: wchart.ColorsDefault(),
		Height:       512,
		BarWidth:     20,
		Bars:         []chartdraw.Value{},
	}
	highValue := int64(0)
	lowValue := int64(0)

	items := ds.ItemsSorted()
	i := 0
	for _, item := range items {
		graph.Bars = append(
			graph.Bars,
			chartdraw.Value{
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

func MsiToValues(msi maputil.MapStringInt, inclValueInKey bool) []chartdraw.Value {
	values := []chartdraw.Value{}
	keys := msi.Keys(true)
	for _, key := range keys {
		val := msi.MustGet(key, 0)
		if inclValueInKey {
			key += " (" + strconv.Itoa(val) + ")"
		}
		values = append(values,
			chartdraw.Value{
				Label: key,
				Value: float64(val)})
	}
	return values
}
