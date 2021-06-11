package table

import (
	"testing"
)

var readFileTests = []struct {
	filename   string
	columns    []string
	colNameLen []int
}{
	{"testdata/simple.csv", []string{"bazqux", "foobar"}, []int{6, 6}},
	{"testdata/utf8bom.csv", []string{"foobar", "bazqux"}, []int{6, 6}},
}

// TestReadFile reads a CSV file.
func TestReadFile(t *testing.T) {
	for _, tt := range readFileTests {
		tbl, err := ReadFile(tt.filename, ',', true)
		if err != nil {
			t.Errorf("table.ReadFile(\"%s\",...) Error: [%v]",
				tt.filename, err.Error())
		}
		if len(tbl.Columns) != len(tt.columns) {
			t.Errorf("table.ReadFile(\"%s\",...) Column Number mismatch: want [%d] got [%d]",
				tt.filename, len(tt.columns), len(tbl.Columns))
		}

		for i, wantColName := range tt.columns {
			tryColName := tbl.Columns[i]
			if wantColName != tryColName {
				t.Errorf("table.ReadFile(\"%s\",...) Column Name mismatch: want [%s] got [%s]",
					tt.filename, wantColName, tryColName)
			}
			if len(tryColName) != tt.colNameLen[i] {
				t.Errorf("table.ReadFile(\"%s\",...) Column Name length mismatch: want [%d] got [%d]",
					tt.filename, tt.colNameLen[i], len(tryColName))
			}
		}
	}
}
