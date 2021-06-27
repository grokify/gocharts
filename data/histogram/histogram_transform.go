package histogram

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
	for binName, binCount := range hist.Bins {
		newHist.Add(xfFunc(binName), binCount)
	}
	newHist.Inflate()
	return newHist
}

// TransformBinNamesMap modifies bin names and returns a new
// histogram. `matchType` can be set to `prefix` to match name
// prefixes instead of exact match.
func (hist *Histogram) TransformBinNamesMap(xfMap map[string]string, matchType string) *Histogram {
	matchType = strings.ToLower(strings.TrimSpace(matchType))
	if matchType == "prefix" {
		return hist.transformBinNamesPrefix(xfMap)
	}
	return hist.transformBinNamesExactMatch(xfMap)
}

// transformBinNamesExactMatch modifies bin names and returns a new
// histogram.
func (hist *Histogram) transformBinNamesExactMatch(xfMap map[string]string) *Histogram {
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

// transformBinNamesPrefix modifies bin names and returns a new
// histogram.
func (hist *Histogram) transformBinNamesPrefix(xfMap map[string]string) *Histogram {
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
