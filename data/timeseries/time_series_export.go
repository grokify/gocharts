package timeseries

import (
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/mogo/os/osutil"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/time/timeutil"

	"github.com/grokify/gocharts/v2/data/table"
)

func TimeFormatRFC3339(dt time.Time) string {
	return dt.Format(time.RFC3339)
}

func TimeFormatNiceMonth(dt time.Time) string {
	return dt.Format("Jan '06")
}

func TimeFormatNiceQuarter(dt time.Time) string {
	return timeutil.FormatQuarterYYQ(dt)
}

// Table generates a `table.Table` given a `TimeSeries`.
func (ts *TimeSeries) Table(tableName, dateColumnName, countColumnName string, dtFmt func(dt time.Time) string) table.Table {
	// previously only took dateColumnName as a parameter.
	if len(strings.TrimSpace(dateColumnName)) == 0 {
		switch ts.Interval {
		case timeutil.Month:
			dateColumnName = "Month"
		case timeutil.Quarter:
			dateColumnName = "Quarter"
		default:
			dateColumnName = "Date"
		}
	}
	if len(strings.TrimSpace(countColumnName)) == 0 {
		countColumnName = "Count"
	}
	tbl := table.NewTable(tableName)
	tbl.Columns = []string{dateColumnName, countColumnName}
	tbl.FormatMap = map[int]string{}
	if ts.IsFloat {
		tbl.FormatMap[1] = table.FormatFloat
	} else {
		tbl.FormatMap[1] = table.FormatInt
	}
	if dtFmt == nil {
		dtFmt = func(dt time.Time) string {
			return dt.Format(time.RFC3339)
		}
		tbl.FormatMap[0] = table.FormatTime
	}
	itemsSorted := ts.ItemsSorted()
	for _, item := range itemsSorted {
		row := []string{
			dtFmt(item.Time)}
		if ts.IsFloat {
			row = append(row, strconv.FormatFloat(item.ValueFloat, 'f', -1, 64))
		} else {
			row = append(row, strconv.Itoa(int(item.Value)))
		}
		tbl.Rows = append(tbl.Rows, row)
	}
	return tbl
}

type TableMonthXOXOpts struct {
	AddMOMGrowth           bool
	MOMGrowthPct           float64
	MOMBaseMonth           time.Time
	MOMTargetName          string
	MOMPerformanceName     string
	momBaseMonthContinuous uint64
	momBaseTimeItem        TimeItem
	momBaseTimeItemExists  bool
}

func (ts *TimeSeries) TableMonthXOX(timeFmtColName, seriesName, valuesName, yoyName, qoqName, momName string, opts *TableMonthXOXOpts) table.Table {
	if len(strings.TrimSpace(seriesName)) == 0 {
		seriesName = "Series"
	}
	if len(strings.TrimSpace(valuesName)) == 0 {
		valuesName = "Values"
	}
	if len(strings.TrimSpace(yoyName)) == 0 {
		yoyName = "YoY"
	}
	if len(strings.TrimSpace(qoqName)) == 0 {
		qoqName = "QoQ"
	}
	if len(strings.TrimSpace(momName)) == 0 {
		momName = "MoM"
	}
	if opts == nil {
		opts = &TableMonthXOXOpts{}
	}
	if opts.AddMOMGrowth {
		minDt, _ := ts.MinMaxTimes()
		if opts.MOMBaseMonth.Before(minDt) {
			opts.MOMBaseMonth = minDt
		}
		opts.MOMBaseMonth = timeutil.NewTimeMore(opts.MOMBaseMonth.UTC(), 0).MonthStart()
		opts.momBaseMonthContinuous = month.TimeToMonthContinuous(opts.MOMBaseMonth)
		momBaseTimeItem, err := ts.Get(opts.MOMBaseMonth)
		if err != nil {
			opts.momBaseTimeItemExists = false
		} else {
			opts.momBaseTimeItemExists = true
			opts.momBaseTimeItem = momBaseTimeItem
		}
		if len(strings.TrimSpace(opts.MOMTargetName)) == 0 {
			opts.MOMTargetName = momName + " Target Value"
		}
		if len(strings.TrimSpace(opts.MOMPerformanceName)) == 0 {
			opts.MOMPerformanceName = momName + " Performance"
		}
	}
	tsm := ts.ToMonth(true)
	tbl := table.NewTable("")
	cols := []string{seriesName}
	times := tsm.Times(true)
	for _, dt := range times {
		cols = append(cols, dt.Format(timeFmtColName))
	}
	tbl.Columns = cols
	tbl.FormatMap = map[int]string{
		-1: table.FormatFloat,
		0:  table.FormatString}

	yoy := tsm.TimeSeriesMonthYOY()
	qoq := tsm.TimeSeriesMonthQOQ()
	mom := tsm.TimeSeriesMonthMOM()

	valData := []string{valuesName}
	yoyData := []string{yoyName}
	qoqData := []string{qoqName}
	momData := []string{momName}
	momGrowthTargets := []string{opts.MOMTargetName}
	momGrowthPerform := []string{opts.MOMPerformanceName}

	for _, dt := range times {
		tiVal, err := tsm.Get(dt)
		if err != nil {
			panic("internal time not found")
		}
		valData = append(valData, strconvutil.FormatFloat64Simple(tiVal.Float64()))
		tiYOY, err := yoy.Get(dt)
		if err != nil {
			yoyData = append(yoyData, "0")
		} else {
			yoyData = append(yoyData, strconvutil.FormatFloat64Simple(tiYOY.Float64()))
		}
		tiQOQ, err := qoq.Get(dt)
		if err != nil {
			qoqData = append(qoqData, "0")
		} else {
			qoqData = append(qoqData, strconvutil.FormatFloat64Simple(tiQOQ.Float64()))
		}
		tiMOM, err := mom.Get(dt)
		if err != nil {
			momData = append(momData, "0")
		} else {
			momData = append(momData, strconvutil.FormatFloat64Simple(tiMOM.Float64()))
		}
		if opts.AddMOMGrowth {
			if dt.After(opts.MOMBaseMonth) && opts.momBaseTimeItemExists {
				diffMonths := int(month.TimeToMonthContinuous(dt) - opts.momBaseMonthContinuous)
				targetValue := opts.momBaseTimeItem.Float64() * math.Pow(1+opts.MOMGrowthPct, float64(diffMonths))
				momGrowthTargets = append(momGrowthTargets, strconvutil.FormatFloat64Simple(targetValue))
				actualValue := tiVal.Float64()
				diff := 0.0
				if targetValue != 0 {
					diff = (actualValue - targetValue) / targetValue
				}
				momGrowthPerform = append(momGrowthPerform, strconvutil.FormatFloat64Simple(diff))
			} else {
				momGrowthTargets = append(momGrowthTargets, "0")
				momGrowthPerform = append(momGrowthPerform, "0")
			}
		}
	}
	tbl.Rows = [][]string{valData, yoyData, qoqData, momData}
	if opts.AddMOMGrowth {
		tbl.Rows = append(tbl.Rows, momGrowthTargets, momGrowthPerform)
	}
	return tbl
}

func (ts *TimeSeries) TableYearYOY(seriesName, valuesName, yoyName string) table.Table {
	if len(strings.TrimSpace(seriesName)) == 0 {
		seriesName = "Series"
	}
	if len(strings.TrimSpace(valuesName)) == 0 {
		valuesName = "Values"
	}
	if len(strings.TrimSpace(yoyName)) == 0 {
		yoyName = "YoY"
	}
	tbl := table.NewTable(ts.SeriesName)
	cols := []string{seriesName}
	times := ts.Times(true)
	for _, dt := range times {
		cols = append(cols, dt.Format("2006"))
	}
	tbl.Columns = cols
	tbl.FormatMap = map[int]string{
		-1: table.FormatFloat,
		0:  table.FormatString}

	yoy := ts.TimeSeriesMonthYOY()
	valData := []string{valuesName}
	yoyData := []string{yoyName}

	for _, dt := range times {
		tiVal, err := ts.Get(dt)
		if err != nil {
			tiVal = TimeItem{
				Time:    dt,
				IsFloat: ts.IsFloat}
		}
		valData = append(valData, strconvutil.FormatFloat64Simple(tiVal.Float64()))
		tiYOY, err := yoy.Get(dt)
		if err != nil {
			yoyData = append(yoyData, "0")
		} else {
			yoyData = append(yoyData, strconvutil.FormatFloat64Simple(tiYOY.Float64()))
		}
	}
	tbl.Rows = [][]string{valData, yoyData}
	return tbl
}

// WriteJSON writes the data to a JSON file. To write a minimized JSON
// file use an empty string for `prefix` and `indent`.
func (ts *TimeSeries) WriteJSON(filename string, perm os.FileMode, prefix, indent string) error {
	return osutil.WriteFileJSON(filename, ts, perm, prefix, indent)
}

// WriteXLSX writes an XSLX file given a `TimeSeries`
func (ts *TimeSeries) WriteXLSX(filename string, sheetName, dateColumnName, countColumnName string) error {
	tbl := ts.Table(sheetName, dateColumnName, countColumnName, nil)
	return table.WriteXLSX(filename, []*table.Table{&tbl})
}
