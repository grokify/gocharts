package frequency

import (
	"sort"
	"strconv"

	"github.com/grokify/gocharts/data/point"
)

// Frequency stats is used to count how many times
// an item appears and how many times number of
// appearances appear.
type FrequencyStats struct {
	Name      string
	Items     map[string]int
	Counts    map[string]int // how many items have counts.
	ItemCount uint
}

func NewFrequencyStats(name string) FrequencyStats {
	return FrequencyStats{
		Name:      name,
		Items:     map[string]int{},
		Counts:    map[string]int{},
		ItemCount: 0}
}

func (fs *FrequencyStats) AddInt(i int) {
	fs.AddString(strconv.Itoa(i))
}

func (fs *FrequencyStats) AddStringMore(s string, count int) {
	for i := 0; i < count; i++ {
		fs.AddString(s)
	}
}

func (fs *FrequencyStats) AddString(s string) {
	if _, ok := fs.Items[s]; !ok {
		fs.Items[s] = 0
	}
	fs.Items[s]++
}

func (fs *FrequencyStats) Inflate() {
	fs.Counts = map[string]int{}
	for _, itemCount := range fs.Items {
		countString := strconv.Itoa(itemCount)
		if _, ok := fs.Counts[countString]; !ok {
			fs.Counts[countString] = 0
		}
		fs.Counts[countString]++
	}
	fs.ItemCount = uint(len(fs.Items))
}

func (fs *FrequencyStats) ItemsSlice() []string {
	strs := []string{}
	for key := range fs.Items {
		strs = append(strs, key)
	}
	return strs
}

func (fs *FrequencyStats) ItemsSliceSorted() []string {
	items := fs.ItemsSlice()
	sort.Strings(items)
	return items
}

func (fs *FrequencyStats) TotalCount() uint64 {
	totalCount := 0
	for _, itemCount := range fs.Items {
		totalCount += itemCount
	}
	return uint64(totalCount)
}

func (fs *FrequencyStats) Stats() point.PointSet {
	pointSet := point.NewPointSet()
	for itemName, itemCount := range fs.Items {
		point := point.Point{
			Name:        itemName,
			AbsoluteInt: int64(itemCount)}
		// Percentage:  float64(itemCount) / float64(totalCount) * 100}
		pointSet.PointsMap[itemName] = point
	}
	pointSet.Inflate()
	return pointSet
}
