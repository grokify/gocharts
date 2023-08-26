package histogram

import (
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/olekukonko/tablewriter"

	"github.com/grokify/gocharts/v2/data/point"
	"github.com/grokify/gocharts/v2/data/table"
)

// Histogram is used to count how many times an item appears and how many times number
// of appearances appear. It can be used with simple string keys or `map[string]string`
// keys which are converted to soerted query strings.
type Histogram struct {
	Name        string
	Bins        map[string]int
	Counts      map[string]int // how many items have counts.
	Percentages map[string]float64
	// BinCount    uint
	// Sum         int
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

func (hist *Histogram) Add(binName string, binCount int) {
	hist.Bins[binName] += binCount
}

func (hist *Histogram) AddBulk(m map[string]int) {
	for k, v := range m {
		hist.Add(k, v)
	}
}

// AddMap provides a helper function to automatically create url encoded string keys.
// This can be used with `TableMap` to generate tables with arbitrary columns easily.
func (hist *Histogram) AddMap(binMap map[string]string, binCount int) {
	m := maputil.MapStringString(binMap)
	key := m.Encode()
	hist.Add(key, binCount)
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
	// hist.BinCount = uint(len(hist.Bins))

	hist.Percentages = map[string]float64{}
	for binName, binVal := range hist.Bins {
		hist.Percentages[binName] = float64(binVal) / float64(sum)
	}
	// hist.Sum = sum
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

func (hist *Histogram) Sum() int {
	binSum := 0
	for _, c := range hist.Bins {
		binSum += c
	}
	return binSum
}

func (hist *Histogram) KeyCount() int {
	return len(hist.Bins)
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
func (hist *Histogram) WriteTableASCII(w io.Writer, header []string, sortBy string, inclTotal bool) {
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

// MapKeys returns a list of keys using query string keys.
func (hist *Histogram) MapKeys() ([]string, error) {
	keys := map[string]int{}
	for qry := range hist.Bins {
		m, err := maputil.ParseMapStringString(qry)
		if err != nil {
			return []string{}, err
		}
		for k := range m {
			keys[k]++
		}
	}
	return maputil.Keys(keys), nil
}

// MayKeyValues returns a list of keys using query string keys.
func (hist *Histogram) MayKeyValues(key string, dedupe bool) ([]string, error) {
	vals := []string{}
	for qry := range hist.Bins {
		m, err := maputil.ParseMapStringString(qry)
		if err != nil {
			return []string{}, err
		}
		if v, ok := m[key]; ok {
			vals = append(vals, v)
		}
	}
	if dedupe {
		vals = slicesutil.Dedupe(vals)
	}
	return vals, nil
}

/*
func mapStringStringSubset(m map[string]string, keys []string, inclUnknown, trimSpace, inclEmpty bool) map[string]string {
	newMap := map[string]string{}
	keyMap := map[string]int{}
	for i, k := range keys {
		keyMap[k] = i
	}
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if trimSpace {
				v = strings.TrimSpace(v)
			}
			if !inclEmpty && v == "" {
				continue
			}
			newMap[k] = v
		} else if inclEmpty {
			newMap[k] = ""
		}
	}
	return newMap
}
*/

// TableMap is used to generate a table using map keys.
func (hist *Histogram) TableMap(mapCols []string, colNameBinCount string) (*table.Table, error) {
	if strings.TrimSpace(colNameBinCount) == "" {
		colNameBinCount = "Count"
	}

	// create histogram with minimized aggregate map keys to aggregate exclude non-desired
	// properties from the key for aggregation.
	histSubset := NewHistogram("")
	for binName, binCount := range hist.Bins {
		binMap, err := maputil.ParseMapStringString(binName)
		if err != nil {
			return nil, err
		}
		newBinMap := binMap.Subset(mapCols, false, true, true)
		// newBinMap := mapStringStringSubset(binMap, mapCols, true, false, true)
		// fmtutil.PrintJSON(newBinMap)
		histSubset.AddMap(newBinMap, binCount)
	}

	tbl := table.NewTable(hist.Name)
	tbl.Columns = append(mapCols, colNameBinCount)

	for binName, binCount := range histSubset.Bins {
		binMap, err := maputil.ParseMapStringString(binName)
		if err != nil {
			return nil, err
		}
		binVals := binMap.Gets(true, mapCols)

		tbl.Rows = append(tbl.Rows,
			append(binVals, strconv.Itoa(binCount)),
		)
	}

	tbl.FormatMap = map[int]string{len(tbl.Columns) - 1: "int"}
	return &tbl, nil
}

/*
// TableMap is used to generate a table using map keys.
func (hist *Histogram) TableMapOld(mapCols []string, colNameBinCount string) (*table.Table, error) {
	tbl := table.NewTable(hist.Name)
	if strings.TrimSpace(colNameBinCount) == "" {
		colNameBinCount = "Count"
	}
	tbl.Columns = mapCols
	tbl.Columns = append(mapCols, colNameBinCount)
	for binName, binCount := range hist.Bins {
		binMap, err := maputil.ParseMapStringString(binName)
		if err != nil {
			return nil, err
		}
		binVals := binMap.Gets(true, mapCols)

		tbl.Rows = append(tbl.Rows,
			append(binVals, strconv.Itoa(binCount)),
		)
	}
	tbl.FormatMap = map[int]string{len(tbl.Columns) - 1: "int"}
	return &tbl, nil
}
*/

func (hist *Histogram) WriteXLSX(filename, sheetname, colNameBinName, colNameBinCount string) error {
	tbl := hist.Table(colNameBinName, colNameBinCount)
	return tbl.WriteXLSX(filename, sheetname)
}
