package histogram

import (
	"strings"
)

// TransformBinNames modifies bin names and returns a new
// histogram.
func TransformBinNames(hist *Histogram, xfFunc func(input string) string) *Histogram {
	if hist == nil {
		return nil
	}
	newHist := NewHistogram()
	for binName, binFreq := range hist.BinsFrequency {
		newHist.Add(xfFunc(binName), binFreq)
	}
	newHist.Inflate()
	return newHist
}

// TransformBinNamesExactMatch modifies bin names and returns a new
// histogram.
func TransformBinNamesExactMatch(hist *Histogram, xfMap map[string]string) *Histogram {
	if hist == nil {
		return nil
	}
	return TransformBinNames(hist,
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
func TransformBinNamesByPrefix(hist *Histogram, xfMap map[string]string) *Histogram {
	if hist == nil {
		return nil
	}
	return TransformBinNames(hist,
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
