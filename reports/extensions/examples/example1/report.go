package main

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/grokify/gotilla/encoding/csvutil"

	"github.com/grokify/go-rickshaw/reports/extensions"
	"github.com/grokify/go-rickshaw/reports/extensions/templates"
)

func main() {
	inputfile := "data.csv"
	outputfile := "report.html"

	csv, fi, err := csvutil.NewReader(inputfile, rune(','), false)
	if err != nil {
		panic(fmt.Sprintf("ERROR %v\n", err))
	}

	rickshawData := rickshawextensions.NewRickshawData()

	idx := -1
	for {
		idx += 1
		record, err := csv.Read()
		if err == io.EOF {
			break
		}
		if idx == 0 {
			continue
		}
		monthData := rickshawextensions.MonthData{
			SeriesName: record[0],
			MonthS:     record[1],
			YearS:      record[2],
			ValueS:     record[3]}
		monthData.Inflate()

		item, err := monthData.RickshawItem()
		if err != nil {
			panic(fmt.Sprintf("ERR_BAD_RICKSHAW_ITEM: %v\n", err))
		}
		rickshawData.AddItem(item)

	}
	fi.Close()
	rickshawDataFormatted := rickshawData.Formatted()

	tmplData := rickshawextensions.TemplateData{
		ReportName:            "Fruit Report",
		RickshawURL:           "https://grokify.github.io/rickshaw",
		RickshawDataFormatted: rickshawDataFormatted,
		IncludeDataTable:      true}

	ioutil.WriteFile(outputfile, []byte(templates.RickshawExtensionsReport(tmplData)), 0644)

	fmt.Println("DONE")
}
