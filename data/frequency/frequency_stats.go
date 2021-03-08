package frequency

import (
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/grokify/gocharts/data/point"
	"github.com/grokify/simplego/type/maputil"
	"github.com/olekukonko/tablewriter"
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

func (fstats *FrequencyStats) AddStringMore(s string, count int) {
	if _, ok := fstats.Items[s]; !ok {
		fstats.Items[s] = 0
	}
	fstats.Items[s] += count
}

func (fstats *FrequencyStats) AddString(s string) {
	fstats.AddStringMore(s, 1)
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

const (
	SortNameAsc   = maputil.SortNameAsc
	SortNameDesc  = maputil.SortNameDesc
	SortValueAsc  = maputil.SortValueAsc
	SortValueDesc = maputil.SortValueDesc
)

// ItemCounts returns sorted item names and values.
func (fs *FrequencyStats) ItemCounts(sortBy string) []maputil.Record {
	msi := maputil.MapStringInt(fs.Items)
	return msi.Sorted(sortBy)
}

// WriteTable writes an ASCII Table. For CLI apps, pass `os.Stdout` for `io.Writer`.
func (fs *FrequencyStats) WriteTableASCII(writer io.Writer, header []string, sortBy string, inclTotal bool) {
	rows := [][]string{}
	sortedItems := fs.ItemCounts(sortBy)
	for _, sortedItem := range sortedItems {
		rows = append(rows, []string{
			sortedItem.Name, strconv.Itoa(sortedItem.Value)})
	}

	if len(header) == 0 {
		header = []string{"Name", "Value"}
	} else if len(header) == 1 {
		header[1] = "Value"
	}
	header[0] = strings.TrimSpace(header[0])
	header[1] = strings.TrimSpace(header[1])
	if len(header[0]) == 0 {
		header[0] = "Name"
	}
	if len(header[1]) == 0 {
		header[1] = "Value"
	}

	table := tablewriter.NewWriter(writer)
	table.SetHeader(header)
	if inclTotal {
		table.SetFooter([]string{
			"Total",
			strconv.Itoa(int(fs.TotalCount())),
		}) // Add Footer
	}
	table.SetBorder(false) // Set Border to false
	table.AppendBulk(rows) // Add Bulk Data
	table.Render()
}
