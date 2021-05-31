package histogram

import "github.com/grokify/simplego/type/stringsutil"

// HistogramSetsCounts returns UID counts. When used with
// NewHistogramSetsCSV(), it can provide a sanity check
// for raw record counts against aggregate query values,
// e.g. compare counts of raw records to GROUP BY counts.
type HistogramSetsCounts struct {
	UidCounts     map[string]map[string]uint
	UidCountsKey1 map[string]uint
	UidCountsKey2 map[string]uint
	Key1Names     []string
	Key2Names     []string
}

func (hcounts *HistogramSetsCounts) preflate() {
	hcounts.Key1Names = []string{}
	hcounts.Key2Names = []string{}
	hcounts.UidCountsKey1 = map[string]uint{}
	hcounts.UidCountsKey2 = map[string]uint{}
}

func (hcounts *HistogramSetsCounts) Inflate() {
	hcounts.preflate()
	for key1Name, key1Vals := range hcounts.UidCounts {
		hcounts.Key1Names = append(hcounts.Key1Names, key1Name)
		if _, ok := hcounts.UidCountsKey1[key1Name]; !ok {
			hcounts.UidCountsKey1[key1Name] = uint(0)
		}
		for key2Name, k1k2Count := range key1Vals {
			hcounts.Key2Names = append(hcounts.Key2Names, key2Name)
			if _, ok := hcounts.UidCountsKey1[key1Name]; !ok {
				hcounts.UidCountsKey2[key2Name] = uint(0)
			}
			hcounts.UidCountsKey1[key1Name] += k1k2Count
			hcounts.UidCountsKey2[key2Name] += k1k2Count
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
		UidCounts:     map[string]map[string]uint{},
		UidCountsKey1: map[string]uint{},
		UidCountsKey2: map[string]uint{}}
	if len(hsets.HistogramSetMap) == 0 {
		return hcounts
	}

	for hsetName, hset := range hsets.HistogramSetMap {
		hcountsGroup, ok := hcounts.UidCounts[hsetName]
		if !ok {
			hcountsGroup = map[string]uint{}
		}
		for histName, hist := range hset.HistogramMap {
			hcountsGroup[histName] = uint(len(hist.Bins))
		}
		hcounts.UidCounts[hsetName] = hcountsGroup
	}

	hcounts.Inflate()
	return hcounts
}
