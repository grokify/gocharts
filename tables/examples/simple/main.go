package main

import (
	"fmt"

	"github.com/grokify/gocharts/tables"
)

func main() {
	table := tables.TableData{
		Id:    "simpletable",
		Style: tables.StyleSimple,
		Rows: [][]string{
			{"foo", "bar"},
			{"1", "2"}}}

	output := tables.SimpleTable(table)
	fmt.Println(output)
}
