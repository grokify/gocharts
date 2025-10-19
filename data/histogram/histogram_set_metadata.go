package histogram

import "sort"

type HistogramSetMetadata struct {
	Names          []string `json:"names,omitempty"`
	NameCount      int      `json:"nameCount,omitempty"`
	UniqueBinCount int      `json:"uniqueBinCount,omitempty"`
}

func NewHistogramSetMetadata(histSet *HistogramSet) *HistogramSetMetadata {
	if histSet == nil {
		return &HistogramSetMetadata{Names: []string{}}
	}
	return buildHistogramSetMetadata(histSet)
}

func buildHistogramSetMetadata(hs *HistogramSet) *HistogramSetMetadata {
	meta := &HistogramSetMetadata{Names: []string{}}
	names := []string{}
	uniqueBins := map[string]int{}
	for name, h := range hs.Items {
		h.Inflate()
		hs.Items[name] = h
		names = append(names, name)
		for binName, binCount := range h.Items {
			if _, ok := uniqueBins[binName]; !ok {
				uniqueBins[binName] = 0
			}
			uniqueBins[binName] += binCount
		}
	}
	meta.UniqueBinCount = len(uniqueBins)
	sort.Strings(names)
	meta.Names = names
	meta.NameCount = len(names)
	return meta
}
