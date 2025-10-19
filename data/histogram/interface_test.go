package histogram

import (
	"testing"
)

var histogramAnyTests = []struct {
	v   HistogramAny
	sum int
}{
	{&Histogram{Items: map[string]int{"foo": 10}}, 10},
	{&HistogramSet{
		Items: map[string]*Histogram{
			"bar": {Items: map[string]int{"foo": 10}}}}, 10},
	{
		&HistogramSets{
			Items: map[string]*HistogramSet{
				"baz": {
					Items: map[string]*Histogram{
						"bar": {Items: map[string]int{"foo": 10}}}},
			},
		},
		10,
	},
	{&Histogram{
		Items: map[string]int{"foo": 10, "bar": 2}}, 12},
	{&HistogramSet{
		Items: map[string]*Histogram{
			"bar": {Items: map[string]int{
				"foo": 10,
				"bar": 20,
			}},
			"baz": {Items: map[string]int{
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
