package histogram

import (
	"strings"
)

type HistogramSets struct {
	Name            string
	HistogramSetMap map[string]*HistogramSet
}

func NewHistogramSets() *HistogramSets {
	return &HistogramSets{HistogramSetMap: map[string]*HistogramSet{}}
}

func (hsets *HistogramSets) Add(setKey1, setKey2, binName string, binValue int, trimSpace bool) {
	if trimSpace {
		setKey1 = strings.TrimSpace(setKey1)
		setKey2 = strings.TrimSpace(setKey2)
		binName = strings.TrimSpace(binName)
	}
	fset, ok := hsets.HistogramSetMap[setKey1]
	if !ok {
		fset = NewHistogramSet(setKey1)
	}
	fset.Add(setKey2, binName, binValue)
	hsets.HistogramSetMap[setKey1] = fset
}

func (hsets *HistogramSets) Flatten(name string) *HistogramSet {
	hsetFlat := NewHistogramSet(name)
	for _, hset := range hsets.HistogramSetMap {
		for histName, hist := range hset.HistogramMap {
			for binName, binCount := range hist.Bins {
				hsetFlat.Add(histName, binName, binCount)
			}
		}
	}
	return hsetFlat
}

func (hsets *HistogramSets) Counts() *HistogramSetsCounts {
	return NewHistogramSetsCounts(*hsets)
}
