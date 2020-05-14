package frequency

import "github.com/grokify/gotilla/type/stringsutil"

type FrequencySetsCounts struct {
	UidCounts     map[string]map[string]uint
	UidCountsKey1 map[string]uint
	UidCountsKey2 map[string]uint
	Key1Names     []string
	Key2Names     []string
}

func (fcounts *FrequencySetsCounts) Inflate() {
	for key1Name, key1Vals := range fcounts.UidCounts {
		fcounts.Key1Names = append(fcounts.Key1Names, key1Name)
		if _, ok := fcounts.UidCountsKey1[key1Name]; !ok {
			fcounts.UidCountsKey1[key1Name] = uint(0)
		}
		for key2Name, k1k2Count := range key1Vals {
			fcounts.Key2Names = append(fcounts.Key2Names, key2Name)
			if _, ok := fcounts.UidCountsKey1[key1Name]; !ok {
				fcounts.UidCountsKey2[key2Name] = uint(0)
			}
			fcounts.UidCountsKey1[key1Name] += k1k2Count
			fcounts.UidCountsKey2[key2Name] += k1k2Count
		}
	}
	fcounts.Key1Names = stringsutil.SliceCondenseSpace(
		fcounts.Key1Names, true, true)
	fcounts.Key2Names = stringsutil.SliceCondenseSpace(
		fcounts.Key2Names, true, true)
}

func NewFrequencySetsCounts(fsets FrequencySets) FrequencySetsCounts {
	fcounts := FrequencySetsCounts{
		Key1Names:     []string{},
		Key2Names:     []string{},
		UidCounts:     map[string]map[string]uint{},
		UidCountsKey1: map[string]uint{},
		UidCountsKey2: map[string]uint{}}
	if len(fsets.FrequencySetMap) == 0 {
		return fcounts
	}

	for fsetGroupName, fsetGroup := range fsets.FrequencySetMap {
		//fcounts.Key1Names = append(fcounts.Key1Names, fsetGroupName)
		fcountsGroup, ok := fcounts.UidCounts[fsetGroupName]
		if !ok {
			fcountsGroup = map[string]uint{}
		}
		for fstatsName, fstats := range fsetGroup.FrequencyMap {
			fcountsGroup[fstatsName] = uint(len(fstats.Items))
		}
		fcounts.UidCounts[fsetGroupName] = fcountsGroup
	}

	fcounts.Inflate()
	return fcounts
}
