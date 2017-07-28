package c3

import (
	"encoding/json"
)

type TemplateData struct {
	HeaderHTML             string
	ReportName             string
	ReportLink             string
	IncludeDataTable       bool
	IncludeDataTableTotals bool
	C3Chart                C3Chart
}

func (td *TemplateData) FormattedDataJSON() []byte {
	bytes, err := json.Marshal(td.C3Chart)
	if err != nil {
		return []byte("")
	}
	return bytes
}
