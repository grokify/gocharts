package wchart

import (
	"time"

	"github.com/grokify/gotilla/strconv/strconvutil"
	"github.com/grokify/gotilla/time/month"
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

// TicksAndGridlinesMonths takes a begin and end time and converts it to
// `[]chart.Tick` and `[]chart.GridLine.`.
func TicksAndGridlinesMonths(timeBegin, timeEnd time.Time, styleMajor, styleMinor chart.Style, timeFormat string, quarterOnly bool) ([]chart.Tick, []chart.GridLine) {
	monthcBegin := month.TimeToMonthContinuous(timeBegin)
	monthcEnd := month.TimeToMonthContinuous(timeEnd)
	ticks := []chart.Tick{}
	gridlines := []chart.GridLine{}
	if monthcBegin > monthcEnd {
		tmp := monthcBegin
		monthcBegin = monthcEnd
		monthcEnd = tmp
	}
	for i := monthcBegin; i <= monthcEnd; i++ {
		iTime := month.MonthContinuousToTime(i)
		if i == monthcBegin {
			ticks = append(ticks, chart.Tick{Value: float64(i)})
		} else if i == monthcEnd {
			ticks = append(ticks, chart.Tick{Value: float64(i)})
		} else if quarterOnly {
			if month.MonthContinuousIsQuarterBegin(i) {
				ticks = append(ticks, chart.Tick{
					Value: float64(i)})
				if len(timeFormat) > 0 {
					ticks[len(ticks)-1].Label = iTime.Format(timeFormat)
				}
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
			if len(timeFormat) > 0 {
				ticks[len(ticks)-1].Label = iTime.Format(timeFormat)
			}
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
