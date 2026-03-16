// Package wchart compiles ChartIR to go-analyze/charts for PNG/SVG rendering.
package wchart

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
	"github.com/grokify/gocharts/v2/charts/chartir"
	"github.com/grokify/gocharts/v2/charts/wchart"
)

// ChartType is the interface for renderable charts.
type ChartType interface {
	Render(rp chartdraw.RendererProvider, w io.Writer) error
}

// Compiler converts ChartIR to go-analyze/charts types.
type Compiler struct {
	// Width is the chart width in pixels. Default: 600.
	Width int

	// Height is the chart height in pixels. Default: 400.
	Height int

	// ColorPalette is the color palette to use. Default: wchart.ColorsDefault().
	ColorPalette chartdraw.ColorPalette
}

// NewCompiler creates a new compiler with default settings.
func NewCompiler() *Compiler {
	return &Compiler{
		Width:        600,
		Height:       400,
		ColorPalette: wchart.ColorsDefault(),
	}
}

// RenderPNG renders the ChartIR to PNG format.
func (c *Compiler) RenderPNG(ir *chartir.ChartIR, w io.Writer) error {
	chart, err := c.Compile(ir)
	if err != nil {
		return err
	}
	return chart.Render(chartdraw.PNG, w)
}

// RenderSVG renders the ChartIR to SVG format.
func (c *Compiler) RenderSVG(ir *chartir.ChartIR, w io.Writer) error {
	chart, err := c.Compile(ir)
	if err != nil {
		return err
	}
	return chart.Render(chartdraw.SVG, w)
}

// Compile converts ChartIR to a ChartType.
func (c *Compiler) Compile(ir *chartir.ChartIR) (ChartType, error) {
	if len(ir.Marks) == 0 {
		return nil, fmt.Errorf("chartir: no marks defined")
	}

	// Determine chart type from first mark's geometry
	geometry := ir.Marks[0].Geometry

	switch geometry {
	case chartir.GeometryBar:
		return c.compileBarChart(ir)
	case chartir.GeometryLine, chartir.GeometryArea:
		return c.compileLineChart(ir)
	case chartir.GeometryScatter:
		return c.compileScatterChart(ir)
	case chartir.GeometryPie:
		return c.compilePieChart(ir)
	default:
		return nil, fmt.Errorf("chartir: unsupported geometry: %s", geometry)
	}
}

func (c *Compiler) compileBarChart(ir *chartir.ChartIR) (*chartdraw.BarChart, error) {
	mark := ir.Marks[0]
	dataset := ir.GetDataset(mark.DatasetID)
	if dataset == nil {
		return nil, fmt.Errorf("chartir: dataset not found: %s", mark.DatasetID)
	}

	// Get label and value columns
	labelCol := mark.Encode.Y // For horizontal bar charts
	if labelCol == "" {
		labelCol = mark.Encode.X
	}
	valueCol := mark.Encode.X
	if mark.Encode.Y != "" && mark.Encode.X != "" {
		valueCol = mark.Encode.X
		labelCol = mark.Encode.Y
	}
	if mark.Encode.Value != "" {
		valueCol = mark.Encode.Value
	}

	labels := dataset.GetStringValues(labelCol)
	values := dataset.GetFloat64Values(valueCol)

	if len(labels) != len(values) {
		return nil, fmt.Errorf("chartir: label/value count mismatch")
	}

	bars := make([]chartdraw.Value, len(labels))
	for i := range labels {
		bars[i] = chartdraw.Value{
			Label: labels[i],
			Value: values[i],
		}
	}

	chart := &chartdraw.BarChart{
		Title:        ir.Title,
		Height:       c.Height,
		Width:        c.Width,
		ColorPalette: c.ColorPalette,
		Bars:         bars,
		Background: chartdraw.Style{
			Padding: chartdraw.Box{Top: 40},
		},
	}

	return chart, nil
}

func (c *Compiler) compileLineChart(ir *chartir.ChartIR) (*chartdraw.Chart, error) {
	chart := &chartdraw.Chart{
		Title:        ir.Title,
		Height:       c.Height,
		Width:        c.Width,
		ColorPalette: c.ColorPalette,
	}

	// Add series for each mark
	for _, mark := range ir.Marks {
		dataset := ir.GetDataset(mark.DatasetID)
		if dataset == nil {
			return nil, fmt.Errorf("chartir: dataset not found: %s", mark.DatasetID)
		}

		xValues := dataset.GetFloat64Values(mark.Encode.X)
		yValues := dataset.GetFloat64Values(mark.Encode.Y)

		if len(xValues) != len(yValues) {
			return nil, fmt.Errorf("chartir: x/y value count mismatch")
		}

		series := chartdraw.ContinuousSeries{
			Name:    mark.Name,
			XValues: xValues,
			YValues: yValues,
		}

		if mark.Style != nil && mark.Style.Color != "" {
			series.Style = chartdraw.Style{
				StrokeColor: colorFromHex(mark.Style.Color),
			}
		}

		chart.Series = append(chart.Series, series)
	}

	// Configure axes
	if xAxis := ir.GetXAxis(); xAxis != nil {
		chart.XAxis = chartdraw.XAxis{
			Name: xAxis.Name,
		}
	}
	if yAxis := ir.GetYAxis(); yAxis != nil {
		chart.YAxis = chartdraw.YAxis{
			Name: yAxis.Name,
		}
	}

	return chart, nil
}

func (c *Compiler) compileScatterChart(ir *chartir.ChartIR) (*chartdraw.Chart, error) {
	chart := &chartdraw.Chart{
		Title:        ir.Title,
		Height:       c.Height,
		Width:        c.Width,
		ColorPalette: c.ColorPalette,
	}

	// Add series for each mark
	for i, mark := range ir.Marks {
		dataset := ir.GetDataset(mark.DatasetID)
		if dataset == nil {
			return nil, fmt.Errorf("chartir: dataset not found: %s", mark.DatasetID)
		}

		xValues := dataset.GetFloat64Values(mark.Encode.X)
		yValues := dataset.GetFloat64Values(mark.Encode.Y)

		if len(xValues) != len(yValues) {
			return nil, fmt.Errorf("chartir: x/y value count mismatch")
		}

		// Use scatter series (continuous series with dot style)
		color := c.ColorPalette.GetSeriesColor(i)
		if mark.Style != nil && mark.Style.Color != "" {
			color = colorFromHex(mark.Style.Color)
		}

		dotWidth := 5.0
		if mark.Style != nil && mark.Style.SymbolSize != nil {
			dotWidth = *mark.Style.SymbolSize
		}

		series := chartdraw.ContinuousSeries{
			Name:    mark.Name,
			XValues: xValues,
			YValues: yValues,
			Style: chartdraw.Style{
				StrokeWidth: 0, // No line
				DotWidth:    dotWidth,
				DotColor:    color,
			},
		}

		chart.Series = append(chart.Series, series)
	}

	// Configure axes
	if xAxis := ir.GetXAxis(); xAxis != nil {
		chart.XAxis = chartdraw.XAxis{
			Name: xAxis.Name,
		}
	}
	if yAxis := ir.GetYAxis(); yAxis != nil {
		chart.YAxis = chartdraw.YAxis{
			Name: yAxis.Name,
		}
	}

	return chart, nil
}

// colorFromHex parses a hex color string like "#ff0000" or "ff0000".
func colorFromHex(hex string) drawing.Color {
	hex = strings.TrimPrefix(hex, "#")
	return drawing.ColorFromHex(hex)
}

func (c *Compiler) compilePieChart(ir *chartir.ChartIR) (*chartdraw.PieChart, error) {
	mark := ir.Marks[0]
	dataset := ir.GetDataset(mark.DatasetID)
	if dataset == nil {
		return nil, fmt.Errorf("chartir: dataset not found: %s", mark.DatasetID)
	}

	nameCol := mark.Encode.Name
	if nameCol == "" {
		nameCol = mark.Encode.Category
	}
	valueCol := mark.Encode.Value

	names := dataset.GetStringValues(nameCol)
	values := dataset.GetFloat64Values(valueCol)

	if len(names) != len(values) {
		return nil, fmt.Errorf("chartir: name/value count mismatch")
	}

	pieValues := make([]chartdraw.Value, len(names))
	for i := range names {
		pieValues[i] = chartdraw.Value{
			Label: names[i],
			Value: values[i],
		}
	}

	chart := &chartdraw.PieChart{
		Title:        ir.Title,
		Height:       c.Height,
		Width:        c.Width,
		ColorPalette: c.ColorPalette,
		Values:       pieValues,
	}

	return chart, nil
}
