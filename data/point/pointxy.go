package point

type PointXYs []PointXY

type PointXY struct {
	X float64
	Y float64
}

func (pts PointXYs) XAndYSeries() ([]float64, []float64) {
	var xs []float64
	var ys []float64
	for _, pt := range pts {
		xs = append(xs, pt.X)
		ys = append(ys, pt.Y)
	}
	return xs, ys
}
