package histogram

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/olekukonko/tablewriter"

	"github.com/grokify/gocharts/v2/data/point"
	"github.com/grokify/gocharts/v2/data/table"
)

var (
	ErrHistogramCannotBeNil    = errors.New("histogram cannot be nil")
	ErrHistogramSetCannotBeNil = errors.New("histogram set cannot be nil")
)

// Histogram is used to count how many times an item appears and how many times number
// of appearances appear. It can be used with simple string keys or `map[string]string`
// keys which are converted to soerted query strings.
type Histogram struct {
	Name        string
	Bins        map[string]int
	Counts      map[string]int // how many items have counts.
	Percentages map[string]float64
	Order       []string // bin ordering for formatting.
}

func NewHistogram(name string) *Histogram {
	return &Histogram{
		Name:        name,
		Bins:        map[string]int{},
		Counts:      map[string]int{},
		Percentages: map[string]float64{},
		// BinCount:    0}
	}
}

// ReadFileHistogramBins reads a JSON file consisting of a `map[string]int` and
// populates a `Histogram`.
func ReadFileHistogramBins(filename string) (*Histogram, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	msi := map[string]int{}
	if err := json.Unmarshal(b, &msi); err != nil {
		return nil, err
	}
	h := NewHistogram("")
	h.Bins = msi
	return h, nil
}

func (hist *Histogram) Add(binName string, binCount int) {
	hist.Bins[binName] += binCount
}

func (hist *Histogram) AddBulk(m map[string]int) {
	for k, v := range m {
		hist.Add(k, v)
	}
}

func (hist *Histogram) GetOrDefault(binName string, def int) int {
	if v, ok := hist.Bins[binName]; ok {
		return v
	} else {
		return def
	}
}

func (hist *Histogram) Inflate() {
	hist.Counts = map[string]int{}
	sum := 0
	for _, binVal := range hist.Bins {
		countString := strconv.Itoa(binVal)
		if _, ok := hist.Counts[countString]; !ok {
			hist.Counts[countString] = 0
		}
		hist.Counts[countString]++
		sum += binVal
	}

	hist.Percentages = map[string]float64{}
	for binName, binVal := range hist.Bins {
		hist.Percentages[binName] = float64(binVal) / float64(sum)
	}
}

func (hist *Histogram) BinNames() []string {
	return hist.ItemNames()
}

func (hist *Histogram) BinNamesMore(inclOrdered, inclUnordered, inclEmpty bool) []string {
	var names []string
	if inclOrdered {
		seen := map[string]int{}
		for _, ord := range hist.Order {
			seen[ord]++
			if !inclEmpty {
				if v, ok := hist.Bins[ord]; !ok || v == 0 {
					continue
				}
			}
			names = append(names, ord)
		}
		if inclUnordered {
			allNames := hist.BinNames()
			for _, name := range allNames {
				if _, ok := seen[name]; ok {
					continue
				}
				seen[name]++
				if !inclEmpty {
					if v, ok := hist.Bins[name]; !ok || v == 0 {
						continue
					}
				}
				names = append(names, name)
			}
		}
	} else {
		allNames := hist.BinNames()
		for _, name := range allNames {
			if !inclEmpty {
				if v, ok := hist.Bins[name]; !ok || v == 0 {
					continue
				}
			}
			names = append(names, name)
		}
	}

	return names
}

func (hist *Histogram) BinNameExists(binName string) bool {
	if _, ok := hist.Bins[binName]; ok {
		return true
	} else {
		return false
	}
}

func (hist *Histogram) BinValue(binName string) (int, error) {
	if v, ok := hist.Bins[binName]; ok {
		return v, nil
	} else {
		return -1, errors.New("bin not found")
	}
}

func (hist *Histogram) BinValueOrDefault(binName string, def int) int {
	if c, err := hist.BinValue(binName); err != nil {
		return def
	} else {
		return c
	}
}

func (hist *Histogram) BinValuesOrDefault(binNames []string, def int) []int {
	var out []int
	for _, binName := range binNames {
		out = append(out, hist.BinValueOrDefault(binName, def))
	}
	return out
}

func (hist *Histogram) ItemCount() uint {
	return uint(len(hist.Bins))
}

func (hist *Histogram) ItemNames() []string {
	return maputil.Keys(hist.Bins)
}

func (hist *Histogram) Map() map[string]int {
	out := map[string]int{}
	for binName, binCount := range hist.Bins {
		out[binName] += binCount
	}
	return out
}

func (hist *Histogram) MapAdd(m map[string]int) {
	for binName, binCount := range m {
		hist.Add(binName, binCount)
	}
}

// OrderOrDefault returns a list of histogram bin names defaulting to
// ordered names and falling back to sorted bin names. If an order is
// provided, the non-explicitly listed bin names can be included at
// the end or not included.
func (hist *Histogram) OrderOrDefault(inclUnordered bool) []string {
	s1, _ := stringsutil.SliceOrderExplicit(
		maputil.StringKeys(hist.Bins, nil),
		hist.Order,
		inclUnordered)
	return s1
	/*
		return SliceDedupeOrdered(
			maputil.StringKeys(hist.Bins, nil),
			hist.Order,
			true,
			inclUnordered)
	*/
	/*
		var out []string
		if len(hist.Order) > 0 {
			seen := map[string]int{}
			for _, o := range hist.Order {
				if _, ok := seen[o]; !ok {
					out = append(out, o)
					seen[o]++
				}
			}
			if inclUnordered {
				asc := maputil.StringKeys(hist.Bins, nil)
				for _, o := range asc {
					if _, ok := seen[o]; !ok {
						out = append(out, o)
						seen[o]++
					}
				}
			}
		} else {
			out = maputil.StringKeys(hist.Bins, nil)
			sort.Strings(out)
		}
		return out
	*/
}

func (hist *Histogram) Sum() int {
	binSum := 0
	for _, c := range hist.Bins {
		binSum += c
	}
	return binSum
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

// ItemCounts returns sorted bin names and values.
func (hist *Histogram) ItemCounts(sortBy string) maputil.Records {
	msi := maputil.MapStringInt(hist.Bins)
	return msi.Sorted(sortBy)
}

// ItemValuesOrdered returns bin names and values sorted by the `Order` field.
// Unordered bins are not included.
func (hist *Histogram) ItemValuesOrdered() maputil.Records {
	var recs maputil.Records
	for _, ord := range hist.Order {
		if v, ok := hist.Bins[ord]; ok {
			recs = append(recs, maputil.Record{Name: ord, Value: v})
		} else {
			recs = append(recs, maputil.Record{Name: ord, Value: 0})
		}
	}
	return recs
}

// Percentile returns a percentile where all the keys are integers. If it encounters
// a non-integer key, it will return an error.
func (hist *Histogram) Percentile(x int) (float32, error) {
	countTotal := 0
	countLess := 0
	for k, v := range hist.Bins {
		kint, err := strconv.Atoi(k)
		if err != nil {
			return 0, err
		}
		countTotal += v
		if kint < x {
			countLess += v
		}
	}
	return float32(countLess) / float32(countTotal), nil
}

// WriteTable writes an ASCII Table. For CLI apps, pass `os.Stdout` for `io.Writer`.
func (hist *Histogram) WriteTableASCII(w io.Writer, header []string, sortBy string, inclTotal bool) {
	var rows [][]string
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

	table := tablewriter.NewWriter(w)
	table.SetHeader(header)
	if inclTotal {
		table.SetFooter([]string{
			"Total",
			strconv.Itoa(hist.Sum()),
		}) // Add Footer
	}
	table.SetBorder(false) // Set Border to false
	table.AppendBulk(rows) // Add Bulk Data
	table.Render()
}

func (hist *Histogram) Table(colNameBinName, colNameBinCount string) *table.Table {
	tbl := table.NewTable(hist.Name)
	tbl.Columns = []string{colNameBinName, colNameBinCount}
	for binName, binCount := range hist.Bins {
		tbl.Rows = append(tbl.Rows,
			[]string{binName, strconv.Itoa(binCount)})
	}
	tbl.FormatMap = map[int]string{1: "int"}
	return &tbl
}

func (hist *Histogram) WriteXLSX(filename, sheetname, colNameBinName, colNameBinCount string) error {
	tbl := hist.Table(colNameBinName, colNameBinCount)
	return tbl.WriteXLSX(filename, sheetname)
}
