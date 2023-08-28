package histogram

type HistogramAny interface {
	BinNames() []string
	ItemCount() uint
	ItemNames() []string
	Sum() int
}
