package sheet

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ExcelMaxColCount = 16384
	ExcelMaxRowCount = 1048576
	Alphabet         = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ZZ               = uint32(702)
	ZZIndex          = uint32(701)
)

var alphabetSlice = strings.Split(Alphabet, "")

func quotient26ToPrefix(quotient uint32) string {
	// " ABCDEF"
	// "0123456"
	// Max Quotient =26
	if quotient > 27 {
		panic(fmt.Sprintf("quotient26 out of range [0..26][%v]", quotient))
	} else if quotient == 0 {
		return ""
	}
	prefix := alphabetSlice[quotient-1] // A=1
	return prefix
}

/*
func quotient25ToPrefix(quotient uint32) string {
	// " ABCDEF"
	// "0123456"
	// Max Quotient =26
	if quotient > 26 {
		panic(fmt.Sprintf("quotient25 out of range [0..26][%v]", quotient))
	} else if quotient == 0 {
		return ""
	}
	prefix := alphabetSlice[quotient] // A=1
	return prefix
}
*/

func remainder26ToSuffix(remainder uint32) string {
	if remainder == 0 {
		return "Z"
	}
	if remainder > 25 {
		panic(fmt.Sprintf("remainder out of range: [%v]", remainder))
	}
	letter := alphabetSlice[remainder-1]
	return letter
}

func remainder25ToSuffix(remainder uint32) string {
	if remainder > 25 {
		panic(fmt.Sprintf("remainder out of range: [%v]", remainder))
	}
	letter := alphabetSlice[remainder]
	return letter
}

func ColNumberToLetters(colNumber uint32) string {
	if colNumber == 0 {
		panic("row cannot be zero. Row is 1 indexed.")
	} else if colNumber > ZZ {
		panic("colNumber greater than 702 (ZZ) is not currently supported")
	}
	quotient := uint32(0)
	if colNumber > 26 {
		quotient = (colNumber - 1) / 26
	}
	remainder := colNumber % 26
	if colNumber < 26 {
		remainder = colNumber
	}
	prefix := quotient26ToPrefix(quotient)
	suffix := remainder26ToSuffix(remainder)
	return strings.TrimSpace(prefix) + strings.TrimSpace(suffix)
}

func ColIndexToLetters(colIndex uint32) string {
	if colIndex > ZZIndex {
		panic(
			fmt.Sprintf("colIndex [%v] greater than index:701 (ZZ) is not currently supported",
				colIndex))
	}
	quotient := uint32(0)
	if colIndex > 25 {
		quotient = colIndex / 26
	}
	remainder := colIndex % 26
	prefix := quotient26ToPrefix(quotient)
	suffix := remainder25ToSuffix(remainder)
	return strings.TrimSpace(prefix) + strings.TrimSpace(suffix)
}

func ColLettersToIndex(letters string) (uint32, error) {
	letters = strings.ToUpper(strings.TrimSpace(letters))
	var result uint32 = 0
	for i := range len(letters) {
		c := letters[i]
		if c < 'A' || c > 'Z' {
			return 0, fmt.Errorf("invalid column letter (%q)", c)
		}
		result = result*26 + uint32(c-'A'+1)
	}
	return result, nil
}

// CoordinatesToSheetLocation converts x, y integer coordinates
// to a spreadsheet location such as "AA1" for col 27, row 1.
func CoordinatesToSheetLocation(colIdx, rowIdx uint32) string {
	colLet := ColIndexToLetters(colIdx)
	return colLet + strconv.Itoa(int(rowIdx+1))
}

// CoordinateNumbersToSheetLocation converts x, y integer coordinates
// to a spreadsheet location such as "AA1" for col 27, row 1.
func CoordinateNumbersToSheetLocation(colNum, rowNum uint32) string {
	return CoordinatesToSheetLocation(colNum-1, rowNum-1)
	//colLet := ColNumberToLetters(colNum)
	//return colLet + strconv.Itoa(int(rowNum))
}
