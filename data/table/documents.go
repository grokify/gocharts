package table

import (
	"github.com/grokify/gocharts/data/histogram"
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
	hg := histogram.NewHistogram()

	//histogram := map[string]int{}
	for _, doc := range ds.Documents {
		if val, ok := doc[key]; ok {
			if _, ok := hg.BinsFrequency[val]; !ok {
				hg.BinsFrequency[val] = 0
			}
			hg.BinsFrequency[val] += 1
		}
	}
	hg.Inflate()
	if ds.Meta.Histograms == nil {
		ds.Meta.Histograms = map[string]histogram.Histogram{}
	}
	ds.Meta.Histograms[key] = hg
}

type DocumentsSetMeta struct {
	Count      int                            `json:"count"`
	Histograms map[string]histogram.Histogram `json:"histograms"`
}

func NewDocumentsSetMeta() DocumentsSetMeta {
	return DocumentsSetMeta{
		Histograms: map[string]histogram.Histogram{}}
}
