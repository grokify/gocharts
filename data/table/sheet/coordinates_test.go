package sheet

import (
	"testing"
)

var quotient26ToPrefixTests = []struct {
	quotient uint
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

var colNumberToLettersTests = []struct {
	colNumber     uint
	rowNumber     uint
	colIndex      uint
	rowIndex      uint
	prefix        string
	suffix        string
	letters       string
	sheetLocation string
}{
	{1, 1, 0, 0, "", "A", "A", "A1"},
	{36, 1, 35, 0, "A", "J", "AJ", "AJ1"},
	{26, 26, 25, 25, "", "Z", "Z", "Z26"},
	{1, 1024, 0, 1023, "", "A", "A", "A1024"},
	{26, 1024, 25, 1023, "", "Z", "Z", "Z1024"},
	{27, 1024, 26, 1023, "A", "A", "AA", "AA1024"},
	{52, 1024, 51, 1023, "A", "Z", "AZ", "AZ1024"},
	{53, 1024, 52, 1023, "B", "A", "BA", "BA1024"},
	{78, 1024, 77, 1023, "B", "Z", "BZ", "BZ1024"},
	{79, 2048, 78, 2047, "C", "A", "CA", "CA2048"},
	{80, 2048, 79, 2047, "C", "B", "CB", "CB2048"},
	{676, 2048, 675, 2047, "Y", "Z", "YZ", "YZ2048"},
	{677, 2048, 676, 2047, "Z", "A", "ZA", "ZA2048"},
	{702, 4096, 701, 4095, "Z", "Z", "ZZ", "ZZ4096"}}

type Instance struct {
	ColNumber   uint
	ColIndex    uint
	RowNumber   uint
	RowIndex    uint
	Quotient26  uint
	Remainder26 uint
	PrefixWant  string
	PrefixGot   string
	SuffixWant  string
	SuffixGot   string
	LettersWant string
	LettersGot  string
}

func TestColNumberToLetters(t *testing.T) {
	for _, tt := range colNumberToLettersTests {
		letters := ColNumberToLetters(tt.colNumber)
		inst := Instance{
			ColNumber:   tt.colNumber,
			RowNumber:   tt.rowNumber,
			ColIndex:    tt.colIndex,
			RowIndex:    tt.rowIndex,
			PrefixWant:  tt.prefix,
			SuffixWant:  tt.suffix,
			LettersWant: tt.letters,
			LettersGot:  letters}

		if inst.LettersGot != inst.LettersWant {
			t.Errorf("table.ColNumberToLetters(%v) want [%v] got [%v]",
				tt.colNumber, tt.letters, letters)
		}

		sheetLocIdx := CoordinatesToSheetLocation(tt.colIndex, tt.rowIndex)
		if sheetLocIdx != tt.sheetLocation {
			t.Errorf("table.CoordinatesToSheetLocation(%v,%v) want [%v] got [%v]",
				tt.colNumber, tt.rowNumber, tt.sheetLocation, sheetLocIdx)
		}

		sheetLocNum := CoordinateNumbersToSheetLocation(tt.colNumber, tt.rowNumber)
		if sheetLocNum != tt.sheetLocation {
			t.Errorf("table.CoordinateNumbersToSheetLocation(%v,%v) want [%v] got [%v]",
				tt.colNumber, tt.rowNumber, tt.sheetLocation, sheetLocNum)
		}
	}
}
