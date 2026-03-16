package wchart

import (
	"bytes"
	"testing"

	"github.com/grokify/gocharts/v2/charts/chartir"
)

func TestCompileBarChart(t *testing.T) {
	ir := &chartir.ChartIR{
		Title: "Test Bar Chart",
		Datasets: []chartir.Dataset{
			{
				ID: "data",
				Columns: []chartir.Column{
					{Name: "label", Type: chartir.ColumnTypeString},
					{Name: "value", Type: chartir.ColumnTypeNumber},
				},
				Rows: [][]string{
					{"A", "10"},
					{"B", "20"},
					{"C", "15"},
				},
			},
		},
		Marks: []chartir.Mark{
			{
				ID:        "bars",
				DatasetID: "data",
				Geometry:  chartir.GeometryBar,
				Encode: chartir.Encode{
					Y:     "label",
					Value: "value",
				},
			},
		},
	}

	compiler := NewCompiler()
	var buf bytes.Buffer
	err := compiler.RenderPNG(ir, &buf)
	if err != nil {
		t.Fatalf("RenderPNG failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("RenderPNG produced empty output")
	}

	// PNG magic bytes
	if buf.Len() >= 8 {
		magic := buf.Bytes()[:8]
		expected := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		for i := range expected {
			if magic[i] != expected[i] {
				t.Errorf("Invalid PNG magic bytes at position %d", i)
			}
		}
	}
}

func TestCompileScatterChart(t *testing.T) {
	ir := &chartir.ChartIR{
		Title: "Test Scatter Chart",
		Datasets: []chartir.Dataset{
			{
				ID: "data",
				Columns: []chartir.Column{
					{Name: "x", Type: chartir.ColumnTypeNumber},
					{Name: "y", Type: chartir.ColumnTypeNumber},
				},
				Rows: [][]string{
					{"1", "2"},
					{"2", "4"},
					{"3", "3"},
					{"4", "5"},
				},
			},
		},
		Marks: []chartir.Mark{
			{
				ID:        "points",
				DatasetID: "data",
				Geometry:  chartir.GeometryScatter,
				Encode: chartir.Encode{
					X: "x",
					Y: "y",
				},
			},
		},
	}

	compiler := NewCompiler()
	var buf bytes.Buffer
	err := compiler.RenderSVG(ir, &buf)
	if err != nil {
		t.Fatalf("RenderSVG failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("RenderSVG produced empty output")
	}

	// Check for SVG content
	content := buf.String()
	if len(content) > 0 && content[0] != '<' {
		t.Error("RenderSVG output doesn't look like SVG")
	}
}

func TestCompilePieChart(t *testing.T) {
	ir := &chartir.ChartIR{
		Title: "Test Pie Chart",
		Datasets: []chartir.Dataset{
			{
				ID: "data",
				Columns: []chartir.Column{
					{Name: "category", Type: chartir.ColumnTypeString},
					{Name: "value", Type: chartir.ColumnTypeNumber},
				},
				Rows: [][]string{
					{"Red", "30"},
					{"Blue", "50"},
					{"Green", "20"},
				},
			},
		},
		Marks: []chartir.Mark{
			{
				ID:        "slices",
				DatasetID: "data",
				Geometry:  chartir.GeometryPie,
				Encode: chartir.Encode{
					Name:  "category",
					Value: "value",
				},
			},
		},
	}

	compiler := NewCompiler()
	var buf bytes.Buffer
	err := compiler.RenderPNG(ir, &buf)
	if err != nil {
		t.Fatalf("RenderPNG failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("RenderPNG produced empty output")
	}
}
