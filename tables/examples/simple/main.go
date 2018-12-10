package main

import (
	"fmt"

	"github.com/grokify/gocharts/tables"
)

func main() {
	td := tables.TableData{
		Id:    "simpletable",
		Style: tables.StyleSimple,
		Rows: [][]string{
			[]string{"foo", "bar"},
			[]string{"1", "2"},
		},
	}
	output := tables.SimpleTable(td)
	fmt.Println(output)
}
