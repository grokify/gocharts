package main

//go:generate go run main.go

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {
	// In this example we add a second series, and assign it to the secondary y axis, giving that series it's own range.
	// We also enable all of the axes by setting the `Show` propery of their respective styles to `true`.

	graph := chartdraw.Chart{
		Background: chartdraw.Style{
			Padding: chartdraw.Box{
				Top:  20,
				Left: 25,
			},
		},
		XAxis: chartdraw.XAxis{
			TickPosition: chartdraw.TickPositionBetweenTicks,
			ValueFormatter: func(v any) string {
				typed := v.(float64)
				return fmt.Sprintf("%v", typed)
				//			m := timeutil.MonthBegin(time.Now(), int(typed))
				//			return m.Format(timeutil.RFC3339FullDate)
				//typed := v.(float64)
				//typedDate := chartdraw.TimeFromFloat64(typed)
				//return fmt.Sprintf("%d-%d\n%d", typedDate.Month(), typedDate.Day(), typedDate.Year())
			},
		},
		Series: []chartdraw.Series{
			chartdraw.ContinuousSeries{
				Name:    "A",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
			chartdraw.ContinuousSeries{
				Name: "B",
				//YAxis:   chartdraw.YAxisSecondary,
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{50.0, 40.0, 30.0, 20.0, 10.0},
			},
		},
	}

	//note we have to do this as a separate step because we need a reference to graph
	graph.Elements = []chartdraw.Renderable{
		chartdraw.LegendLeft(&graph),
	}

	f, err := os.Create("output.png")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer f.Close()
	if err := graph.Render(chartdraw.PNG, f); err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}
	os.Exit(0)
}
