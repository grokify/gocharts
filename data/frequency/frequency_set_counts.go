package frequency

type FrequencySetsCounts struct {
	UidCounts  map[string]map[string]uint
	GroupNames []string
}

func NewFrequencySetsCounts(fsets FrequencySets) FrequencySetsCounts {
	fcounts := FrequencySetsCounts{
		GroupNames: []string{},
		UidCounts:  map[string]map[string]uint{}}
	if len(fsets.FrequencySetMap) == 0 {
		return fcounts
	}

	for fsetGroupName, fsetGroup := range fsets.FrequencySetMap {
		fcounts.GroupNames = append(fcounts.GroupNames, fsetGroupName)
		fcountsGroup, ok := fcounts.UidCounts[fsetGroupName]
		if !ok {
			fcountsGroup = map[string]uint{}
		}
		for fstatsName, fstats := range fsetGroup.FrequencyMap {
			fcountsGroup[fstatsName] = uint(len(fstats.Items))
		}
		fcounts.UidCounts[fsetGroupName] = fcountsGroup
	}

	return fcounts
}
