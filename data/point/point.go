package point

import (
	"sort"
)

type PointSet struct {
	IsFloat   bool
	PointsMap map[string]Point
}

func NewPointSet() PointSet {
	return PointSet{PointsMap: map[string]Point{}}
}

func (ps *PointSet) Inflate() {
	total := ps.TotalValueFloat()
	for key, point := range ps.PointsMap {
		if ps.IsFloat && !point.IsFloat {
			point.IsFloat = true
			point.AbsoluteFloat = float64(point.AbsoluteInt)
		} else if !ps.IsFloat && point.IsFloat {
			point.IsFloat = false
			point.AbsoluteInt = int64(point.AbsoluteFloat)
		}
		if point.IsFloat {
			point.Percentage = (point.AbsoluteFloat / total) * 100
		} else {
			point.Percentage = (float64(point.AbsoluteInt) / total) * 100
		}
		point.PercentageNot = 100.0 - point.Percentage
		ps.PointsMap[key] = point
	}
}

func (ps *PointSet) TotalValueFloat() float64 {
	total := float64(0)
	for _, point := range ps.PointsMap {
		if point.IsFloat {
			total += point.AbsoluteFloat
		} else {
			total += float64(point.AbsoluteInt)
		}
	}
	return total
}

func (ps *PointSet) Slice(sortDesc bool) []Point {
	points := []Point{}
	for _, point := range ps.PointsMap {
		points = append(points, point)
	}
	if sortDesc {
		if ps.IsFloat {
			sort.Slice(points, func(i, j int) bool {
				return points[i].AbsoluteInt > points[j].AbsoluteInt
			})
		} else {
			sort.Slice(points, func(i, j int) bool {
				return points[i].AbsoluteFloat > points[j].AbsoluteFloat
			})
		}
	}
	return points
}

type Point struct {
	Name          string
	IsFloat       bool
	AbsoluteInt   int64
	AbsoluteFloat float64
	Percentage    float64
	PercentageNot float64
}
