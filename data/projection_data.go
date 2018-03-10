package data

import (
	"fmt"
	"strings"
	"time"
)

type ProjectionDataInt struct {
	Start      int
	Current    int
	Target     int
	Projection int
	Delta      int
}

func (proj *ProjectionDataInt) Inflate() {
	proj.Delta = proj.Projection - proj.Target
}

type ProjectionDataInt64 struct {
	Start      int64
	Current    int64
	Target     int64
	Projection int64
	Delta      int64
}

func (proj *ProjectionDataInt64) CreateProjection(start, current, end int64) {
	proj.Projection = ProjectStraightLineInt64(
		proj.Start, proj.Current,
		start, current, end)
	proj.Inflate()
}

func (proj *ProjectionDataInt64) CreateProjectionTime(start, current, end time.Time) {
	proj.Projection = ProjectStraightLineInt64(
		proj.Start, proj.Current,
		start.UTC().Unix(), current.UTC().Unix(), end.UTC().Unix())
	proj.Inflate()
}

func (proj *ProjectionDataInt64) ToString(strs []string, abbr bool) string {
	outputs := []string{}
	for _, str := range strs {
		switch strings.ToLower(str) {
		case "s":
			if abbr {
				outputs = append(outputs, fmt.Sprintf("S:%v", proj.Start))
			} else {
				outputs = append(outputs, fmt.Sprintf("Start: %v", proj.Start))
			}
		case "c":
			if abbr {
				outputs = append(outputs, fmt.Sprintf("C:%v", proj.Current))
			} else {
				outputs = append(outputs, fmt.Sprintf("Current: %v", proj.Current))
			}
		case "t":
			if abbr {
				outputs = append(outputs, fmt.Sprintf("T:%v", proj.Target))
			} else {
				outputs = append(outputs, fmt.Sprintf("Target: %v", proj.Target))
			}
		case "p":
			if abbr {
				outputs = append(outputs, fmt.Sprintf("P:%v", proj.Projection))
			} else {
				outputs = append(outputs, fmt.Sprintf("Projection: %v", proj.Projection))
			}
		case "d":
			if abbr {
				outputs = append(outputs, fmt.Sprintf("D:%v", proj.Delta))
			} else {
				outputs = append(outputs, fmt.Sprintf("Delta: %v", proj.Delta))
			}
		}
	}
	if abbr {
		return strings.Join(outputs, " ")
	}
	return strings.Join(outputs, ", ")
}

func (proj *ProjectionDataInt64) Inflate() {
	proj.Delta = proj.Projection - proj.Target
}

func ProjectStraightLineInt64(startY, curY, startX, curX, endX int64) int64 {
	deltaCur := curX - startX
	deltaEnd := endX - startX
	return startY + int64(float64(curY)/float64(deltaCur)*float64(deltaEnd))
}
