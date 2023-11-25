package wchart

import (
	"time"

	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/time/quarter"
	"github.com/grokify/mogo/time/timeutil"
	chart "github.com/wcharczuk/go-chart/v2"
)

// Ticks converts a slice of `float64` to a slice of `chart.Tick`. Common
// formatting functions include `strconvutil.Commify` and
// `strconvutil.Int64Abbreviation`.
func Ticks(tickValues []float64, fn strconvutil.Float64ToString) []chart.Tick {
	ticks := []chart.Tick{}
	for _, tickVal := range tickValues {
		ticks = append(ticks, chart.Tick{
			Value: tickVal,
			Label: fn(tickVal)})
	}
	return ticks
}

// TicksInt64 converts a slice of `int64` to a slice of `chart.Tick`. Common
// formatting functions include `strconvutil.Commify` and
// `strconvutil.Int64Abbreviation`.
func TicksInt64(tickValues []int64, fn strconvutil.Int64ToString) []chart.Tick {
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
func GridLines(values []float64, style chart.Style) []chart.GridLine {
	lines := []chart.GridLine{}
	for _, val := range values {
		lines = append(lines, chart.GridLine{
			Style: style,
			Value: val})
	}
	return lines
}

// TicksAndGridlinesTime takes a start and end time and converts it to
// `[]chart.Tick` and `[]chart.GridLine.`.
func TicksAndGridlinesTime(interval timeutil.Interval, timeStart, timeEnd time.Time, styleMajor, styleMinor chart.Style, timeFormat func(time.Time) string, tickInterval, gridInterval timeutil.Interval) ([]chart.Tick, []chart.GridLine) {
	//fmt.Printf("TICK [%v] GRID [%v]\n", tickInterval, gridInterval)

	timecStart := uint64(0)
	timecEnd := uint64(0)
	if interval == timeutil.IntervalMonth {
		timecStart = month.TimeToMonthContinuous(timeStart)
		timecEnd = month.TimeToMonthContinuous(timeEnd)
	} else if interval == timeutil.IntervalQuarter {
		timecStart = quarter.TimeToQuarterContinuous(timeStart)
		timecEnd = quarter.TimeToQuarterContinuous(timeEnd)
	}

	ticks := []chart.Tick{}
	gridlines := []chart.GridLine{}
	if timecStart > timecEnd {
		tmp := timecStart
		timecStart = timecEnd
		timecEnd = tmp
	}
	for i := timecStart; i <= timecEnd; i++ {
		iTime := time.Now()
		if interval == timeutil.IntervalMonth {
			iTime = month.MonthContinuousToTime(i)
		} else if interval == timeutil.IntervalQuarter {
			iTime = quarter.QuarterContinuousToTime(i)
		}
		if i == timecStart {
			ticks = append(ticks, chart.Tick{Value: float64(i)})
			if (tickInterval == timeutil.IntervalQuarter && month.MonthContinuousIsQuarterStart(i)) ||
				(tickInterval == timeutil.IntervalYear && month.MonthContinuousIsYearStart(i)) {
				ticks[len(ticks)-1].Label = timeFormat(iTime)
			}
		} else if i == timecEnd {
			ticks = append(ticks, chart.Tick{Value: float64(i)})
			if (tickInterval == timeutil.IntervalQuarter && month.MonthContinuousIsQuarterStart(i)) ||
				(tickInterval == timeutil.IntervalYear && month.MonthContinuousIsYearStart(i)) {
				ticks[len(ticks)-1].Label = timeFormat(iTime)
			}
		} else {
			if (tickInterval == timeutil.IntervalQuarter && month.MonthContinuousIsQuarterStart(i)) ||
				(tickInterval == timeutil.IntervalYear && month.MonthContinuousIsYearStart(i)) ||
				(tickInterval != timeutil.IntervalQuarter && tickInterval != timeutil.IntervalYear) {
				ticks = append(ticks, chart.Tick{Value: float64(i)})
				ticks[len(ticks)-1].Label = timeFormat(iTime)
			}
			if (gridInterval == timeutil.IntervalQuarter && month.MonthContinuousIsQuarterStart(i)) ||
				(gridInterval == timeutil.IntervalYear && month.MonthContinuousIsYearStart(i)) ||
				(gridInterval != timeutil.IntervalQuarter && gridInterval != timeutil.IntervalYear) {
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

		if 1 == 0 {
			if interval == timeutil.IntervalMonth && month.MonthContinuousIsYearStart(i) {
				//if interval == timeutil.Month && month.MonthContinuousIsQuarterStart(i) {
				ticks = append(ticks, chart.Tick{Value: float64(i)})
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
			} else if interval == timeutil.IntervalQuarter {
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
		} else if 1 == 2 { // monthly
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
