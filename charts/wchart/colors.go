package wchart

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"golang.org/x/image/colornames"
)

func MustGetSVGColor(colorName string) drawing.Color {
	drawingColor, err := GetSVGColor(colorName)
	if err != nil {
		return drawing.ColorBlack
	}
	return drawingColor
}

func GetSVGColor(colorName string) (drawing.Color, error) {
	colorName = strings.ToLower(strings.TrimSpace(colorName))
	if col, ok := colornames.Map[colorName]; ok {
		return ColorImageToDrawing(col), nil
	}
	return drawing.ColorBlack, fmt.Errorf("E_COLOR_NOT_FOUND [%s]", colorName)
}

func ColorImageToDrawing(col color.RGBA) drawing.Color {
	return drawing.Color{R: col.R, G: col.G, B: col.B, A: col.A}
}

var (
	// ColorOrange is orange.
	ColorOrange = drawing.Color{R: 255, G: 165, B: 0, A: 255}

	// ColorGreeen is greeen.
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
