// Package chartir provides a non-polymorphic intermediate representation
// for chart configurations. The IR is designed to be AI-friendly,
// easily validated via JSON Schema, and compiled to various chart formats
// including Apache ECharts, go-analyze/charts (PNG/SVG), and others.
package chartir

// ChartIR is the top-level chart intermediate representation.
// It provides a normalized, non-polymorphic structure that can be
// compiled to various chart output formats.
type ChartIR struct {
	// Title is the chart title text.
	Title string `json:"title,omitempty"`

	// Datasets contains the data sources for the chart.
	// Each dataset is referenced by marks via DatasetID.
	Datasets []Dataset `json:"datasets"`

	// Marks define the visual representations (equivalent to chart series).
	// All marks have the same structure regardless of geometry type.
	Marks []Mark `json:"marks"`

	// Axes define the chart axes. Optional for non-Cartesian charts (e.g., pie).
	Axes []Axis `json:"axes,omitempty"`

	// Legend configures the chart legend.
	Legend *Legend `json:"legend,omitempty"`

	// Tooltip configures hover tooltips.
	Tooltip *Tooltip `json:"tooltip,omitempty"`

	// Grid configures the chart container/grid positioning.
	Grid *Grid `json:"grid,omitempty"`
}

// GetDataset returns the dataset with the given ID, or nil if not found.
func (c *ChartIR) GetDataset(id string) *Dataset {
	for i := range c.Datasets {
		if c.Datasets[i].ID == id {
			return &c.Datasets[i]
		}
	}
	return nil
}

// GetXAxis returns the first horizontal axis, or nil if none.
func (c *ChartIR) GetXAxis() *Axis {
	for i := range c.Axes {
		if c.Axes[i].IsHorizontal() {
			return &c.Axes[i]
		}
	}
	return nil
}

// GetYAxis returns the first vertical axis, or nil if none.
func (c *ChartIR) GetYAxis() *Axis {
	for i := range c.Axes {
		if c.Axes[i].IsVertical() {
			return &c.Axes[i]
		}
	}
	return nil
}
