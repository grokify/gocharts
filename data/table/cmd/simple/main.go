package main

import (
	"fmt"

	"github.com/grokify/gocharts/v2/data/table"
)

func main() {
	tbl := table.Table{
		ID:    "simpletable",
		Style: table.StyleSimple,
		Rows: [][]string{
			{"foo", "bar"},
			{"1", "2"}}}

	output := table.SimpleTable(tbl)
	fmt.Println(output)
}
