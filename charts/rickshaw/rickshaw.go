package rickshaw

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/data"
	"github.com/grokify/gocharts/data/statictimeseries/interval"
	"github.com/grokify/gotilla/time/timeutil"
)

// DataInfoJS is the series item to be sent to the
// Rickshaw extensions JS code.
type DataInfoJs struct {
	Color string `json:"color,omitempty"`
	Data  []Item `json:"data"`
	Name  string `json:"name"`
}

type Item struct {
	SeriesName string    `json:"-"`
	Time       time.Time `json:"-"`
	ValueX     int64     `json:"x"`
	ValueY     int64     `json:"y"`
}

// AppMonth is the input value to be convered into
// Rickshaw items
type MonthData struct {
	SeriesName string
	MonthS     string
	YearS      string
	Dt6        int32
	Value      int64
	ValueS     string
}

func (am *MonthData) Inflate() {
	months_en := []string{}
	json.Unmarshal([]byte(timeutil.MonthsEN), &months_en)

	mon := strings.ToLower(am.MonthS)
	monthn := 0
	for i, try := range months_en {
		if mon == strings.ToLower(try) {
			monthn = i + 1
		}
	}
	if monthn < 1 {
		panic("E_DATA_ERROR")
	}
	dt6s := fmt.Sprintf("%v%02d", am.YearS, monthn)

	i, err := strconv.Atoi(dt6s)
	if err != nil {
		panic("E_DATE_CONVERSION_ERROR")
	} else {
		am.Dt6 = int32(i)
	}

	i2, err := strconv.Atoi(am.ValueS)
	if err == nil {
		am.Value = int64(i2)
	}
}

func (am *MonthData) RickshawItem() (Item, error) {
	dt6Time, err := timeutil.TimeForDt6(am.Dt6)
	if err != nil {
		return Item{}, err
	}
	item := Item{
		SeriesName: am.SeriesName,
		Time:       dt6Time,
		ValueY:     am.Value,
		ValueX:     int64(dt6Time.Unix())}
	return item, nil
}

type RickshawData struct {
	SeriesMap map[string]Series
	MinX      int64
	MaxX      int64
}

func NewMonthlyRickshawDataFromSlotSeriesSet(set data.SlotDataSeriesSet) RickshawData {
	rs := NewRickshawData()

	for seriesName, seriesData := range set.SeriesSet {
		for slotValueX, slotValueY := range seriesData.SeriesData {
			rs.AddItem(Item{
				SeriesName: seriesName,
				Time:       time.Unix(slotValueX, 0),
				ValueX:     slotValueX,
				ValueY:     slotValueY})
		}
	}

	rs.Inflate()
	return rs
}

func NewRickshawData() RickshawData {
	return RickshawData{SeriesMap: map[string]Series{}}
}

// Add an item to the report
func (rd *RickshawData) AddItem(item Item) {
	item.SeriesName = strings.TrimSpace(item.SeriesName)
	if len(item.SeriesName) < 1 {
		panic("E_NO_SERIES_NAME")
	}
	series, ok := rd.SeriesMap[item.SeriesName]
	if !ok {
		series = Series{ItemsMapX: map[int64]Item{}}
	}
	series.ItemsMapX[item.ValueX] = item
	rd.SeriesMap[item.SeriesName] = series
}

func (rd *RickshawData) Inflate() {
	minX := int64(0)
	maxX := int64(0)
	first := true
	for seriesName, series := range rd.SeriesMap {
		series.Inflate()
		if first {
			minX = series.MinX
			maxX = series.MaxX
			first = false
		} else {
			if series.MinX < minX {
				minX = series.MinX
			}
			if series.MaxX > maxX {
				maxX = series.MaxX
			}
		}
		rd.SeriesMap[seriesName] = series
	}
	rd.MinX = minX
	rd.MaxX = maxX
}

func (rd *RickshawData) seriesNames() []string {
	seriesNames := []string{}
	first := true
	for seriesName, series := range rd.SeriesMap {
		seriesNames = append(seriesNames, seriesName)
		series.Inflate()
		if first {
			rd.MinX = series.MinX
			rd.MaxX = series.MaxX
			first = false
			continue
		}
		if series.MinX < rd.MinX {
			rd.MinX = series.MinX
		}
		if series.MaxX > rd.MaxX {
			rd.MaxX = series.MaxX
		}
	}
	sort.Strings(seriesNames)
	return seriesNames
}

// Formatted returns formatted information ready for Rickshaw
func (rd *RickshawData) Formatted() RickshawDataFormatted {
	seriesNames := rd.seriesNames()
	sort.Sort(sort.Reverse(sort.StringSlice(seriesNames)))
	formatted := RickshawDataFormatted{
		SeriesNames: seriesNames}

	seriesSet := []Series{}
	for _, seriesName := range seriesNames {
		seriesSet = append(seriesSet, rd.SeriesMap[seriesName])
	}

	times := map[int64]int32{}
	minDt6 := int32(0)
	maxDt6 := int32(0)
	first := true
	for _, series := range seriesSet {
		for _, item := range series.ItemsMapX {
			dt := time.Unix(item.ValueX, 0)
			dt = dt.UTC()

			i, err := strconv.Atoi(dt.Format(timeutil.DT6))
			if err != nil {
				panic("E_STRCONV")
			}
			dt6 := int32(i)
			if first {
				minDt6 = dt6
				maxDt6 = dt6
				first = false
			} else {
				if dt6 < minDt6 {
					minDt6 = dt6
				}
				if dt6 > maxDt6 {
					maxDt6 = dt6
				}
			}
			times[item.ValueY] = dt6
		}
	}

	timeutil.ParseDt6(minDt6)

	dt6Axis := timeutil.Dt6MinMaxSlice(minDt6, maxDt6)

	seriesSetInflatedSorted := [][]Item{}
	for _, thinSeries := range seriesSet {
		fullSeries := []Item{}
		for _, dt6 := range dt6Axis {
			dt, err := timeutil.TimeForDt6(dt6)
			if err != nil {
				panic("DT6_PARSE_ERROR")
			}
			dt6Epoch := dt.Unix()
			if item, ok := thinSeries.ItemsMapX[dt6Epoch]; ok {
				fullSeries = append(fullSeries, item)
			} else {
				fullSeries = append(fullSeries, Item{
					Time:   time.Unix(dt6Epoch, 0),
					ValueX: dt6Epoch, ValueY: 0})
			}
		}
		seriesSetInflatedSorted = append(seriesSetInflatedSorted, fullSeries)
	}

	formatted.SeriesData = seriesSetInflatedSorted

	formattedSeries := []DataInfoJs{}
	for i, seriesName := range seriesNames {
		data := DataInfoJs{
			Name: seriesName,
			Data: seriesSetInflatedSorted[i]}
		formattedSeries = append(formattedSeries, data)
	}

	formatted.FormattedData = formattedSeries

	return formatted
}

type Series struct {
	ItemsMapX map[int64]Item
	MinX      int64
	MaxX      int64
}

func (s *Series) Inflate() {
	first := true
	s.MinX = 0
	s.MaxX = 0
	for _, item := range s.ItemsMapX {
		if first {
			s.MinX = item.ValueX
			s.MaxX = item.ValueX
			first = false
			continue
		}
		if item.ValueX < s.MinX {
			s.MinX = item.ValueX
		}
		if item.ValueX > s.MaxX {
			s.MaxX = item.ValueX
		}
	}
}

type RickshawDataFormatted struct {
	SeriesNames   []string
	SeriesData    [][]Item
	FormattedData []DataInfoJs
}

func NewRickshawDataFormattedFromDateHistogram(timeset interval.DataSeriesSet) RickshawDataFormatted {
	formatted := RickshawDataFormatted{}
	formatted.SeriesNames = timeset.SeriesNamesSorted()
	formatted.SeriesData = [][]Item{}
	formatted.FormattedData = []DataInfoJs{}

	for _, seriesName := range formatted.SeriesNames {
		seriesData, ok := timeset.OutputSeriesMap[seriesName]
		if !ok {
			panic("series not found")
		}
		formattedData := []Item{}
		dataInfoJs := DataInfoJs{Name: seriesName, Data: []Item{}}
		for _, sourceItem := range seriesData.SortedItems() {
			rsItem := Item{
				SeriesName: sourceItem.SeriesName,
				Time:       sourceItem.Time,
				ValueX:     sourceItem.Time.Unix(),
				ValueY:     sourceItem.Value,
			}
			formattedData = append(formattedData, rsItem)
			dataInfoJs.Data = append(dataInfoJs.Data, rsItem)
		}
		formatted.SeriesData = append(formatted.SeriesData, formattedData)
		formatted.FormattedData = append(formatted.FormattedData, dataInfoJs)
	}
	return formatted
}
