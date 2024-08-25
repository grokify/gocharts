// interfacetests is a separate package for import purposes.
package interfacetests

import (
	"testing"

	"github.com/grokify/gocharts/v2/charts/google"
	"github.com/grokify/gocharts/v2/charts/google/barchart"
	"github.com/grokify/gocharts/v2/charts/google/linechart"
	"github.com/grokify/gocharts/v2/charts/google/piechart"
)

var chartInterfaceTests = []struct {
	v     google.Chart
	title string
}{
	{&barchart.Chart{Title: "foobar"}, "foobar"},
	{&linechart.Chart{Title: "foobar"}, "foobar"},
	{&piechart.Chart{Title: "foobar"}, "foobar"},
}

// TestChartInterface tests interface functions for `HistogramAny` interface.`
func TestChartInterface(t *testing.T) {
	for _, tt := range chartInterfaceTests {
		if tt.v.PageTitle() != tt.title {
			t.Errorf("barchart.TestChartInterface() mismatch: want (%s) got (%s)",
				tt.title, tt.v.PageTitle())
		}
	}
}
