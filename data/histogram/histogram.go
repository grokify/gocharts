package histogram

import (
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/grokify/gocharts/data/point"
	"github.com/grokify/simplego/type/maputil"
	"github.com/olekukonko/tablewriter"
)

// Histogram stats is used to count how many times
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

func (hist *Histogram) Add(s string, count int) {
	if _, ok := hist.Items[s]; ok {
		hist.Items[s] += count
	} else {
		hist.Items[s] = count
	}
}

func (hist *Histogram) Inflate() {
	hist.Counts = map[string]int{}
	sum := int(0)
	for _, itemCount := range hist.Items {
		countString := strconv.Itoa(itemCount)
		if _, ok := hist.Counts[countString]; !ok {
			hist.Counts[countString] = 0
		}
		hist.Counts[countString]++
		sum += itemCount
	}
	hist.ItemCount = uint(len(hist.Items))

	hist.Percentages = map[string]float64{}
	for itemName, itemCount := range hist.Items {
		hist.Percentages[itemName] = float64(itemCount) / float64(sum)
	}
	hist.Sum = sum
}

func (hist *Histogram) ItemsSlice() []string {
	strs := []string{}
	for key := range hist.Items {
		strs = append(strs, key)
	}
	return strs
}

func (hist *Histogram) ItemsSliceSorted() []string {
	items := hist.ItemsSlice()
	sort.Strings(items)
	return items
}

func (hist *Histogram) TotalCount() uint64 {
	totalCount := 0
	for _, itemCount := range hist.Items {
		totalCount += itemCount
	}
	return uint64(totalCount)
}

func (hist *Histogram) Stats() point.PointSet {
	pointSet := point.NewPointSet()
	for itemName, itemCount := range hist.Items {
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
func (hist *Histogram) ItemCounts(sortBy string) []maputil.Record {
	msi := maputil.MapStringInt(hist.Items)
	return msi.Sorted(sortBy)
}

// WriteTable writes an ASCII Table. For CLI apps, pass `os.Stdout` for `io.Writer`.
func (hist *Histogram) WriteTableASCII(writer io.Writer, header []string, sortBy string, inclTotal bool) {
	rows := [][]string{}
	sortedItems := hist.ItemCounts(sortBy)
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
			strconv.Itoa(int(hist.TotalCount())),
		}) // Add Footer
	}
	table.SetBorder(false) // Set Border to false
	table.AppendBulk(rows) // Add Bulk Data
	table.Render()
}
