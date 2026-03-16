package chartir

// LegendPosition defines legend placement.
type LegendPosition string

const (
	LegendPositionTop    LegendPosition = "top"
	LegendPositionBottom LegendPosition = "bottom"
	LegendPositionLeft   LegendPosition = "left"
	LegendPositionRight  LegendPosition = "right"
)

// LegendPositions returns all valid legend position values.
func LegendPositions() []LegendPosition {
	return []LegendPosition{
		LegendPositionTop,
		LegendPositionBottom,
		LegendPositionLeft,
		LegendPositionRight,
	}
}

// Legend defines legend configuration.
type Legend struct {
	// Show controls legend visibility.
	Show bool `json:"show,omitempty"`

	// Position specifies legend placement.
	Position LegendPosition `json:"position,omitempty"`

	// Items lists specific items to show. If empty, auto-generated from marks.
	Items []string `json:"items,omitempty"`
}
