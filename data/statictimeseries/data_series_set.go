// statictimeseriesdata provides tools for adding and formatting
// static time series data for reporting purposes.
package statictimeseries

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gotilla/sort/sortutil"
	"github.com/grokify/gotilla/time/month"
	"github.com/grokify/gotilla/time/timeutil"
	tu "github.com/grokify/gotilla/time/timeutil"
	"github.com/grokify/gotilla/type/stringsutil"
	"github.com/pkg/errors"
)

type DataSeriesSet struct {
	Name     string
	Series   map[string]DataSeries
	Times    []time.Time
	Order    []string
	IsFloat  bool
	Interval timeutil.Interval
}

func NewDataSeriesSet() DataSeriesSet {
	return DataSeriesSet{
		Series: map[string]DataSeries{},
		Times:  []time.Time{},
		Order:  []string{}}
}

func (set *DataSeriesSet) AddItems(items ...DataItem) {
	for _, item := range items {
		set.AddItem(item)
	}
}

func (set *DataSeriesSet) AddItem(item DataItem) {
	item.SeriesName = strings.TrimSpace(item.SeriesName)
	if _, ok := set.Series[item.SeriesName]; !ok {
		set.Series[item.SeriesName] =
			DataSeries{
				SeriesName: item.SeriesName,
				ItemMap:    map[string]DataItem{},
				IsFloat:    item.IsFloat,
				Interval:   set.Interval}
	}
	series := set.Series[item.SeriesName]
	series.AddItem(item)
	set.Series[item.SeriesName] = series

	set.Times = append(set.Times, item.Time)
}

func (set *DataSeriesSet) AddDataSeries(dataSeries ...DataSeries) error {
	for _, ds := range dataSeries {
		ds.SeriesName = strings.TrimSpace(ds.SeriesName)
		if len(ds.SeriesName) == 0 {
			return errors.New("E_DataSeriesSet.AddDataSeries_NO_DataSeries.SeriesName")
		}
		for _, item := range ds.ItemMap {
			if len(item.SeriesName) == 0 || item.SeriesName != ds.SeriesName {
				item.SeriesName = ds.SeriesName
			}
			set.AddItem(item)
		}
	}
	return nil
}

func (set *DataSeriesSet) Inflate() {
	set.Times = set.GetTimeSlice(true)
	if len(set.Order) > 0 {
		set.Order = stringsutil.SliceCondenseSpace(set.Order, true, false)
	} else {
		order := []string{}
		for name := range set.Series {
			order = append(order, name)
		}
		sort.Strings(order)
		set.Order = order
	}
}

func (set *DataSeriesSet) SeriesNames() []string {
	seriesNames := []string{}
	for seriesName := range set.Series {
		seriesNames = append(seriesNames, seriesName)
	}
	sort.Strings(seriesNames)
	return seriesNames
}

func (set *DataSeriesSet) GetSeriesByIndex(index int) (DataSeries, error) {
	if len(set.Order) == 0 && len(set.Series) > 0 {
		set.Inflate()
	}
	if index < len(set.Order) {
		name := set.Order[index]
		if ds, ok := set.Series[name]; ok {
			return ds, nil
		}
	}
	return DataSeries{}, fmt.Errorf("E_CANNOT_FIND_INDEX_[%d]_SET_COUNT_[%d]", index, len(set.Order))
}

func (set *DataSeriesSet) GetItem(seriesName, rfc3339 string) (DataItem, error) {
	di := DataItem{}
	dss, ok := set.Series[seriesName]
	if !ok {
		return di, fmt.Errorf("SeriesName not found [%s]", seriesName)
	}
	item, ok := dss.ItemMap[rfc3339]
	if !ok {
		return di, fmt.Errorf("SeriesName found [%s] Time not found [%s]", seriesName, rfc3339)
	}
	return item, nil
}

func (set *DataSeriesSet) GetTimeSlice(sortAsc bool) sortutil.TimeSlice {
	times := []time.Time{}
	for _, ds := range set.Series {
		for _, item := range ds.ItemMap {
			times = append(times, item.Time)
		}
	}
	return month.TimeSeriesMonth(sortAsc, times...)
}

func (set *DataSeriesSet) TimeStrings() []string {
	times := []string{}
	for _, ds := range set.Series {
		for rfc3339 := range ds.ItemMap {
			times = append(times, rfc3339)
		}
	}
	return stringsutil.SliceCondenseSpace(times, true, true)
}

func (set *DataSeriesSet) MinMaxTimes() (time.Time, time.Time) {
	values := sortutil.TimeSlice{}
	for _, ds := range set.Series {
		min, max := ds.MinMaxTimes()
		values = append(values, min, max)
	}
	sort.Sort(values)
	return values[0], values[len(values)-1]
}

func (set *DataSeriesSet) MinMaxValues() (int64, int64) {
	values := sortutil.Int64Slice{}
	for _, ds := range set.Series {
		min, max := ds.MinMaxValues()
		values = append(values, min, max)
	}
	sort.Sort(values)
	return values[0], values[len(values)-1]
}

func (set *DataSeriesSet) MinMaxValuesFloat64() (float64, float64) {
	values := sort.Float64Slice{}
	for _, ds := range set.Series {
		min, max := ds.MinMaxValuesFloat64()
		values = append(values, min, max)
	}
	sort.Sort(values)
	return values[0], values[len(values)-1]
}

func (set *DataSeriesSet) ToMonth() DataSeriesSet {
	newDss := DataSeriesSet{
		Name:     set.Name,
		Series:   map[string]DataSeries{},
		Times:    set.Times,
		Interval: timeutil.Month,
		Order:    set.Order}
	for name, ds := range set.Series {
		newDss.Series[name] = ds.ToMonth()
	}
	newDss.Times = newDss.GetTimeSlice(true)
	return newDss
}

func (set *DataSeriesSet) ToMonthCumulative(popLast bool) (DataSeriesSet, error) {
	newDss := DataSeriesSet{
		Name:     set.Name,
		Series:   map[string]DataSeries{},
		Times:    set.Times,
		Interval: timeutil.Month,
		Order:    set.Order}
	for name, ds := range set.Series {
		newDs, err := ds.ToMonthCumulative(newDss.Times...)
		if err != nil {
			return newDss, err
		}
		newDss.Series[name] = newDs
	}
	if popLast {
		newDss.PopLast()
	}
	newDss.Times = newDss.GetTimeSlice(true)
	return newDss, nil
}

func (set *DataSeriesSet) PopLast() {
	times := set.GetTimeSlice(true)
	if len(times) == 0 {
		return
	}
	last := times[len(times)-1]
	set.DeleteItemByTime(last)
}

func (set *DataSeriesSet) DeleteItemByTime(dt time.Time) {
	for id, ds := range set.Series {
		ds.DeleteByTime(dt)
		set.Series[id] = ds
	}
}

type RowInt64 struct {
	Name         string
	DisplayName  string
	HavePlusOne  bool
	ValuePlusOne int64
	Values       []int64
}

func (row *RowInt64) Flatten(conv func(v int64) string) []string {
	strs := []string{row.Name}
	for _, v := range row.Values {
		strs = append(strs, conv(v))
	}
	return strs
}

type RowFloat64 struct {
	Name   string
	Values []float64
}

func (row *RowFloat64) Flatten(conv func(v float64) string) []string {
	strs := []string{row.Name}
	for _, v := range row.Values {
		strs = append(strs, conv(v))
	}
	return strs
}

// ReportAxisX generates data for use with `C3Chart.C3Axis.C3AxisX.Categories`.
func ReportAxisX(dss DataSeriesSet, cols int, conv func(time.Time) string) []string {
	var times tu.TimeSlice
	if cols < len(dss.Times) {
		min := len(dss.Times) - cols
		times = dss.Times[min:]
	} else { // cols >= len(dss.Times)
		times = dss.Times
	}
	cats := []string{}
	for _, t := range times {
		cats = append(cats, conv(t))
	}
	return cats
}

// Report generates data for use with `C3Chart.C3ChartData.Columns`.
func Report(dss DataSeriesSet, cols int, lowFirst bool) []RowInt64 {
	rows := []RowInt64{}
	var times tu.TimeSlice
	var timePlus1 time.Time
	havePlus1 := false
	if cols < len(dss.Times) {
		min := len(dss.Times) - cols
		prev := min - 1
		times = dss.Times[min:]
		timePlus1 = dss.Times[prev]
		havePlus1 = true
	} else { // cols >= len(dss.Times)
		times = dss.Times
		if cols > len(dss.Times) {
			timePlus1 = dss.Times[len(dss.Times)-cols-1]
			havePlus1 = true
		}
	}
	timePlus1Rfc := timePlus1.UTC().Format(time.RFC3339)
	if !lowFirst {
		times = sort.Reverse(times).(tu.TimeSlice)
	}
	for _, seriesName := range dss.Order {
		row := RowInt64{
			Name:        seriesName + " Count",
			HavePlusOne: havePlus1,
		}
		if ds, ok := dss.Series[seriesName]; !ok {
			for i := 0; i < cols; i++ {
				row.Values = append(row.Values, 0)
			}
			if havePlus1 {
				row.ValuePlusOne = 0
			}
		} else {
			for _, t := range times {
				rfc := t.UTC().Format(time.RFC3339)
				if item, ok := ds.ItemMap[rfc]; ok {
					row.Values = append(row.Values, item.Value)
				} else {
					row.Values = append(row.Values, 0)
				}
			}
			if havePlus1 {
				if item, ok := ds.ItemMap[timePlus1Rfc]; ok {
					row.ValuePlusOne = item.Value
				} else {
					row.ValuePlusOne = 0
				}
			}
		}
		rows = append(rows, row)
	}
	return rows
}

func ReportFunnelPct(rows []RowInt64) []RowFloat64 {
	pcts := []RowFloat64{}
	if len(rows) < 2 {
		return pcts
	}
	for i := 0; i < len(rows)-1; i++ {
		r := RowFloat64{Name: fmt.Sprintf("Success Pct #%v", i)}
		j := i + 1
		for k := 0; k < len(rows[0].Values); k++ {
			v1 := rows[i].Values[k]
			v2 := rows[j].Values[k]
			pct := float64(v2) / float64(v1)
			r.Values = append(r.Values, pct)
		}
		pcts = append(pcts, r)
	}
	return pcts
}

func ReportGrowthPct(rows []RowInt64) []RowFloat64 {
	grows := []RowFloat64{}
	if len(rows) == 0 {
		return grows
	}
	for i := 0; i < len(rows); i++ {
		r := rows[i]
		grow := RowFloat64{Name: fmt.Sprintf("%v XoX", r.Name)}
		if r.HavePlusOne {
			pct := float64(r.Values[0]) / float64(r.ValuePlusOne)
			grow.Values = append(grow.Values, pct)
		}
		for j := 0; j < len(r.Values)-1; j++ {
			k := j + 1
			pct := float64(r.Values[k]) / float64(r.Values[j])
			grow.Values = append(grow.Values, pct)
		}
		grows = append(grows, grow)
	}
	return grows
}

// DS3ToTable returns a `DataSeriesSetSimple` as a
// `table.TableData`.
func DS3ToTable(ds3 DataSeriesSet, fmtTime func(time.Time) string) (table.TableData, error) {
	tbl := table.NewTableData()
	seriesNames := ds3.SeriesNames()
	tbl.Columns = []string{"Time"}
	tbl.Columns = append(tbl.Columns, seriesNames...)
	timeStrings := ds3.TimeStrings()
	for _, rfc3339 := range timeStrings {
		dt, err := time.Parse(time.RFC3339, rfc3339)
		if err != nil {
			return tbl, err
		}
		line := []string{fmtTime(dt)}
		for _, seriesName := range seriesNames {
			item, err := ds3.GetItem(seriesName, rfc3339)
			if err == nil {
				if item.IsFloat {
					line = append(line, fmt.Sprintf("%.10f", item.ValueFloat))
				} else {
					line = append(line, strconv.Itoa(int(item.Value)))
				}
			} else {
				line = append(line, "0")
			}
		}
		tbl.Records = append(tbl.Records, line)
	}
	return tbl, nil
}

func WriteXLSX(filename string, ds3 DataSeriesSet, fmtTime func(time.Time) string) error {
	tbl, err := DS3ToTable(ds3, fmtTime)
	if err != nil {
		return err
	}
	tbl.FormatFunc = table.FormatStringAndFloats
	return table.WriteXLSX(filename, &tbl)
}
