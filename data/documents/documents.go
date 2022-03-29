package documents

import (
	"github.com/grokify/gocharts/v2/data/histogram"
)

type DocumentsSet struct {
	Meta      DocumentsSetMeta    `json:"meta"`
	Documents []map[string]string `json:"records"`
}

func NewDocumentsSet() DocumentsSet {
	return DocumentsSet{
		Meta:      DocumentsSetMeta{},
		Documents: []map[string]string{}}
}

func (ds *DocumentsSet) Inflate() {
	ds.Meta.Count = len(ds.Documents)
}

func (ds *DocumentsSet) CreateHistogram(key string) {
	hist := histogram.NewHistogram("")

	//histogram := map[string]int{}
	for _, doc := range ds.Documents {
		if val, ok := doc[key]; ok {
			if _, ok := hist.Bins[val]; !ok {
				hist.Bins[val] = 0
			}
			hist.Bins[val] += 1
		}
	}
	hist.Inflate()
	if ds.Meta.Histograms == nil {
		ds.Meta.Histograms = map[string]*histogram.Histogram{}
	}
	ds.Meta.Histograms[key] = hist
}

type DocumentsSetMeta struct {
	Count      int                             `json:"count"`
	Histograms map[string]*histogram.Histogram `json:"histograms"`
}

func NewDocumentsSetMeta() DocumentsSetMeta {
	return DocumentsSetMeta{
		Histograms: map[string]*histogram.Histogram{}}
}
