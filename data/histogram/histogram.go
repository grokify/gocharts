package histogram

type Histogram struct {
	BinsCount          int            `json:"binsCount"`
	BinFrequencyCounts map[string]int `json:"binFrequencyCounts"`
}

func NewHistogram() Histogram {
	return Histogram{BinFrequencyCounts: map[string]int{}}
}

func (h *Histogram) Inflate() {
	h.BinsCount = len(h.BinFrequencyCounts)
}
