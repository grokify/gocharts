package histogram

import (
	"strings"
)

type HistogramSets struct {
	HistogramSetMap map[string]*HistogramSet
}

func NewHistogramSets() *HistogramSets {
	return &HistogramSets{HistogramSetMap: map[string]*HistogramSet{}}
}

func (fsets *HistogramSets) Add(setKey1, setKey2, binName string, binValue int, trimSpace bool) {
	if trimSpace {
		setKey1 = strings.TrimSpace(setKey1)
		setKey2 = strings.TrimSpace(setKey2)
		binName = strings.TrimSpace(binName)
	}
	fset, ok := fsets.HistogramSetMap[setKey1]
	if !ok {
		fset = NewHistogramSet(setKey1)
	}
	fset.Add(setKey2, binName, binValue)
	fsets.HistogramSetMap[setKey1] = fset
}

func (fsets *HistogramSets) Flatten(name string) *HistogramSet {
	fsetFlat := NewHistogramSet(name)
	for _, fset := range fsets.HistogramSetMap {
		for k2, fstats := range fset.HistogramMap {
			for binName, binValue := range fstats.Items {
				fsetFlat.Add(k2, binName, binValue)
			}
		}
	}
	return fsetFlat
}

func (fsets *HistogramSets) Counts() *HistogramSetsCounts {
	return NewHistogramSetsCounts(*fsets)
}
