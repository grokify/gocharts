package histogram

import (
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/grokify/gocharts/data/point"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/simplego/type/maputil"
	"github.com/olekukonko/tablewriter"
)

// Histogram stats is used to count how many times
// an item appears and how many times number of
// appearances appear.
type Histogram struct {
	Name        string
	Bins        map[string]int
	Counts      map[string]int // how many items have counts.
	Percentages map[string]float64
	BinCount    uint
	Sum         int
}

func NewHistogram(name string) *Histogram {
	return &Histogram{
		Name:        name,
		Bins:        map[string]int{},
		Counts:      map[string]int{},
		Percentages: map[string]float64{},
		BinCount:    0}
}

/*
func (hist *Histogram) AddInt(i int) {
	hist.Add(strconv.Itoa(i), 1)
}
*/

func (hist *Histogram) Add(binName string, binCount int) {
	hist.Bins[binName] += binCount
}

func (hist *Histogram) Inflate() {
	hist.Counts = map[string]int{}
	sum := 0
	for _, binCount := range hist.Bins {
		countString := strconv.Itoa(binCount)
		if _, ok := hist.Counts[countString]; !ok {
			hist.Counts[countString] = 0
		}
		hist.Counts[countString]++
		sum += binCount
	}
	hist.BinCount = uint(len(hist.Bins))

	hist.Percentages = map[string]float64{}
	for binName, binCount := range hist.Bins {
		hist.Percentages[binName] = float64(binCount) / float64(sum)
	}
	hist.Sum = sum
}

func (hist *Histogram) BinNames() []string {
	binNames := []string{}
	for binName := range hist.Bins {
		binNames = append(binNames, binName)
	}
	sort.Strings(binNames)
	return binNames
}

func (hist *Histogram) BinNameExists(binName string) bool {
	if _, ok := hist.Bins[binName]; ok {
		return true
	}
	return false
}

func (hist *Histogram) ValueSum() int {
	totalCount := 0
	for _, binCount := range hist.Bins {
		totalCount += binCount
	}
	return totalCount
}

func (hist *Histogram) Stats() point.PointSet {
	pointSet := point.NewPointSet()
	for binName, binCount := range hist.Bins {
		pointSet.PointsMap[binName] = point.Point{
			Name:        binName,
			AbsoluteInt: int64(binCount)}
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
	msi := maputil.MapStringInt(hist.Bins)
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
			strconv.Itoa(hist.ValueSum()),
		}) // Add Footer
	}
	table.SetBorder(false) // Set Border to false
	table.AppendBulk(rows) // Add Bulk Data
	table.Render()
}

func (hist *Histogram) ToTable(colNameBinName, colNameBinCount string) *table.Table {
	tbl := table.NewTable()
	tbl.Name = hist.Name
	tbl.Columns = []string{colNameBinName, colNameBinCount}
	for binName, binCount := range hist.Bins {
		tbl.Rows = append(tbl.Rows,
			[]string{binName, strconv.Itoa(binCount)})
	}
	tbl.FormatMap = map[int]string{1: "int"}
	return &tbl
}

func (hist *Histogram) WriteXLSX(filename, sheetname, colNameBinName, colNameBinCount string) error {
	tbl := hist.ToTable(colNameBinName, colNameBinCount)
	return tbl.WriteXLSX(filename, sheetname)
}
