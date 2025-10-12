package table

import (
	"strings"
	"testing"
)

func TestMarkdownFormat(t *testing.T) {
	var markdownFormatTests = []struct {
		columns []string
		rows    [][]string
	}{
		{[]string{"a", "a"}, [][]string{{"bb", "bbb"}, {"cccc", "ccccc"}}},
	}

	for _, tt := range markdownFormatTests {
		tbl := NewTable("")
		tbl.Columns = tt.columns
		tbl.Rows = tt.rows
		if len(tbl.Rows)+1 < 2 {
			panic("bad test data")
		}
		md := tbl.Markdown("\n", false)
		lines := strings.Split(md, "\n")

		length := 0
		for i, l := range lines {
			if i == 0 {
				length = len(l)
			} else if len(l) != length {
				t.Errorf("Table.Markdown() Mismatch: first line length (%d) line (%d) length (%d)",
					length, i, len(l))
			}
		}
	}
}
