// roadmap2 represents a simplified set of data structures to represent a roadmap.
package roadmap2

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/maputil"
)

// Item represents a "roadmap item" or box on a roadmap slide.
type Item struct {
	Name        string
	Description string
	ReleaseTime time.Time
	StreamName  string
	Meta        map[string]string
}

func NewItem() Item {
	return Item{Meta: map[string]string{}}
}

func (i Item) ReleaseTimeQuarter() int32 {
	tm := timeutil.NewTimeMore(i.ReleaseTime, 0)
	q := tm.Quarter()
	if q == 11 {
		q = -1
	}
	return q
}

type Items []Item

func (ii Items) FilterByMeta(metaFilterAnd map[string][]string) Items {
	var out Items
	for _, try := range ii {
		if maputil.IsSubsetOrValues(try.Meta, metaFilterAnd) {
			out = append(out, try)
		}
	}
	return out
}

func (ii Items) CountsByMetaAttribute(attrName string, keyForMissing *string) map[string]int {
	counts := map[string]int{}
	for _, item := range ii {
		if item.Meta == nil {
			if keyForMissing != nil {
				counts[pointer.Dereference(keyForMissing)]++
			}
		} else {
			if v, ok := item.Meta[attrName]; ok {
				counts[v]++
			} else if keyForMissing != nil {
				counts[pointer.Dereference(keyForMissing)]++
			}
		}
	}
	return counts
}

func (ii Items) NamesByIntervals(intervals []string, sortAsc, inclNonExplicit bool, defaultInterval timeutil.Interval) (map[string][]string, error) {
	out := map[string][]string{}
	for _, item := range ii {
		added := false
		for _, invString := range intervals {
			invString = strings.ToUpper(strings.TrimSpace(invString))
			invTimeRange, err := timeutil.ParseTimeRangeInterval(invString)
			if err != nil {
				return out, err
			}
			if contained, err := invTimeRange.Contains(item.ReleaseTime, true, true); err != nil {
				return out, err
			} else if contained {
				if _, ok := out[invString]; !ok {
					out[invString] = []string{}
				}
				out[invString] = append(out[invString], item.Name)
				added = true
				break
			}
		}
		if !added && inclNonExplicit {
			tm := timeutil.NewTimeMore(item.ReleaseTime, 0)
			qtrString := tm.YearQuarter()
			if _, ok := out[qtrString]; !ok {
				out[qtrString] = []string{}
			}
			out[qtrString] = append(out[qtrString], item.Name)
		}
	}
	if sortAsc {
		for k := range out {
			sort.Strings(out[k])
		}
	}
	return out, nil
}

func (ii Items) NamesByQuarter(sortAsc bool) map[int32][]string {
	out := map[int32][]string{}
	for _, item := range ii {
		q := item.ReleaseTimeQuarter()
		if out[q] == nil {
			out[q] = []string{}
		}
		out[q] = append(out[q], item.Name)
	}
	if sortAsc {
		for k := range out {
			sort.Strings(out[k])
		}
	}
	return out
}

func MapInt32SliceStringToLines(m map[int32][]string, keyPrefix, valPrefix string) []string {
	var lines []string
	for k, vals := range m {
		kStr := strconv.Itoa(int(k))
		kStr = keyPrefix + kStr
		lines = append(lines, kStr)
		for _, v := range vals {
			v = valPrefix + v
			lines = append(lines, v)
		}
	}
	return lines
}

func ItemsFromTable(t *table.Table, nameColIdx, timeColIdx uint, timeLayout string, metaColIdxs []uint) (Items, error) {
	var items Items
	if t == nil {
		return items, errors.New("table cannot be nil")
	}
	for _, metaColIdx := range metaColIdxs {
		if metaColIdx >= uint(len(t.Columns)) {
			return items, errors.New("cols too short")
		}
	}

	for _, row := range t.Rows {
		if nameColIdx >= uint(len(row)) {
			return items, errors.New("row too short")
		} else if timeColIdx >= uint(len(row)) {
			return items, errors.New("row too short")
		}
		item := NewItem()
		item.Name = strings.TrimSpace(row[nameColIdx])
		tStr := row[timeColIdx]
		if strings.TrimSpace(tStr) != "" {
			if t, err := time.Parse(timeLayout, tStr); err != nil {
				return items, err
			} else {
				item.ReleaseTime = t
			}
		}
		for _, metaColIdx := range metaColIdxs {
			if metaColIdx >= uint(len(row)) {
				return items, errors.New("row too short")
			} else {
				item.Meta[t.Columns[metaColIdx]] = row[metaColIdx]
			}
		}
		items = append(items, item)
	}

	return items, nil
}
