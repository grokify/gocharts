package histogram

import (
	"strings"

	"github.com/grokify/mogo/type/maputil"
)

// FilterHistogramNames returns a new `HistogramSet` with only the matching histogram names included.
func (hset *HistogramSet) FilterHistogramNames(inFunc, exFunc func(histname string) bool) *HistogramSet {
	newHset := NewHistogramSet(hset.Name)
	for histName, hist := range hset.HistogramMap {
		if (exFunc != nil && exFunc(histName)) ||
			inFunc == nil ||
			!inFunc(histName) {
			continue
		}
		for binName, binCount := range hist.Bins {
			newHset.Add(histName, binName, binCount)
		}
	}
	return newHset
}

func (hset *HistogramSet) TransformNames(xfFuncHist, xfFuncBin func(input string) string) *HistogramSet {
	if xfFuncHist == nil {
		xfFuncHist = funcStringTransformNoop
	}
	if xfFuncBin == nil {
		xfFuncBin = funcStringTransformNoop
	}
	newHset := NewHistogramSet(hset.Name)
	for histName, hist := range hset.HistogramMap {
		for binName, binCount := range hist.Bins {
			newHset.Add(xfFuncHist(histName), xfFuncBin(binName), binCount)
		}
	}
	return newHset
}

/*
// TransformBinNames modifies histogram names and returns a new histogram set.
func (hset *HistogramSet) TransformBinNames(xfFunc func(input string) string) *HistogramSet {
	newHset := NewHistogramSet(hset.Name)
	for histName, hist := range hset.HistogramMap {
		for binName, binCount := range hist.Bins {
			newHset.Add(histName, xfFunc(binName), binCount)
		}
	}
	return newHset
}

// TransformHistogramNames modifies histogram names and returns a new histogram set.
func (hset *HistogramSet) TransformHistogramNames(xfFunc func(input string) string) *HistogramSet {
	newHset := NewHistogramSet(hset.Name)
	for histName, hist := range hset.HistogramMap {
		for binName, binCount := range hist.Bins {
			newHset.Add(xfFunc(histName), binName, binCount)
		}
	}
	return newHset
}
*/

// TransformBinNamesMap modifies bin names and returns a new `HistogramSet`.
func (hset *HistogramSet) TransformBinNamesMap(xfMap map[string]string, trimSpace bool) *HistogramSet {
	return hset.TransformNames(
		nil,
		func(input string) string {
			if trimSpace {
				input = strings.TrimSpace(input)
			}
			if newBinName, ok := xfMap[input]; ok {
				return newBinName
			}
			return input
		},
	)
}

func (hset *HistogramSet) TransformBinNamesMapSlice(xfMap map[string][]string, dedupe, sortAsc bool, sep string, trimSpace bool) *HistogramSet {
	xfMSS := maputil.MapStringSlice(xfMap)
	return hset.TransformBinNamesMap(xfMSS.FlattenJoin(dedupe, sortAsc, sep), trimSpace)
}

// TransformHistogramNamesMap modifies bin names and returns a new `HistogramSet`. `matchType`
// can be set to `prefix` to match name prefixes instead of exact match.
func (hset *HistogramSet) TransformHistogramNamesMap(xfMap map[string]string, matchType string) *HistogramSet {
	matchType = strings.ToLower(strings.TrimSpace(matchType))
	if matchType == "prefix" {
		return hset.transformHistogramNamesPrefix(xfMap)
	}
	return hset.transformHistogramNamesExactMatch(xfMap)
}

// transformHistogramNamesExactMatch modifies bin names and returns a new histogram.
func (hset *HistogramSet) transformHistogramNamesExactMatch(xfMap map[string]string) *HistogramSet {
	return hset.TransformNames(
		func(oldName string) string {
			for oldNameTry, newName := range xfMap {
				if oldNameTry == oldName {
					return newName
				}
			}
			return oldName
		}, nil,
	)
}

// transformHistogramNamesPrefix modifies bin names and returns a new histogram.
func (hset *HistogramSet) transformHistogramNamesPrefix(xfMap map[string]string) *HistogramSet {
	return hset.TransformNames(
		func(oldName string) string {
			for oldPrefix, newName := range xfMap {
				if strings.Index(oldName, oldPrefix) == 0 {
					return newName
				}
			}
			return oldName
		},
		nil,
	)
}
