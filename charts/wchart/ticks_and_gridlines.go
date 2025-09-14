package wchart

import (
	"time"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/time/quarter"
	"github.com/grokify/mogo/time/timeutil"
)

// Ticks converts a slice of `float64` to a slice of `chartdraw.Tick`. Common
// formatting functions include `strconvutil.Commify` and
// `strconvutil.Int64Abbreviation`.
func Ticks(tickValues []float64, fn strconvutil.Float64ToString) []chartdraw.Tick {
	ticks := []chartdraw.Tick{}
	for _, tickVal := range tickValues {
		ticks = append(ticks, chartdraw.Tick{
			Value: tickVal,
			Label: fn(tickVal)})
	}
	return ticks
}

// TicksInt64 converts a slice of `int64` to a slice of `chartdraw.Tick`. Common
// formatting functions include `strconvutil.Commify` and
// `strconvutil.Int64Abbreviation`.
func TicksInt64(tickValues []int64, fn strconvutil.Int64ToString) []chartdraw.Tick {
	ticks := []chartdraw.Tick{}
	for _, tickVal := range tickValues {
		ticks = append(ticks, chartdraw.Tick{
			Value: float64(tickVal),
			Label: fn(tickVal)})
	}
	return ticks
}

// GridLines creates a `[]chartdraw.GridLine` from a slice of `int64`
// and a style.
func GridLines(values []float64, style chartdraw.Style) []chartdraw.GridLine {
	lines := []chartdraw.GridLine{}
	for _, val := range values {
		lines = append(lines, chartdraw.GridLine{
			Style: style,
			Value: val})
	}
	return lines
}

// TicksAndGridlinesTime takes a start and end time and converts it to
// `[]chartdraw.Tick` and `[]chartdraw.GridLine.`.
func TicksAndGridlinesTime(interval timeutil.Interval, timeStart, timeEnd time.Time, styleMajor, styleMinor chartdraw.Style, timeFormat func(time.Time) string, tickInterval, gridInterval timeutil.Interval) ([]chartdraw.Tick, []chartdraw.GridLine, error) {
	//fmt.Printf("TICK [%v] GRID [%v]\n", tickInterval, gridInterval)
	ticks := []chartdraw.Tick{}
	gridlines := []chartdraw.GridLine{}

	timecStart := uint32(0)
	timecEnd := uint32(0)
	if interval == timeutil.IntervalMonth {
		if timecStartTry, err := month.TimeToMonthContinuous(timeStart); err != nil {
			return ticks, gridlines, err
		} else {
			timecStart = timecStartTry
		}
		if timecEndTry, err := month.TimeToMonthContinuous(timeStart); err != nil {
			return ticks, gridlines, err
		} else {
			timecEnd = timecEndTry
		}
	} else if interval == timeutil.IntervalQuarter {
		if timecStartTry, err := quarter.TimeToQuarterContinuous(timeStart); err != nil {
			return ticks, gridlines, err
		} else {
			timecStart = timecStartTry
		}
		if timecEndTry, err := quarter.TimeToQuarterContinuous(timeEnd); err != nil {
			return ticks, gridlines, err
		} else {
			timecEnd = timecEndTry
		}
	}

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
			ticks = append(ticks, chartdraw.Tick{Value: float64(i)})
			if (tickInterval == timeutil.IntervalQuarter && month.MonthContinuousIsQuarterStart(i)) ||
				(tickInterval == timeutil.IntervalYear && month.MonthContinuousIsYearStart(i)) {
				ticks[len(ticks)-1].Label = timeFormat(iTime)
			}
		} else if i == timecEnd {
			ticks = append(ticks, chartdraw.Tick{Value: float64(i)})
			if (tickInterval == timeutil.IntervalQuarter && month.MonthContinuousIsQuarterStart(i)) ||
				(tickInterval == timeutil.IntervalYear && month.MonthContinuousIsYearStart(i)) {
				ticks[len(ticks)-1].Label = timeFormat(iTime)
			}
		} else {
			if (tickInterval == timeutil.IntervalQuarter && month.MonthContinuousIsQuarterStart(i)) ||
				(tickInterval == timeutil.IntervalYear && month.MonthContinuousIsYearStart(i)) ||
				(tickInterval != timeutil.IntervalQuarter && tickInterval != timeutil.IntervalYear) {
				ticks = append(ticks, chartdraw.Tick{Value: float64(i)})
				ticks[len(ticks)-1].Label = timeFormat(iTime)
			}
			if (gridInterval == timeutil.IntervalQuarter && month.MonthContinuousIsQuarterStart(i)) ||
				(gridInterval == timeutil.IntervalYear && month.MonthContinuousIsYearStart(i)) ||
				(gridInterval != timeutil.IntervalQuarter && gridInterval != timeutil.IntervalYear) {
				if iTime.Month() == 1 {
					gridlines = append(gridlines, chartdraw.GridLine{
						Style: styleMajor,
						Value: float64(i)})
				} else {
					gridlines = append(gridlines, chartdraw.GridLine{
						Style: styleMinor,
						Value: float64(i)})
				}
			}
		}

		if 1 == 0 {
			if interval == timeutil.IntervalMonth && month.MonthContinuousIsYearStart(i) {
				//if interval == timeutil.Month && month.MonthContinuousIsQuarterStart(i) {
				ticks = append(ticks, chartdraw.Tick{Value: float64(i)})
				ticks[len(ticks)-1].Label = timeFormat(iTime)
				if iTime.Month() == 1 {
					gridlines = append(gridlines, chartdraw.GridLine{
						Style: styleMajor,
						Value: float64(i)})
				} else {
					gridlines = append(gridlines, chartdraw.GridLine{
						Style: styleMinor,
						Value: float64(i)})
				}
			} else if interval == timeutil.IntervalQuarter {
				ticks = append(ticks, chartdraw.Tick{
					Value: float64(i)})
				ticks[len(ticks)-1].Label = timeFormat(iTime)
				if iTime.Month() == 1 {
					gridlines = append(gridlines, chartdraw.GridLine{
						Style: styleMajor,
						Value: float64(i)})
				} else {
					gridlines = append(gridlines, chartdraw.GridLine{
						Style: styleMinor,
						Value: float64(i)})
				}
			}
		} else if 1 == 2 { // monthly
			ticks = append(ticks, chartdraw.Tick{
				Value: float64(i)})
			ticks[len(ticks)-1].Label = timeFormat(iTime)
			if iTime.Month() == 1 {
				gridlines = append(gridlines, chartdraw.GridLine{
					Style: styleMajor,
					Value: float64(i)})
			} else {
				gridlines = append(gridlines, chartdraw.GridLine{
					Style: styleMinor,
					Value: float64(i)})
			}
		}
	}

	return ticks, gridlines, nil
}
