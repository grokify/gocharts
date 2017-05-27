package rickshawextensions

import (
	"encoding/json"
	"fmt"
)

type TemplateData struct {
	ReportName            string
	RickshawURL           string
	RickshawDataFormatted RickshawDataFormatted
	ItemType              string
	IncludeDataTable      bool
}

func (td *TemplateData) FormattedDataJSON() []byte {
	bytes, err := json.Marshal(td.RickshawDataFormatted.FormattedData)
	if err != nil {
		return []byte("")
	}
	return bytes
}

func (td *TemplateData) TableData() ([]string, [][]string) {
	dataRows := [][]string{}
	headRow := []string{td.ItemType}
	haveHeader := false
	if len(td.RickshawDataFormatted.FormattedData) > 0 {
		for _, series := range td.RickshawDataFormatted.FormattedData {
			if len(series.Data) > 0 {
				if !haveHeader {
					for _, item := range series.Data {
						headRow = append(
							headRow, fmt.Sprintf("%v %v",
								item.Time.Month().String()[0:3],
								item.Time.Year()))
					}
					haveHeader = true
				}
				dataRow := []string{series.Name}
				for _, item := range series.Data {
					dataRow = append(dataRow, fmt.Sprintf("%v", item.ValueY))
				}
				dataRows = append([][]string{dataRow}, dataRows...)
			}
		}
	}
	return headRow, dataRows
}
