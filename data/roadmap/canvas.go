// roadmap provides data for generating roadmaps
package roadmap

import (
	"fmt"
	"time"

	"github.com/grokify/gotilla/math/mathutil"
	tu "github.com/grokify/gotilla/time/timeutil"
)

type Canvas struct {
	MinTime time.Time
	MaxTime time.Time
	MinX    int64
	MaxX    int64
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
		return fmt.Errorf("Max is < min: min [%v] max [%v]", qtrMin, qtrMax)
	}
	err := can.SetMinQuarter(qtrMin)
	if err != nil {
		return err
	}
	return can.SetMaxQuarter(qtrMax)
}

func (can *Canvas) SetMinQuarter(qtr int32) error {
	qt, err := tu.QuarterInt32Start(qtr)
	if err != nil {
		return err
	}
	can.MinTime = qt
	can.MinX = qt.Unix()
	return nil
}

func (can *Canvas) SetMaxQuarter(qtr int32) error {
	qt, err := tu.QuarterInt32End(qtr)
	if err != nil {
		return err
	}
	can.MaxTime = qt
	can.MaxX = qt.Unix()
	return nil
}

func (can *Canvas) SetRange(cells int32) {
	can.Range = mathutil.RangeInt64{
		Min:   can.MinX,
		Max:   can.MaxX,
		Cells: cells,
	}
}

func (can *Canvas) AddItem(i Item) {
	can.Items = append(can.Items, i)
}

func (can *Canvas) InflateItems() error {
	for i, item := range can.Items {
		item, err := can.InflateItem(item)
		if err != nil {
			return err
		}
		can.Items[i] = item
	}
	return nil
}

func (can *Canvas) InflateItem(item Item) (Item, error) {
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
	qt, err := tu.QuarterInt32Start(qtr)
	if err != nil {
		return err
	}
	i.MinTime = qt
	return nil
}

func (i *Item) SetMaxQuarter(qtr int32) error {
	qt, err := tu.QuarterInt32End(qtr)
	if err != nil {
		return err
	}
	i.MaxTime = qt
	return nil
}

func GetCanvasQuarter(start, end int32) (Canvas, error) {
	qs, err := tu.QuarterInt32Start(start)
	if err != nil {
		return Canvas{}, err
	}
	qe, err := tu.QuarterInt32End(end)
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
