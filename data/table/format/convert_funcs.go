package format

import (
	"regexp"
)

var rxComma = regexp.MustCompile(`,`)

func ConvertDecommify(input string) (string, error) {
	return rxComma.ReplaceAllString(input, ""), nil
}
