package roadmap

import (
	"github.com/grokify/mogo/math/mathutil"
	"github.com/grokify/mogo/time/timeutil"
)

// QuarterInt32sBeginEnd converts relative and
// default quarter times to absolute int32
// quarter numbers.
func QuarterInt32sBeginEnd(begin, end int32) (int32, int32) {
	if begin == 0 && end == 0 {
		begin = -1
		end = 4
	}
	return mathutil.MinMaxInt32(
		timeutil.QuartersInt32RelToAbs(begin, end))
}
