// statictimeseriesdata provides tools for adding and formatting
// static time series data for reporting purposes.
package statictimeseries

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/data/table"
	tu "github.com/grokify/gotilla/time/timeutil"
	"github.com/grokify/gotilla/type/stringsutil"
)

type TimeSeriesSimple struct {
	Name        string
	DisplayName string
	Times       []time.Time
}

func NewTimeSeriesSimple(name, displayName string) TimeSeriesSimple {
	return TimeSeriesSimple{
		Name:        name,
		DisplayName: displayName,
		Times:       []time.Time{}}
}

func (tss *TimeSeriesSimple) ToDataSeriesQuarter() DataSeries {
	ds := NewDataSeries()
	ds.SeriesName = tss.Name
	for _, t := range tss.Times {
		ds.AddItem(DataItem{
			SeriesName: tss.Name,
			Time:       tu.QuarterStart(t),
			Value:      int64(1)})
	}
	return ds
}

type TimeSeriesFunnel struct {
	Series map[string]TimeSeriesSimple
	Order  []string
}

func (tsf *TimeSeriesFunnel) Times() []time.Time {
	times := []time.Time{}
	for _, s := range tsf.Series {
		times = append(times, s.Times...)
	}
	return times
}

func (tsf *TimeSeriesFunnel) TimesSorted() []time.Time {
	times := tsf.Times()
	return tu.Sort(times)
}

type DataSeriesSetSimple struct {
	Name   string
	Series map[string]DataSeries
	Times  []time.Time
	Order  []string
}

func NewDataSeriesSetSimple() DataSeriesSetSimple {
	return DataSeriesSetSimple{
		Series: map[string]DataSeries{},
		Times:  []time.Time{},
		Order:  []string{}}
}

func (set *DataSeriesSetSimple) AddItems(items ...DataItem) {
	for _, item := range items {
		set.AddItem(item)
	}
}

func (set *DataSeriesSetSimple) AddItem(item DataItem) {
	item.SeriesName = strings.TrimSpace(item.SeriesName)
	if _, ok := set.Series[item.SeriesName]; !ok {
		set.Series[item.SeriesName] =
			DataSeries{
				SeriesName: item.SeriesName,
				ItemMap:    map[string]DataItem{}}
	}
	series := set.Series[item.SeriesName]
	series.AddItem(item)
	set.Series[item.SeriesName] = series

	set.Times = append(set.Times, item.Time)
}

func (set *DataSeriesSetSimple) AddDataSeries(ds DataSeries) {
	for _, item := range ds.ItemMap {
		if len(item.SeriesName) == 0 {
			item.SeriesName = ds.SeriesName
		}
		set.AddItem(item)
	}
}

func (set *DataSeriesSetSimple) Inflate() {
	if len(set.Times) == 0 {
		set.Times = set.getTimes()
	}
	if len(set.Order) == 0 {
		order := []string{}
		for name := range set.Series {
			order = append(order, name)
		}
		sort.Strings(order)
		set.Order = order
	}
}

func (set *DataSeriesSetSimple) SeriesNames() []string {
	seriesNames := []string{}
	for seriesName := range set.Series {
		seriesNames = append(seriesNames, seriesName)
	}
	sort.Strings(seriesNames)
	return seriesNames
}

func (set *DataSeriesSetSimple) GetItem(seriesName, rfc3339 string) (DataItem, error) {
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

func (set *DataSeriesSetSimple) getTimes() []time.Time {
	times := []time.Time{}
	for _, ds := range set.Series {
		for _, item := range ds.ItemMap {
			times = append(times, item.Time)
		}
	}
	return times
}

func (set *DataSeriesSetSimple) TimeStrings() []string {
	times := []string{}
	for _, ds := range set.Series {
		for rfc3339 := range ds.ItemMap {
			times = append(times, rfc3339)
		}
	}
	return stringsutil.SliceCondenseSpace(times, true, true)
}

func (tsf *TimeSeriesFunnel) DataSeriesSetByQuarter() (DataSeriesSetSimple, error) {
	dss := DataSeriesSetSimple{Order: tsf.Order}
	seriesMapQuarter := map[string]DataSeries{}

	allTimes := []time.Time{}
	for _, s := range tsf.Series {
		allTimes = append(allTimes, s.Times...)
	}

	if len(allTimes) == 0 {
		return dss, errors.New("No times")
	}
	earliest, err := tu.Earliest(allTimes, false)
	if err != nil {
		return dss, err
	}
	latest, err := tu.Latest(allTimes, false)
	if err != nil {
		return dss, err
	}
	earliestQuarter := tu.QuarterStart(earliest)
	latestQuarter := tu.QuarterStart(latest)

	sliceQuarter := tu.QuarterSlice(earliestQuarter, latestQuarter)
	dss.Times = sliceQuarter

	for name, tss := range tsf.Series {
		dataSeries := tss.ToDataSeriesQuarter()
		dataSeries.SeriesName = tss.Name
		for _, q := range sliceQuarter {
			q = q.UTC()
			rfc := q.Format(time.RFC3339)
			if _, ok := dataSeries.ItemMap[rfc]; !ok {
				dataSeries.AddItem(DataItem{
					SeriesName: tss.Name,
					Time:       q,
					Value:      int64(0)})
			}
		}
		seriesMapQuarter[name] = dataSeries
	}
	dss.Series = seriesMapQuarter
	return dss, nil
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
func ReportAxisX(dss DataSeriesSetSimple, cols int, conv func(time.Time) string) []string {
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
func Report(dss DataSeriesSetSimple, cols int, lowFirst bool) []RowInt64 {
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
func DS3ToTable(ds3 DataSeriesSetSimple, fmtTime func(time.Time) string) (table.TableData, error) {
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

func WriteXLSX(filename string, ds3 DataSeriesSetSimple, fmtTime func(time.Time) string) error {
	tbl, err := DS3ToTable(ds3, fmtTime)
	if err != nil {
		return err
	}
	tf := &table.TableFormatter{
		Table:     &tbl,
		Formatter: table.FormatStringAndFloats}
	return table.WriteXLSXFormatted(filename, tf)
}
