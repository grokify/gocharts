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
type Histogram struct {
	Name        string
	Items       map[string]int
	Counts      map[string]int // how many items have counts.
	Percentages map[string]float64
	ItemCount   uint
	Sum         int
}

func NewHistogram(name string) *Histogram {
	return &Histogram{
		Name:        name,
		Items:       map[string]int{},
		Counts:      map[string]int{},
		Percentages: map[string]float64{},
		ItemCount:   0}
}

/*
func (fs *FrequencyStats) AddInt(i int) {
	fs.AddString(strconv.Itoa(i), 1)
}
*/

func (fstats *Histogram) Add(s string, count int) {
	if _, ok := fstats.Items[s]; ok {
		fstats.Items[s] += count
	} else {
		fstats.Items[s] = count
	}
}

func (fs *Histogram) Inflate() {
	fs.Counts = map[string]int{}
	sum := int(0)
	for _, itemCount := range fs.Items {
		countString := strconv.Itoa(itemCount)
		if _, ok := fs.Counts[countString]; !ok {
			fs.Counts[countString] = 0
		}
		fs.Counts[countString]++
		sum += itemCount
	}
	fs.ItemCount = uint(len(fs.Items))

	fs.Percentages = map[string]float64{}
	for itemName, itemCount := range fs.Items {
		fs.Percentages[itemName] = float64(itemCount) / float64(sum)
	}
	fs.Sum = sum
}

func (fs *Histogram) ItemsSlice() []string {
	strs := []string{}
	for key := range fs.Items {
		strs = append(strs, key)
	}
	return strs
}

func (fs *Histogram) ItemsSliceSorted() []string {
	items := fs.ItemsSlice()
	sort.Strings(items)
	return items
}

func (fs *Histogram) TotalCount() uint64 {
	totalCount := 0
	for _, itemCount := range fs.Items {
		totalCount += itemCount
	}
	return uint64(totalCount)
}

func (fs *Histogram) Stats() point.PointSet {
	pointSet := point.NewPointSet()
	for itemName, itemCount := range fs.Items {
		pointSet.PointsMap[itemName] = point.Point{
			Name:        itemName,
			AbsoluteInt: int64(itemCount)}
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
func (fs *Histogram) ItemCounts(sortBy string) []maputil.Record {
	msi := maputil.MapStringInt(fs.Items)
	return msi.Sorted(sortBy)
}

// WriteTable writes an ASCII Table. For CLI apps, pass `os.Stdout` for `io.Writer`.
func (fs *Histogram) WriteTableASCII(writer io.Writer, header []string, sortBy string, inclTotal bool) {
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
