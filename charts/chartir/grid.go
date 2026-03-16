package chartir

// Grid defines the chart container/grid positioning.
// Values can be percentages (e.g., "10%") or pixel values (e.g., "50").
type Grid struct {
	// Left is the distance from the left edge.
	Left string `json:"left,omitempty"`

	// Right is the distance from the right edge.
	Right string `json:"right,omitempty"`

	// Top is the distance from the top edge.
	Top string `json:"top,omitempty"`

	// Bottom is the distance from the bottom edge.
	Bottom string `json:"bottom,omitempty"`

	// Width is the grid width.
	Width string `json:"width,omitempty"`

	// Height is the grid height.
	Height string `json:"height,omitempty"`

	// ContainLabel adjusts grid to contain axis labels.
	ContainLabel bool `json:"containLabel,omitempty"`
}
