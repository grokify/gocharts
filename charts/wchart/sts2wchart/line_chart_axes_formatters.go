package sts2wchart

import (
	"fmt"
	"regexp"

	"github.com/grokify/mogo/strconv/strconvutil"
)

func YAxisTickFormatSimple(raw float64) string {
	return strconvutil.Int64Abbreviation(int64(raw))
}

func YAxisTickFormatPercent(raw float64) string {
	return fmt.Sprintf("%.1f%%", raw*100)
}

func YAxisTickFormatDollars(raw float64) string {
	abbr := strconvutil.Int64Abbreviation(int64(raw))
	return "$" + abbr
}

var rxMrr = regexp.MustCompile(`(?i)\bmrr\b`)

func FormatYTickFunc(seriesName string) func(float64) string {
	return func(val float64) string {
		abbr := strconvutil.Int64Abbreviation(int64(val))
		if rxMrr.MatchString(seriesName) {
			return "$" + abbr
		}
		return abbr
	}
}
