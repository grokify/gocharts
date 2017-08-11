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
