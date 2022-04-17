package format

import "time"

func ParseTimeJanuary2006(s string) (time.Time, error) {
	return time.Parse("January 2006", s)
}
