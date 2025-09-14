package wchart

import (
	"errors"
	"io"

	"github.com/grokify/gocharts/v2/data/tablef64"
	chart "github.com/go-analyze/charts/chartdraw"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Table struct {
	Table *tablef64.Table
}

func (tbl *Table) BoxPlotColNameWritePNGFile(filename, title, colName string) error {
	vals, err := tbl.Table.ValuesColumnName(colName)
	if err != nil {
		return err
	}
	boxPlot, err := tbl.BoxPlot(vals)
	if err != nil {
		return err
	}
	p := plot.New()
	if len(title) > 0 {
		p.Title.Text = title
	}
	p.Add(boxPlot)

	return p.Save(3*vg.Inch, 3*vg.Inch, filename)
}

func (tbl *Table) BoxPlot(values []float64) (*plotter.BoxPlot, error) {
	return plotter.NewBoxPlot(vg.Length(15), 0.0, plotter.Values(values))
}

/*
func boxPlotWritePNG(filename, title string, values []float64) error {
	p := plot.New()
	if len(title) > 0 {
		p.Title.Text = title
	}
	box, err := plotter.NewBoxPlot(vg.Length(15), 0.0, plotter.Values(values))
	if err != nil {
		return err
	}
	p.Add(box)

	return p.Save(3*vg.Inch, 3*vg.Inch, filename)
*/

func (tbl *Table) ScatterPlotWritePNG(w io.Writer, title, xColName, yColName string, dotWidth float64) error {
	if graph, err := tbl.ScatterPlot(title, xColName, yColName, dotWidth); err != nil {
		return err
	} else {
		return WritePNG(w, graph)
	}
}

func (tbl *Table) ScatterPlotWritePNGFile(filename, title, xColName, yColName string, dotWidth float64) error {
	if graph, err := tbl.ScatterPlot(title, xColName, yColName, dotWidth); err != nil {
		return err
	} else {
		return WritePNGFile(filename, graph)
	}
}

func (tbl *Table) ScatterPlot(title, xColName, yColName string, dotWidth float64) (chart.Chart, error) {
	if dotWidth < 0 {
		dotWidth = 3
	}
	if tbl.Table == nil {
		return chart.Chart{}, errors.New("table cannot be nil")
	}
	pts, err := tbl.Table.PointXYsColumnNames(xColName, yColName)
	if err != nil {
		return chart.Chart{}, err
	}
	xs, ys := pts.XAndYSeries()
	graph := chart.Chart{
		Title: title,
		XAxis: chart.XAxis{
			Name: xColName,
		},
		YAxis: chart.YAxis{
			Name: yColName,
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeWidth: chart.Disabled,
					DotWidth:    dotWidth,
				},
				XValues: xs,
				YValues: ys,
			},
		},
	}
	return graph, nil
}
