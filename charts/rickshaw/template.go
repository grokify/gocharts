package rickshaw

import (
	"encoding/json"
	"fmt"
)

type TemplateData struct {
	HeaderHTML             string
	ReportName             string
	ReportLink             string
	RickshawURL            string
	RickshawDataFormatted  RickshawDataFormatted
	ItemType               string
	IncludeDataTable       bool
	IncludeDataTableTotals bool
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
		allSeriesSubtotal := int64(0)
		for _, series := range td.RickshawDataFormatted.FormattedData {
			if len(series.Data) > 0 {
				if !haveHeader {
					for _, item := range series.Data {
						dt := item.Time.UTC()
						headRow = append(
							headRow, fmt.Sprintf("%v %v",
								dt.Month().String()[0:3],
								dt.Year()))
					}
					haveHeader = true
				}
				dataRow := []string{series.Name}
				seriesYSubtotal := int64(0)
				for _, item := range series.Data {
					dataRow = append(dataRow, fmt.Sprintf("%v", item.ValueY))
					seriesYSubtotal += item.ValueY
				}
				if td.IncludeDataTableTotals {
					dataRow = append(dataRow, fmt.Sprintf("%v", seriesYSubtotal))
					allSeriesSubtotal += seriesYSubtotal
				}
				dataRows = append([][]string{dataRow}, dataRows...)
			}
		}
		if td.IncludeDataTableTotals {
			headRow = append(headRow, "Total")
			dataRow := []string{"Total"}
			for i := 0; i < len(headRow)-2; i++ {
				dataRow = append(dataRow, "")
			}
			dataRow = append(dataRow, fmt.Sprintf("%v", allSeriesSubtotal))
			dataRows = append(dataRows, dataRow)
		}
	}
	return headRow, dataRows
}
