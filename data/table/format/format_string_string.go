package format

import "regexp"

func FormatStringRemoveControls(s string) (string, error) {
	// Tableau CSVs may have `\x00` control chars.
	// use `Table.ForamatRows()` on all columns to remove.
	return regexp.MustCompile(`\x00`).ReplaceAllString(s, ""), nil
}
