package main

import (
	"fmt"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gocharts/tables"
)

func main() {
	tbl := table.Table{
		ID:    "simpletable",
		Style: tables.StyleSimple,
		Records: [][]string{
			{"foo", "bar"},
			{"1", "2"}}}

	output := table.SimpleTable(tbl)
	fmt.Println(output)
}
