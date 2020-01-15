package histogram

import (
	"sort"
)

type Histogram struct {
	BinCount      int            `json:"binCount"`
	BinsFrequency map[string]int `json:"binsFrequency"`
}

func NewHistogram() Histogram {
	return Histogram{BinsFrequency: map[string]int{}}
}

func (h *Histogram) Inflate() {
	h.BinCount = len(h.BinsFrequency)
}

func (h *Histogram) Add(bin string, count int) {
	if freq, ok := h.BinsFrequency[bin]; !ok {
		h.BinsFrequency[bin] = count
	} else {
		h.BinsFrequency[bin] = freq + count
	}
}

type HistogramSet struct {
	Meta         HistogramSetMetadata `json:"meta,omitempty"`
	HistogramMap map[string]Histogram `json:"histograms"`
}

type HistogramSetMetadata struct {
	Names          []string `json:"names,omitempty"`
	NameCount      int      `json:"nameCount,omitempty"`
	UniqueBinCount int      `json:"uniqueBinCount,omitempty"`
}

func NewHistogramSetMetadata() HistogramSetMetadata {
	return HistogramSetMetadata{Names: []string{}}
}

func NewHistogramSet() HistogramSet {
	return HistogramSet{
		Meta:         NewHistogramSetMetadata(),
		HistogramMap: map[string]Histogram{}}
}

func (hs *HistogramSet) Add(name, bin string, count int) {
	if hs.HistogramMap == nil {
		hs.HistogramMap = map[string]Histogram{}
	}
	if _, ok := hs.HistogramMap[name]; !ok {
		hs.HistogramMap[name] = NewHistogram()
	}
	h := hs.HistogramMap[name]
	h.Add(bin, count)
	hs.HistogramMap[name] = h
}

func (hs *HistogramSet) Inflate() {
	names := []string{}
	uniqueBins := map[string]int{}
	for name, h := range hs.HistogramMap {
		h.Inflate()
		hs.HistogramMap[name] = h
		names = append(names, name)
		for binName, binFreq := range h.BinsFrequency {
			if _, ok := uniqueBins[binName]; !ok {
				uniqueBins[binName] = 0
			}
			uniqueBins[binName] += binFreq
		}
	}
	hs.Meta.UniqueBinCount = len(uniqueBins)
	if hs.Meta.Names == nil {
		hs.Meta.Names = []string{}
	}
	sort.Strings(names)
	hs.Meta.Names = names
	hs.Meta.NameCount = len(names)
}
