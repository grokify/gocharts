package histogram

import "github.com/grokify/mogo/type/stringsutil"

// HistogramSetsCounts returns UID counts. When used with
// NewHistogramSetsCSV(), it can provide a sanity check
// for raw record counts against aggregate query values,
// e.g. compare counts of raw records to GROUP BY counts.
type HistogramSetsCounts struct {
	UIDCounts     map[string]map[string]uint
	UIDCountsKey1 map[string]uint
	UIDCountsKey2 map[string]uint
	Key1Names     []string
	Key2Names     []string
}

func (hcounts *HistogramSetsCounts) preflate() {
	hcounts.Key1Names = []string{}
	hcounts.Key2Names = []string{}
	hcounts.UIDCountsKey1 = map[string]uint{}
	hcounts.UIDCountsKey2 = map[string]uint{}
}

func (hcounts *HistogramSetsCounts) Inflate() {
	hcounts.preflate()
	for key1Name, key1Vals := range hcounts.UIDCounts {
		hcounts.Key1Names = append(hcounts.Key1Names, key1Name)
		if _, ok := hcounts.UIDCountsKey1[key1Name]; !ok {
			hcounts.UIDCountsKey1[key1Name] = uint(0)
		}
		for key2Name, k1k2Count := range key1Vals {
			hcounts.Key2Names = append(hcounts.Key2Names, key2Name)
			if _, ok := hcounts.UIDCountsKey1[key1Name]; !ok {
				hcounts.UIDCountsKey2[key2Name] = uint(0)
			}
			hcounts.UIDCountsKey1[key1Name] += k1k2Count
			hcounts.UIDCountsKey2[key2Name] += k1k2Count
		}
	}
	hcounts.Key1Names = stringsutil.SliceCondenseSpace(
		hcounts.Key1Names, true, true)
	hcounts.Key2Names = stringsutil.SliceCondenseSpace(
		hcounts.Key2Names, true, true)
}

func NewHistogramSetsCounts(hsets HistogramSets) *HistogramSetsCounts {
	hcounts := &HistogramSetsCounts{
		Key1Names:     []string{},
		Key2Names:     []string{},
		UIDCounts:     map[string]map[string]uint{},
		UIDCountsKey1: map[string]uint{},
		UIDCountsKey2: map[string]uint{}}
	if len(hsets.Items) == 0 {
		return hcounts
	}

	for hsetName, hset := range hsets.Items {
		hcountsGroup, ok := hcounts.UIDCounts[hsetName]
		if !ok {
			hcountsGroup = map[string]uint{}
		}
		for histName, hist := range hset.Items {
			hcountsGroup[histName] = uint(len(hist.Items))
		}
		hcounts.UIDCounts[hsetName] = hcountsGroup
	}

	hcounts.Inflate()
	return hcounts
}
