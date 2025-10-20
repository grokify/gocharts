package histogram

import (
	"slices"
	"strings"

	"github.com/grokify/mogo/type/maputil"
)

type HistogramSets struct {
	Name  string
	Items map[string]*HistogramSet
	Order []string
}

func NewHistogramSets(name string) *HistogramSets {
	return &HistogramSets{
		Name:  name,
		Items: map[string]*HistogramSet{}}
}

func (hsets *HistogramSets) Add(hsetName, histName, binName string, binCount int, trimSpace bool) {
	if trimSpace {
		hsetName = strings.TrimSpace(hsetName)
		histName = strings.TrimSpace(histName)
		binName = strings.TrimSpace(binName)
	}
	hset, ok := hsets.Items[hsetName]
	if !ok {
		hset = NewHistogramSet(hsetName)
	}
	hset.Add(histName, binName, binCount)
	hsets.Items[hsetName] = hset
}

func (hsets *HistogramSets) Flatten(name string) *HistogramSet {
	hsetFlat := NewHistogramSet(name)
	for _, hset := range hsets.Items {
		for histName, hist := range hset.Items {
			for binName, binCount := range hist.Items {
				hsetFlat.Add(histName, binName, binCount)
			}
		}
	}
	return hsetFlat
}

func (hsets *HistogramSets) BinNames() []string {
	binNamesMap := map[string]int{}
	hsets.Visit(func(hsetName, histName, binName string, binCount int) {
		binNamesMap[binName] = 1
	})
	return maputil.Keys(binNamesMap)
}

// BinValue the value of a bin.
func (hsets *HistogramSets) BinValue(hsetName, histName, binName string) int {
	if hset, ok := hsets.Items[hsetName]; !ok || hset == nil {
		return 0
	} else {
		return hset.BinValue(histName, binName)
	}
}

func (hsets *HistogramSets) BinSumsByHset() *Histogram {
	sums := NewHistogram("Bin Sums")
	for hsetName, hset := range hsets.Items {
		for _, hist := range hset.Items {
			for _, binVal := range hist.Items {
				sums.Add(hsetName, binVal)
			}
		}
	}
	return sums
}

func (hsets *HistogramSets) Counts() *HistogramSetsCounts {
	return NewHistogramSetsCounts(*hsets)
}

// ItemCount returns the number of histogram sets.
func (hsets *HistogramSets) ItemCount() uint {
	return uint(len(hsets.Items))
}

func (hsets *HistogramSets) ItemNames() []string {
	return maputil.Keys(hsets.Items)
}

func (hsets *HistogramSets) Map() map[string]map[string]map[string]int {
	out := map[string]map[string]map[string]int{}
	for hsetName, hset := range hsets.Items {
		if _, ok := out[hsetName]; !ok {
			out[hsetName] = map[string]map[string]int{}
		}
		for histName, hist := range hset.Items {
			if _, ok := out[histName]; !ok {
				out[histName] = map[string]map[string]int{}
			}
			for binName, binCount := range hist.Items {
				out[hsetName][histName][binName] += binCount
			}
		}
	}
	return out
}

func (hsets *HistogramSets) MapAdd(m map[string]map[string]map[string]int, trimSpace bool) {
	for hsetName, hsetMap := range m {
		for histName, histMap := range hsetMap {
			for binName, binCount := range histMap {
				hsets.Add(hsetName, histName, binName, binCount, trimSpace)
			}
		}
	}
}

func (hsets *HistogramSets) Sum() int {
	sum := 0
	for _, hset := range hsets.Items {
		for _, hist := range hset.Items {
			for _, binSum := range hist.Items {
				sum += binSum
			}
		}
	}
	return sum
}

func (hsets *HistogramSets) UpdateSetOrders(setOrders map[string][]string) {
	for k, vs := range setOrders {
		if hset, ok := hsets.Items[k]; ok {
			hset.Order = slices.Clone(vs)
			hsets.Items[k] = hset
		}
	}
}

func (hsets *HistogramSets) Visit(visit func(hsetName, histName, binName string, binCount int)) {
	for hsetName, hset := range hsets.Items {
		for histName, hist := range hset.Items {
			for binName, binCount := range hist.Items {
				visit(hsetName, histName, binName, binCount)
			}
		}
	}
}
