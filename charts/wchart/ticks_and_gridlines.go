package wchart

import (
	"time"

	"github.com/grokify/gotilla/strconv/strconvutil"
	"github.com/grokify/gotilla/time/month"
	"github.com/grokify/gotilla/time/quarter"
	"github.com/grokify/gotilla/time/timeutil"
	"github.com/wcharczuk/go-chart"
)

// Ticks converts a slice of `int64` to a slice of `chart.Tick`. Common
// formatting functions include `strconvutil.Commify` and
// `strconvutil.Int64Abbreviation`.
func Ticks(tickValues []int64, fn strconvutil.Int64ToString) []chart.Tick {
	ticks := []chart.Tick{}
	for _, tickVal := range tickValues {
		ticks = append(ticks, chart.Tick{
			Value: float64(tickVal),
			Label: fn(tickVal)})
	}
	return ticks
}

// GridLines creates a `[]chart.GridLine` from a slice of `int64`
// and a style.
func GridLines(values []int64, style chart.Style) []chart.GridLine {
	lines := []chart.GridLine{}
	for _, val := range values {
		lines = append(lines, chart.GridLine{
			Style: style,
			Value: float64(val)})
	}
	return lines
}

// TicksAndGridlinesTime takes a begin and end time and converts it to
// `[]chart.Tick` and `[]chart.GridLine.`.
func TicksAndGridlinesTime(interval timeutil.Interval, timeBegin, timeEnd time.Time, styleMajor, styleMinor chart.Style, timeFormat func(time.Time) string, quarterOnly bool) ([]chart.Tick, []chart.GridLine) {
	timecBegin := uint64(0)
	timecEnd := uint64(0)
	if interval == timeutil.Month {
		timecBegin = month.TimeToMonthContinuous(timeBegin)
		timecEnd = month.TimeToMonthContinuous(timeEnd)
	} else if interval == timeutil.Quarter {
		timecBegin = quarter.TimeToQuarterContinuous(timeBegin)
		timecEnd = quarter.TimeToQuarterContinuous(timeEnd)
	}

	ticks := []chart.Tick{}
	gridlines := []chart.GridLine{}
	if timecBegin > timecEnd {
		tmp := timecBegin
		timecBegin = timecEnd
		timecEnd = tmp
	}
	for i := timecBegin; i <= timecEnd; i++ {
		iTime := time.Now()
		if interval == timeutil.Month {
			iTime = month.MonthContinuousToTime(i)
		} else if interval == timeutil.Quarter {
			iTime = quarter.QuarterContinuousToTime(i)
		}
		if i == timecBegin {
			ticks = append(ticks, chart.Tick{Value: float64(i)})
		} else if i == timecEnd {
			ticks = append(ticks, chart.Tick{Value: float64(i)})
		} else if quarterOnly {
			if interval == timeutil.Month && month.MonthContinuousIsQuarterBegin(i) {
				ticks = append(ticks, chart.Tick{
					Value: float64(i)})
				ticks[len(ticks)-1].Label = timeFormat(iTime)
				if iTime.Month() == 1 {
					gridlines = append(gridlines, chart.GridLine{
						Style: styleMajor,
						Value: float64(i)})
				} else {
					gridlines = append(gridlines, chart.GridLine{
						Style: styleMinor,
						Value: float64(i)})
				}
			} else if interval == timeutil.Quarter {
				ticks = append(ticks, chart.Tick{
					Value: float64(i)})
				ticks[len(ticks)-1].Label = timeFormat(iTime)
				if iTime.Month() == 1 {
					gridlines = append(gridlines, chart.GridLine{
						Style: styleMajor,
						Value: float64(i)})
				} else {
					gridlines = append(gridlines, chart.GridLine{
						Style: styleMinor,
						Value: float64(i)})
				}
			}
		} else { // monthly
			ticks = append(ticks, chart.Tick{
				Value: float64(i)})
			ticks[len(ticks)-1].Label = timeFormat(iTime)
			if iTime.Month() == 1 {
				gridlines = append(gridlines, chart.GridLine{
					Style: styleMajor,
					Value: float64(i)})
			} else {
				gridlines = append(gridlines, chart.GridLine{
					Style: styleMinor,
					Value: float64(i)})
			}
		}
	}

	return ticks, gridlines
}
