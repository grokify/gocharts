package interval

import (
	"time"
)

type Events struct {
	Events []Even
}

type Event interface {
	Time() time.Time
}

func (evs *Events) IntervalCount(dt time.Time, interval timeutil.Interval, dow time.Weekday) (int, error) {
	dtStart := timeutil.IntervalStart(dt, interval, dow)
	count := 0
	for _, ev := range evs.Events {
		dtThis := ev.Time()
		dtThisInt, err := timeutil.IntervalStart(dtThis, interval, dow)
		if dtThisInt == dtStart {
			count += 1
		}
	}
	return count
}

func (evs *Events) YoY(dt time.Time, interval timeutil.Interval, dow time.Weekday) float64 {
	dtStart, err := timeutil.IntervalStart(dt, interval, dow)
	thisCount := evs.IntervalCount(dtStart, interval, dow)
	dt

}
