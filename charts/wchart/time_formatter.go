package wchart

import (
	"strings"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

// TimeFormatter provides a struct that satisifies the
// `github.com/wcharczuk/go-chart.ValueFormatter` interface
// using a given time layout string.
type TimeFormatter struct {
	Layout string // time format string
}

func (tvf *TimeFormatter) FormatTime(v any) string {
	tvf.Layout = strings.TrimSpace(tvf.Layout)
	if len(tvf.Layout) == 0 {
		tvf.Layout = time.RFC3339
	}
	return timeutil.FormatTimeMulti(tvf.Layout, v)
}
