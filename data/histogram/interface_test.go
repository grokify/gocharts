package histogram

import (
	"testing"
)

var histogramAnyTests = []struct {
	v   HistogramAny
	sum int
}{
	{&Histogram{Bins: map[string]int{"foo": 10}}, 10},
	{&HistogramSet{
		HistogramMap: map[string]*Histogram{
			"bar": {Bins: map[string]int{"foo": 10}}}}, 10},
	{
		&HistogramSets{
			HistogramSetMap: map[string]*HistogramSet{
				"baz": {
					HistogramMap: map[string]*Histogram{
						"bar": {Bins: map[string]int{"foo": 10}}}},
			},
		},
		10,
	},
	{&Histogram{
		Bins: map[string]int{"foo": 10, "bar": 2}}, 12},
	{&HistogramSet{
		HistogramMap: map[string]*Histogram{
			"bar": {Bins: map[string]int{
				"foo": 10,
				"bar": 20,
			}},
			"baz": {Bins: map[string]int{
				"foo": 12,
				"bar": 11,
			}}},
	}, 53,
	},
}

// TestHistogramAny tests interface functions for `HistogramAny` interface.`
func TestHistogramAny(t *testing.T) {
	for _, tt := range histogramAnyTests {
		sum := tt.v.Sum()
		if sum != tt.sum {
			t.Errorf("histogramAny.Sum() mismatch: want (%d) got (%d)",
				tt.sum, sum)
		}
	}
}
