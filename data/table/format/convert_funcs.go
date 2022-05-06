package format

import (
	"regexp"
	"strings"
)

func ConvertDecommify(s string) (string, error) {
	return strings.Replace(s, ",", "", -1), nil
}

func ConvertRemoveControls(s string) (string, error) {
	// Tableau CSVs may have `\x00` control chars.
	// use `Table.ForamatRows()` on all columns to remove.
	return regexp.MustCompile(`\x00`).ReplaceAllString(s, ""), nil
}
