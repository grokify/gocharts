package histogram

type HistogramAny interface {
	ItemCount() uint
	Sum() int
}
