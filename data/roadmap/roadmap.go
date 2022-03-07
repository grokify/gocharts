// roadmap provides data for generating roadmaps
package roadmap

import (
	"time"

	"github.com/grokify/mogo/math/mathutil"
	tu "github.com/grokify/mogo/time/timeutil"
)

type Feature struct {
	Name      string
	URL       string
	StartTime time.Time
	EndTime   time.Time
	StartIdx  int32
	EndIdx    int32
}

func (f *Feature) SortTimes() {
	f.StartTime, f.EndTime = tu.MinMax(f.StartTime, f.EndTime)
}

func (f *Feature) AddIndexes(range64 mathutil.RangeInt64) error {
	f.SortTimes()
	idx0, err := range64.CellIndexForValue(f.StartTime.Unix())
	if err != nil {
		return err
	}
	idx1, err := range64.CellIndexForValue(f.EndTime.Unix())
	if err != nil {
		return err
	}
	f.StartIdx = idx0
	f.EndIdx = idx1
	return nil
}

type Initiative struct {
	Name     string
	Features []Feature
	Rows     [][]Feature // Array of array of non-overlapping features
}

func (init *Initiative) BuildRows(start, end time.Time, range64 mathutil.RangeInt64) ([][]Feature, error) {
	start, end = tu.MinMax(start, end)
	rows := [][]Feature{}
	seen := map[string]int{}
	for _, f := range init.Features {
		if _, ok := seen[f.Name]; ok {
			continue
		}
		f.SortTimes()
		if tu.IsLessThan(f.EndTime, start, false) ||
			tu.IsGreaterThan(f.StartTime, end, false) {
			continue
		}
		err := f.AddIndexes(range64)
		if err != nil {
			return rows, err
		}
		goodRow := -1
		for j, row := range rows {
			isGoodRow := true
			for _, existingFeature := range row {
				if mathutil.IsOverlapSortedInt32(
					f.StartIdx, f.EndIdx,
					existingFeature.StartIdx, existingFeature.EndIdx,
				) {
					isGoodRow = false
					continue
				}
			}
			if isGoodRow {
				goodRow = j
				continue
			}
		}
		if goodRow == -1 {
			rows = append(rows, []Feature{f})
		} else {
			rows[goodRow] = append(rows[goodRow], f)
		}
	}
	return rows, nil
}

type Roadmap struct {
	Initiatives     []Initiative
	ReportStartTime time.Time
	ReportEndTime   time.Time
	Cells           int32
	Range64         mathutil.RangeInt64
}

func NewRoadmap(reportStartTime, reportEndTime time.Time, numCells int32) Roadmap {
	reportStartTime, reportEndTime = tu.MinMax(reportStartTime, reportEndTime)
	return Roadmap{
		Initiatives:     []Initiative{},
		ReportStartTime: reportStartTime,
		ReportEndTime:   reportEndTime,
		Cells:           numCells,
		Range64: mathutil.RangeInt64{
			Min:   reportStartTime.Unix(),
			Max:   reportEndTime.Unix(),
			Cells: numCells,
		},
	}
}

func (rm *Roadmap) Build() error {
	for i, init := range rm.Initiatives {
		_, err := init.BuildRows(rm.ReportStartTime, rm.ReportEndTime, rm.Range64)
		if err != nil {
			return err
		}
		rm.Initiatives[i] = init
	}
	return nil
}

/*

1 2 3 4
Grid

xx x xxxx x x xx x

StartTime
EndTime
*/
