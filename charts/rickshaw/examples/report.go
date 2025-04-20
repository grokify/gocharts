package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/grokify/mogo/encoding/csvutil"

	"github.com/grokify/gocharts/v2/charts/rickshaw"
)

func main() {
	inputfile := "data.csv"
	outputfile := "report.html"

	csv, fi, err := csvutil.NewReaderFile(inputfile, rune(','))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	rickshawData := rickshaw.NewRickshawData()

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
		monthData := rickshaw.MonthData{
			SeriesName: record[0],
			MonthS:     record[1],
			YearS:      record[2],
			ValueS:     record[3]}
		if err := monthData.Inflate(); err != nil {
			slog.Error(err.Error())
			os.Exit(2)
		}

		if item, err := monthData.RickshawItem(); err != nil {
			slog.Error(err.Error(), "msg", "ERR_BAD_RICKSHAW_ITEM")
			os.Exit(3)
		} else {
			rickshawData.AddItem(item)
		}
	}
	fi.Close()
	rickshawDataFormatted, err := rickshawData.Formatted()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(4)
	}

	tmplData := rickshaw.TemplateData{
		ReportName:            "Fruit Report",
		RickshawURL:           "https://grokify.github.io/rickshaw",
		RickshawDataFormatted: rickshawDataFormatted,
		IncludeDataTable:      true}

	if err := os.WriteFile(outputfile,
		[]byte(rickshaw.RickshawExtensionsReport(tmplData)), 0600); err != nil {
		slog.Error(err.Error())
		os.Exit(5)
	}

	fmt.Println("DONE")
	os.Exit(0)
}
