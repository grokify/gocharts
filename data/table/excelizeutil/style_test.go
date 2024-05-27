package excelizeutil

import (
	"testing"
)

var readFileTests = []struct {
	filename string
}{
	{"testdata/style_test.xlsx"},
}

// TestReadFile reads a XSLX file.
func TestReadFile(t *testing.T) {
	for _, tt := range readFileTests {
		xlsx, err := NewFile(tt.filename)
		if err != nil {
			t.Errorf("excelizeutil.NewFile(\"%s\",...) Error: [%v]", tt.filename, err.Error())
		}
		cellStyleA1, err := xlsx.File.GetCellStyle("Sheet1", "A2")
		if err != nil {
			t.Errorf(`xlsx.File.GetCellStyle("Sheet1","A1") Error: [%v]`, err.Error())
		}
		// fmt.Printf("STYLE (%v)\n", cellStyleA1)

		cellStyle, err := xlsx.File.GetStyle(cellStyleA1)
		if err != nil {
			t.Errorf(`xlsx.File.GetStyle("Sheet1","A1") Error: [%v]`, err.Error())
		}
		if cellStyle == nil {
			t.Errorf(`xlsx.File.GetStyle("Sheet1","A1") [%s]`, "nil returned")
		}
		// fmtutil.PrintJSON(cellStyle)
	}
}
