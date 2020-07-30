package frequency

import (
	"strings"
)

type FrequencySets struct {
	FrequencySetMap map[string]FrequencySet
}

func NewFrequencySets() FrequencySets {
	return FrequencySets{FrequencySetMap: map[string]FrequencySet{}}
}

func (fsets *FrequencySets) Add(key1, key2, uid string, trimSpace bool) {
	if trimSpace {
		key1 = strings.TrimSpace(key1)
		key2 = strings.TrimSpace(key2)
		uid = strings.TrimSpace(uid)
	}
	fset, ok := fsets.FrequencySetMap[key1]
	if !ok {
		fset = NewFrequencySet(key1)
	}
	fset.AddString(key2, uid)
	fsets.FrequencySetMap[key1] = fset
}

func (fsets *FrequencySets) Flatten(name string) FrequencySet {
	fsetFlat := NewFrequencySet(name)
	for _, fset := range fsets.FrequencySetMap {
		for k2, fstats := range fset.FrequencyMap {
			for item, count := range fstats.Items {
				fsetFlat.AddStringMore(k2, item, count)
			}
		}
	}
	return fsetFlat
}

func (fsets *FrequencySets) Counts() FrequencySetsCounts {
	return NewFrequencySetsCounts(*fsets)
}
