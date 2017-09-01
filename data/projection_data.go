package data

type ProjectionDataInt struct {
	Target     int
	Current    int
	Projection int
	Delta      int
}

func (proj *ProjectionDataInt) Inflate() {
	proj.Delta = proj.Projection - proj.Target
}

func ProjectStraightLineInt64(startY int64, curY int64, startX int64, curX int64, endX int64) int64 {
	deltaCur := curX - startX
	deltaEnd := endX - startX
	return startY + int64(float64(curY)/float64(deltaCur)*float64(deltaEnd))
}
