package frequency

import "strconv"

// Frequency stats is used to count how many times
// an item appears and how many times number of
// appearances appear.
type FrequencyStats struct {
	Name   string
	Items  map[string]int
	Counts map[string]int // how many items have counts.
}

func NewFrequencyStats(name string) FrequencyStats {
	return FrequencyStats{
		Name:   name,
		Items:  map[string]int{},
		Counts: map[string]int{}}
}

func (fs *FrequencyStats) AddItem(s string) {
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
}
