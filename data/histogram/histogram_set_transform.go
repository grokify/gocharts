package histogram

import (
	"strings"
)

// TransformHistogramNames modifies histogram names and returns a new
// histogram set.
func (hset *HistogramSet) TransformHistogramNames(xfFunc func(input string) string) *HistogramSet {
	if hset == nil {
		return nil
	}
	newHset := NewHistogramSet(hset.Name)
	for histName, hist := range hset.HistogramMap {
		for binName, binCount := range hist.Bins {
			newHset.Add(xfFunc(histName), binName, binCount)
		}
	}
	return newHset
}

// TransformHistogramNamesMap modifies bin names and returns a new
// `HistogramSet`. `matchType` can be set to `prefix` to match name
// prefixes instead of exact match.
func (hset *HistogramSet) TransformHistogramNamesMap(xfMap map[string]string, matchType string) *HistogramSet {
	matchType = strings.ToLower(strings.TrimSpace(matchType))
	if matchType == "prefix" {
		return hset.transformHistogramNamesPrefix(xfMap)
	}
	return hset.transformHistogramNamesExactMatch(xfMap)
}

// transformHistogramNamesExactMatch modifies bin names and returns a new
// histogram.
func (hset *HistogramSet) transformHistogramNamesExactMatch(xfMap map[string]string) *HistogramSet {
	if hset == nil {
		return nil
	}
	return hset.TransformHistogramNames(
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

// transformHistogramNamesPrefix modifies bin names and returns a new
// histogram.
func (hset *HistogramSet) transformHistogramNamesPrefix(xfMap map[string]string) *HistogramSet {
	if hset == nil {
		return nil
	}
	return hset.TransformHistogramNames(
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
