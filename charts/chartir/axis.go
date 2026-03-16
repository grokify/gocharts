package chartir

// AxisType defines the axis scale type.
type AxisType string

const (
	AxisTypeCategory AxisType = "category"
	AxisTypeValue    AxisType = "value"
	AxisTypeTime     AxisType = "time"
	AxisTypeLog      AxisType = "log"
)

// AxisTypes returns all valid axis type values.
func AxisTypes() []AxisType {
	return []AxisType{
		AxisTypeCategory,
		AxisTypeValue,
		AxisTypeTime,
		AxisTypeLog,
	}
}

// AxisPosition defines axis placement.
type AxisPosition string

const (
	AxisPositionBottom AxisPosition = "bottom"
	AxisPositionTop    AxisPosition = "top"
	AxisPositionLeft   AxisPosition = "left"
	AxisPositionRight  AxisPosition = "right"
)

// AxisPositions returns all valid axis position values.
func AxisPositions() []AxisPosition {
	return []AxisPosition{
		AxisPositionBottom,
		AxisPositionTop,
		AxisPositionLeft,
		AxisPositionRight,
	}
}

// Axis defines a chart axis.
type Axis struct {
	// ID uniquely identifies this axis.
	ID string `json:"id"`

	// Type specifies the axis scale type.
	Type AxisType `json:"type"`

	// Position specifies where the axis is placed.
	Position AxisPosition `json:"position"`

	// Name is the axis label/title.
	Name string `json:"name,omitempty"`

	// Min is the minimum axis value. If nil, auto-calculated.
	Min *float64 `json:"min,omitempty"`

	// Max is the maximum axis value. If nil, auto-calculated.
	Max *float64 `json:"max,omitempty"`
}

// IsHorizontal returns true if the axis position is horizontal (top/bottom).
func (a Axis) IsHorizontal() bool {
	return a.Position == AxisPositionBottom || a.Position == AxisPositionTop
}

// IsVertical returns true if the axis position is vertical (left/right).
func (a Axis) IsVertical() bool {
	return a.Position == AxisPositionLeft || a.Position == AxisPositionRight
}
