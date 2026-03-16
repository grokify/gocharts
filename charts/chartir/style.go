package chartir

// Style defines visual styling properties.
// The structure is intentionally flat and simple to avoid
// deep nesting that would require polymorphic handling.
type Style struct {
	// Color is the primary color (hex, rgb, or named color).
	Color string `json:"color,omitempty"`

	// Opacity is the transparency level (0.0 to 1.0).
	Opacity *float64 `json:"opacity,omitempty"`

	// BorderColor is the border/stroke color.
	BorderColor string `json:"borderColor,omitempty"`

	// BorderWidth is the border/stroke width in pixels.
	BorderWidth *float64 `json:"borderWidth,omitempty"`

	// Line/Area specific

	// Smooth enables smooth curves for line/area geometries.
	Smooth bool `json:"smooth,omitempty"`

	// AreaOpacity sets opacity for area fill (0.0 to 1.0).
	AreaOpacity *float64 `json:"areaOpacity,omitempty"`

	// LineWidth sets line stroke width in pixels.
	LineWidth *float64 `json:"lineWidth,omitempty"`

	// Bar specific

	// BarWidth sets bar width (number or percentage string).
	BarWidth any `json:"barWidth,omitempty"`

	// BarGap sets gap between bars (percentage string).
	BarGap string `json:"barGap,omitempty"`

	// BorderRadius sets bar corner radius.
	BorderRadius any `json:"borderRadius,omitempty"`

	// Point/Symbol specific

	// Symbol sets the marker symbol type.
	Symbol string `json:"symbol,omitempty"`

	// SymbolSize sets the marker symbol size.
	SymbolSize *float64 `json:"symbolSize,omitempty"`

	// Radar specific

	// Shape sets radar chart shape (polygon or circle).
	Shape string `json:"shape,omitempty"`

	// Funnel specific

	// FunnelAlign sets funnel alignment (left, center, right).
	FunnelAlign string `json:"funnelAlign,omitempty"`

	// FunnelSort sets funnel sort direction (ascending, descending, none).
	FunnelSort string `json:"funnelSort,omitempty"`

	// FunnelGap sets gap between funnel segments.
	FunnelGap *float64 `json:"funnelGap,omitempty"`

	// Gauge specific

	// StartAngle sets gauge start angle in degrees.
	StartAngle *float64 `json:"startAngle,omitempty"`

	// EndAngle sets gauge end angle in degrees.
	EndAngle *float64 `json:"endAngle,omitempty"`

	// GaugeMin sets gauge minimum value.
	GaugeMin *float64 `json:"gaugeMin,omitempty"`

	// GaugeMax sets gauge maximum value.
	GaugeMax *float64 `json:"gaugeMax,omitempty"`
}
