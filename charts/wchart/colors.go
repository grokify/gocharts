package wchart

import (
	"image/color"

	"github.com/grokify/mogo/image/colors"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

// MustParseColor returns a `drawing.Color` for the
// SVG color name. If the color name is not found,
// black is returned.
func MustParseColor(colorName string) drawing.Color {
	drawingColor, err := ParseColor(colorName)
	if err != nil {
		return drawing.ColorBlack
	}
	return drawingColor
}

// ParseColor returns a `drawing.Color` for the
// SVG color name. If the color name is not found,
// an error is returned.
func ParseColor(colorName string) (drawing.Color, error) {
	rgba, err := colors.Parse(colorName)
	if err != nil {
		return drawing.ColorBlack, err
	}
	return ColorImageToDrawing(rgba), nil
}

// ColorImageToDrawing converts a `color.RGBA` value
// to a `drawing.Color` value.
func ColorImageToDrawing(col color.RGBA) drawing.Color {
	return drawing.Color{R: col.R, G: col.G, B: col.B, A: col.A}
}

var (
	// ColorOrange is orange.
	ColorOrange = drawing.Color{R: 255, G: 165, B: 0, A: 255}

	// ColorGreen is green.
	ColorGreen = drawing.Color{R: 0, G: 255, B: 0, A: 255}
)

type Colors struct {
	BackgroundColorVal       drawing.Color
	BackgroundStrokeColorVal drawing.Color
	CanvasColorVal           drawing.Color
	CanvasStrokeColorVal     drawing.Color
	AxisStrokeColorVal       drawing.Color
	TextColorVal             drawing.Color
	SeriesColorVal           drawing.Color
}

func ColorsDefault() Colors {
	return Colors{
		BackgroundColorVal:       chart.DefaultBackgroundColor,
		BackgroundStrokeColorVal: chart.DefaultBackgroundStrokeColor,
		CanvasColorVal:           chart.DefaultCanvasColor,
		CanvasStrokeColorVal:     chart.DefaultCanvasStrokeColor,
		AxisStrokeColorVal:       chart.DefaultAxisColor,
		TextColorVal:             chart.DefaultTextColor,
		SeriesColorVal:           chart.DefaultFillColor}
}

func (c Colors) BackgroundColor() drawing.Color         { return c.BackgroundColorVal }
func (c Colors) BackgroundStrokeColor() drawing.Color   { return c.BackgroundStrokeColorVal }
func (c Colors) CanvasColor() drawing.Color             { return c.CanvasColorVal }
func (c Colors) CanvasStrokeColor() drawing.Color       { return c.CanvasStrokeColorVal }
func (c Colors) AxisStrokeColor() drawing.Color         { return c.AxisStrokeColorVal }
func (c Colors) TextColor() drawing.Color               { return c.TextColorVal }
func (c Colors) GetSeriesColor(index int) drawing.Color { return c.SeriesColorVal }
