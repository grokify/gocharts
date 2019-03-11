package histogram

type Histogram struct {
	BinCount      int            `json:"binCount"`
	BinsFrequency map[string]int `json:"binsFrequency"`
}

func NewHistogram() Histogram {
	return Histogram{BinsFrequency: map[string]int{}}
}

func (h *Histogram) Inflate() {
	h.BinCount = len(h.BinsFrequency)
}
