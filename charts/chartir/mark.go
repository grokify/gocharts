package chartir

// Geometry defines the visual representation type.
// This replaces polymorphic series.type with a simple enum.
type Geometry string

const (
	GeometryLine    Geometry = "line"
	GeometryBar     Geometry = "bar"
	GeometryPie     Geometry = "pie"
	GeometryScatter Geometry = "scatter"
	GeometryArea    Geometry = "area"
	GeometryRadar   Geometry = "radar"
	GeometryFunnel  Geometry = "funnel"
	GeometryGauge   Geometry = "gauge"
	GeometryHeatmap Geometry = "heatmap"
	GeometryTreemap Geometry = "treemap"
	GeometrySankey  Geometry = "sankey"
)

// Geometries returns all valid geometry values.
func Geometries() []Geometry {
	return []Geometry{
		GeometryLine,
		GeometryBar,
		GeometryPie,
		GeometryScatter,
		GeometryArea,
		GeometryRadar,
		GeometryFunnel,
		GeometryGauge,
		GeometryHeatmap,
		GeometryTreemap,
		GeometrySankey,
	}
}

// CoordinateSystem defines the coordinate system type.
type CoordinateSystem string

const (
	CoordinateCartesian2D CoordinateSystem = "cartesian2d"
	CoordinatePolar       CoordinateSystem = "polar"
	CoordinateRadial      CoordinateSystem = "radial"
)

// CoordinateSystems returns all valid coordinate system values.
func CoordinateSystems() []CoordinateSystem {
	return []CoordinateSystem{
		CoordinateCartesian2D,
		CoordinatePolar,
		CoordinateRadial,
	}
}

// Mark defines a visual mark (equivalent to chart series).
// All marks have the same structure regardless of geometry type.
// The compiler handles geometry-specific transformations.
type Mark struct {
	// ID uniquely identifies this mark.
	ID string `json:"id"`

	// DatasetID references the dataset to use for this mark.
	DatasetID string `json:"datasetId"`

	// Geometry specifies the visual representation type.
	Geometry Geometry `json:"geometry"`

	// CoordinateSystem specifies the coordinate system.
	// Defaults to cartesian2d if not specified.
	CoordinateSystem CoordinateSystem `json:"coordinateSystem,omitempty"`

	// Encode maps data columns to visual channels.
	Encode Encode `json:"encode"`

	// Style defines visual styling properties.
	Style *Style `json:"style,omitempty"`

	// Stack groups marks for stacking. Marks with the same stack
	// value are stacked together.
	Stack string `json:"stack,omitempty"`

	// Smooth enables smooth curves for line/area geometries.
	Smooth bool `json:"smooth,omitempty"`

	// Name is the display name for this mark in legends/tooltips.
	Name string `json:"name,omitempty"`
}

// Encode maps data columns to visual channels.
// Not all fields are used by all geometry types.
type Encode struct {
	// X maps to the x-axis (for Cartesian geometries).
	X string `json:"x,omitempty"`

	// Y maps to the y-axis (for Cartesian geometries).
	Y string `json:"y,omitempty"`

	// Value maps to the primary value (for pie, funnel, gauge, treemap).
	Value string `json:"value,omitempty"`

	// Name maps to the name/label (for pie chart segments, etc.).
	Name string `json:"name,omitempty"`

	// Size maps to mark size (for scatter plots).
	Size string `json:"size,omitempty"`

	// Color maps to mark color (for color encoding by data).
	Color string `json:"color,omitempty"`

	// Category maps to categorical grouping (for pie, funnel, treemap).
	Category string `json:"category,omitempty"`

	// Indicator maps to radar indicator (for radar charts).
	Indicator string `json:"indicator,omitempty"`

	// Source maps to source node (for sankey diagrams).
	Source string `json:"source,omitempty"`

	// Target maps to target node (for sankey diagrams).
	Target string `json:"target,omitempty"`

	// Heat maps to heat intensity (for heatmaps).
	Heat string `json:"heat,omitempty"`
}
