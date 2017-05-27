package rickshawextensions

import (
	"encoding/json"
)

type TemplateData struct {
	ReportName            string
	RickshawURL           string
	RickshawDataFormatted RickshawDataFormatted
}

func (td *TemplateData) FormattedDataJSON() []byte {
	bytes, err := json.Marshal(td.RickshawDataFormatted.FormattedData)
	if err != nil {
		return []byte("")
	}
	return bytes
}
