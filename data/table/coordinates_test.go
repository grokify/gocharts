package table

import (
	"testing"
)

var quotient26ToPrefixTests = []struct {
	quotient uint32
	prefix   string
}{
	{0, ""},
	{1, "A"},
	{2, "B"},
	{3, "C"}}

func TestQuotient26ToPrefix(t *testing.T) {
	for _, tt := range quotient26ToPrefixTests {
		prefix := quotient26ToPrefix(tt.quotient)
		if prefix != tt.prefix {
			t.Errorf("table.Quotient26ToPrefix: with [%v] want [%v] got [%v]",
				tt.quotient, tt.prefix, prefix)
		}
	}
}

var colNumberToLetters = []struct {
	colNumber     uint32
	rowNumber     uint32
	prefix        string
	suffix        string
	letters       string
	sheetLocation string
}{
	{1, 1024, "", "A", "A", "A1024"},
	{26, 1024, "", "Z", "Z", "Z1024"},
	{27, 1024, "A", "A", "AA", "AA1024"},
	{52, 1024, "A", "Z", "AZ", "AZ1024"},
	{53, 1024, "B", "A", "BA", "BA1024"},
	{78, 1024, "B", "Z", "BZ", "BZ1024"},
	{79, 2048, "C", "A", "CA", "CA2048"},
	{80, 2048, "C", "B", "CB", "CB2048"},
	{676, 2048, "Y", "Z", "YZ", "YZ2048"},
	{677, 2048, "Z", "A", "ZA", "ZA2048"},
	{702, 4096, "Z", "Z", "ZZ", "ZZ4096"}}

type Instance struct {
	ColNumber   uint32
	Quotient26  uint32
	Remainder26 uint32
	PrefixWant  string
	PrefixGot   string
	SuffixWant  string
	SuffixGot   string
	LettersWant string
	LettersGot  string
}

func TestColNumberToLetters(t *testing.T) {
	for _, tt := range colNumberToLetters {
		letters := ColNumberToLetters(tt.colNumber)
		inst := Instance{
			ColNumber:   tt.colNumber,
			PrefixWant:  tt.prefix,
			SuffixWant:  tt.suffix,
			LettersWant: tt.letters,
			LettersGot:  letters}
		//fmtutil.PrintJSON(inst)
		if inst.LettersGot != inst.LettersWant {
			t.Errorf("table.ColNumberToLetters(%v) want [%v] got [%v]",
				tt.colNumber, tt.letters, letters)
		}
		sheetLoc := CoordinatesToSheetLocation(tt.colNumber, tt.rowNumber)
		if sheetLoc != tt.sheetLocation {
			t.Errorf("table.CoordinatesToSheetLocation(%v,%v) want [%v] got [%v]",
				tt.colNumber, tt.rowNumber, tt.sheetLocation, sheetLoc)
		}
	}
}
