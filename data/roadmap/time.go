package roadmap

import (
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/ordered"
)

// QuartersBeginEnd converts relative and
// default quarter times to absolute int32
// quarter numbers.
func QuartersBeginEnd(begin, end int) (int, int) {
	if begin == 0 && end == 0 {
		begin = -1
		end = 4
	}
	return ordered.MinMax(
		timeutil.QuartersRelToAbs(begin, end))
}
