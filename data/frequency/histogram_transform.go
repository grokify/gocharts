package frequency

import (
	"strings"
)

// TransformBinNames modifies bin names and returns a new
// histogram.
func (hist *Histogram) TransformBinNames(xfFunc func(input string) string) *Histogram {
	if hist == nil {
		return nil
	}
	newHist := NewHistogram(hist.Name)
	for binName, binFreq := range hist.Items {
		newHist.Add(xfFunc(binName), binFreq)
	}
	newHist.Inflate()
	return newHist
}

// TransformBinNamesExactMatch modifies bin names and returns a new
// histogram.
func (hist *Histogram) TransformBinNamesExactMatch(xfMap map[string]string) *Histogram {
	if hist == nil {
		return nil
	}
	return hist.TransformBinNames(
		func(oldName string) string {
			for oldNameTry, newName := range xfMap {
				if oldNameTry == oldName {
					return newName
				}
			}
			return oldName
		},
	)
}

// TransformBinNamesByPrefix modifies bin names and returns a new
// histogram.
func (hist *Histogram) TransformBinNamesByPrefix(xfMap map[string]string) *Histogram {
	if hist == nil {
		return nil
	}
	return hist.TransformBinNames(
		func(oldName string) string {
			for oldPrefix, newName := range xfMap {
				if strings.Index(oldName, oldPrefix) == 0 {
					return newName
				}
			}
			return oldName
		},
	)
}
