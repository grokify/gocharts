// roadmap provides data for generating roadmaps
package roadmap

import (
	"fmt"
	"time"

	"github.com/grokify/mogo/math/mathutil"
	tu "github.com/grokify/mogo/time/timeutil"
)

const canvasLogName string = "roadmap.Canvas"

type Canvas struct {
	MinTime time.Time
	MaxTime time.Time
	MinX    int64 // time
	MaxX    int64 // time
	MaxY    int64
	MinY    int64
	MinCell int32
	MaxCell int32
	Range   mathutil.RangeInt64
	Items   []Item
	Rows    [][]Item
}

func (can *Canvas) SetMinMaxQuarter(qtrMin, qtrMax int32) error {
	if qtrMax < qtrMin {
		return fmt.Errorf("max is < min: min [%v] max [%v]", qtrMin, qtrMax)
	}
	err := can.SetMinQuarter(qtrMin)
	if err != nil {
		return err
	}
	return can.SetMaxQuarter(qtrMax)
}

func (can *Canvas) SetMinQuarter(qtr int32) error {
	qt, err := tu.QuarterInt32StartTime(qtr)
	if err != nil {
		return err
	}
	can.MinTime = qt
	can.MinX = qt.Unix()
	return nil
}

func (can *Canvas) SetMaxQuarter(qtr int32) error {
	qt, err := tu.QuarterInt32EndTime(qtr)
	if err != nil {
		return err
	}
	can.MaxTime = qt
	can.MaxX = qt.Unix()
	return nil
}

func (can *Canvas) SetRangeCells(cells int32) {
	can.Range = mathutil.RangeInt64{
		Min:   can.MinX,
		Max:   can.MaxX,
		Cells: cells}
}

func (can *Canvas) AddItem(item Item) {
	// The following if statements check if the item is entirely
	// or partially out of bounds of the canvas. If it is entirely
	// out of bounds, it is not added. If it is partially out of bounds
	// the item times are adjusted to fit.
	if item.MinTime.After(can.MaxTime) ||
		item.MaxTime.Before(can.MinTime) {
		return
	}
	if item.MinTime.Before(can.MinTime) && item.MaxTime.After(can.MinTime) {
		item.MinTime = can.MinTime
	}
	if item.MaxTime.After(can.MaxTime) && item.MinTime.Before(can.MaxTime) {
		item.MaxTime = can.MaxTime
	}
	can.Items = append(can.Items, item)
}

func (can *Canvas) InflateItems() error {
	for i, item := range can.Items {
		// fmt.Printf("ITEM [%v]\n", i)
		// fmtutil.PrintJSON(item)
		item, err := can.InflateItem(item)
		if err != nil {
			return err
		}
		can.Items[i] = item
	}
	return nil
}

func (can *Canvas) InflateItem(item Item) (Item, error) {
	if tu.IsZeroAny(item.MinTime) && tu.IsZeroAny(item.MaxTime) {
		return item, fmt.Errorf("func %s.InflateItem() Error: Need NonZero Time For [%v][%v][%v]", canvasLogName, item.Name, "item.MinTime", "item.MaxTime")
	} else if tu.IsZeroAny(item.MinTime) {
		return item, fmt.Errorf("func %s.InflateItem() Error: Need NonZero Time For [%v][%v]", canvasLogName, item.Name, "item.MinTime")
	} else if tu.IsZeroAny(item.MaxTime) {
		return item, fmt.Errorf("func %s.InflateItem() Error: Need NonZero Time For [%v][%v]", canvasLogName, item.Name, "item.MaxTime")
	}

	cellMargin := int32(1)
	cr, err := can.Range.CellRange()
	if err != nil {
		return item, err
	}
	margin := cr + 1
	cell0, err := can.Range.CellIndexForValue(item.MinTime.Unix())
	if err != nil {
		return item, err
	}
	item.MinCell = cell0 + cellMargin
	min0, _, err := can.Range.CellMinMax(item.MinCell)
	if err != nil {
		return item, err
	}
	item.Min = min0 + margin

	cell1, err := can.Range.CellIndexForValue(item.MaxTime.Unix())
	if err != nil {
		return item, err
	}
	item.MaxCell = cell1 - cellMargin

	_, max1, err := can.Range.CellMinMax(item.MaxCell)
	if err != nil {
		return item, err
	}
	item.Max = max1 - margin

	return item, nil
}

func (can *Canvas) BuildRows() {
	rows := [][]Item{}
ITEMS:
	for _, item := range can.Items {
		if len(rows) == 0 {
			rows = append(rows, []Item{item})
			continue ITEMS
		}
	ROWS:
		for j, row := range rows {
			if len(row) == 0 {
				rows[j] = append(rows[j], item)
				continue ITEMS
			}
			for _, match := range row {
				if mathutil.IsOverlapSortedInt64(item.Min, item.Max, match.Min, match.Max) {
					continue ROWS
				}
			}
			rows[j] = append(rows[j], item)
			continue ITEMS
		}
		rows = append(rows, []Item{item})
	}
	can.Rows = rows
}

/*
	type Item struct {
		MinTime time.Time
		MaxTime time.Time
		MinCell int32 // Inflated by Canvas
		MaxCell int32 // Inflated by Canvas
		Min     int64 // Inflated by Canvas
		Max     int64 // Inflated by Canvas
		Name    string
		URL     string
		Color   string
	}

	func (i *Item) SetMinMaxQuarter(qtrMin, qtrMax int32) error {
		if qtrMax < qtrMin {
			return fmt.Errorf("Max is < min: min [%v] max [%v]", qtrMin, qtrMax)
		}
		err := i.SetMinQuarter(qtrMin)
		if err != nil {
			return err
		}
		return i.SetMaxQuarter(qtrMax)
	}

	func (i *Item) SetMinQuarter(qtr int32) error {
		qt, err := tu.QuarterInt32StartTime(qtr)
		if err != nil {
			return err
		}
		i.MinTime = qt
		return nil
	}

	func (i *Item) SetMaxQuarter(qtr int32) error {
		qt, err := tu.QuarterInt32EndTime(qtr)
		if err != nil {
			return err
		}
		i.MaxTime = qt
		return nil
	}
*/

func GetCanvasQuarter(start, end int32) (Canvas, error) {
	qs, err := tu.QuarterInt32StartTime(start)
	if err != nil {
		return Canvas{}, err
	}
	qe, err := tu.QuarterInt32EndTime(end)
	if err != nil {
		return Canvas{}, err
	}
	return Canvas{
		MinTime: qs,
		MinX:    qs.Unix(),
		MaxTime: qe,
		MaxX:    qe.Unix(),
	}, nil
}
