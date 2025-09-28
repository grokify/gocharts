package histogram

import (
	"strings"

	"github.com/grokify/mogo/type/maputil"
)

type MatchType string

const (
	MatchTypePrefix MatchType = "prefix"
	MatchTypeExact  MatchType = "exact"
)

// TransformBinNames modifies bin names and returns a new histogram.
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
func (hist *Histogram) TransformBinNamesMap(xfMap map[string]string, matchType MatchType) *Histogram {
	if matchType == MatchTypePrefix {
		return hist.transformBinNamesPrefix(xfMap)
	}
	return hist.transformBinNamesExactMatch(xfMap)
}

func (hist *Histogram) TransformBinNamesMapSlice(xfMap map[string][]string, matchType MatchType, dedupe, sortAsc bool, sep string, def []string) *Histogram {
	xfMSS := maputil.MapStringSlice(xfMap)
	return hist.TransformBinNamesMap(xfMSS.FlattenJoin(dedupe, sortAsc, sep), matchType)
}

// transformBinNamesExactMatch modifies bin names and returns a new histogram.
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

// transformBinNamesPrefix modifies bin names and returns a new histogram.
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
