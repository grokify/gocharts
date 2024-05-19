package piechart

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/grokify/mogo/strconv/strconvutil"
)

type PieChartData struct {
	Data    []PieChartDatum
	IsFloat bool
}

func (d *PieChartData) AddInts(v map[string]int) {
	for k, v := range v {
		d.AddInt64(k, int64(v))
	}
}

func (d *PieChartData) AddInt64(s string, v int64) {
	if d.IsFloat {
		d.Data = append(d.Data, PieChartDatum{
			Name:     s,
			ValFloat: float64(v),
			IsFloat:  true,
		})
	} else {
		d.Data = append(d.Data, PieChartDatum{
			Name:   s,
			ValInt: v,
		})
	}
}

func (d *PieChartData) Sort() {
	slices.SortFunc(d.Data, func(a, b PieChartDatum) int {
		if d.IsFloat {
			if n := cmp.Compare(a.ValFloat, b.ValFloat); n != 0 {
				return -1 * n
			} else {
				return cmp.Compare(a.Name, b.Name)
			}
		} else {
			if n := cmp.Compare(a.ValInt, b.ValInt); n != 0 {
				return -1 * n
			} else {
				return cmp.Compare(a.Name, b.Name)
			}
		}
	})
}

type PieChartDatum struct {
	Name     string
	ValInt   int64
	ValFloat float64
	IsFloat  bool
}

func (datum PieChartDatum) NameWithCount(defaultName string) string {
	var parts []string
	if len(datum.Name) > 0 {
		parts = append(parts, datum.Name)
	} else {
		parts = append(parts, defaultName)
	}
	if datum.IsFloat {
		parts = append(parts, fmt.Sprintf("(%s)", strconvutil.Ftoa(datum.ValFloat)))
	} else {
		parts = append(parts, fmt.Sprintf("(%d)", datum.ValInt))
	}
	return strings.Join(parts, " ")
}

/*
func (d PieChartDatum) JSONArray() ([]byte, error) {
	if d.IsFloat {
		var s = struct {
			n string
			v float64
		}{d.Name, d.ValFloat}
		return json.Marshal(s)
	} else {
		var s = struct {
			n string
			v int64
		}{d.Name, d.ValInt}
		return json.Marshal(s)
	}
}
*/
