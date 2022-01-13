package main

//go:generate go run main.go

import (
	"fmt"
	"os"

	"github.com/wcharczuk/go-chart/v2"
)

func main() {
	// In this example we add a second series, and assign it to the secondary y axis, giving that series it's own range.
	// We also enable all of the axes by setting the `Show` propery of their respective styles to `true`.

	graph := chart.Chart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 25,
			},
		},
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				return fmt.Sprintf("%v", typed)
				//			m := timeutil.MonthBegin(time.Now(), int(typed))
				//			return m.Format(timeutil.RFC3339FullDate)
				//typed := v.(float64)
				//typedDate := chart.TimeFromFloat64(typed)
				//return fmt.Sprintf("%d-%d\n%d", typedDate.Month(), typedDate.Day(), typedDate.Year())
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "A",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
			chart.ContinuousSeries{
				Name: "B",
				//YAxis:   chart.YAxisSecondary,
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{50.0, 40.0, 30.0, 20.0, 10.0},
			},
		},
	}

	//note we have to do this as a separate step because we need a reference to graph
	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}
