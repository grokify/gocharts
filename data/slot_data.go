package data

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/grokify/simplego/time/timeutil"
	"github.com/grokify/simplego/type/maputil"
)

type TimeThin struct {
	EpochMs int64
	Time    time.Time
}

type SlotData struct {
	SeriesName string
	SlotValue  int64
	SlotNumber int64
}

type SlotDataSeriesSet struct {
	SeriesSet      map[string]SlotDataSeries
	MinSlotValue   TimeThin
	MaxSlotValue   TimeThin
	Interval       string
	CanonicalSlots []TimeThin
}

func NewSlotDataSeriesSet() SlotDataSeriesSet {
	return SlotDataSeriesSet{SeriesSet: map[string]SlotDataSeries{}}
}

func (set *SlotDataSeriesSet) AddData(data SlotData) error {
	seriesName := strings.TrimSpace(data.SeriesName)
	if len(seriesName) == 0 {
		return errors.New("No Series Name")
	}
	set.CreateSeriesIfNotExists(seriesName)
	series := set.SeriesSet[seriesName]
	series.AddSlotData(data)
	set.SeriesSet[seriesName] = series
	return nil
}

func (set *SlotDataSeriesSet) CreateSeriesIfNotExists(seriesName string) {
	seriesName = strings.TrimSpace(seriesName)
	if len(seriesName) < 1 {
		seriesName = "_none"
	}
	if _, ok := set.SeriesSet[seriesName]; !ok {
		set.SeriesSet[seriesName] = SlotDataSeries{
			SeriesName: seriesName,
			SeriesData: map[int64]int64{}}
	}
}

func (set *SlotDataSeriesSet) Inflate() {
	set.InflateMinMaxX()
}

func (set *SlotDataSeriesSet) InflateMinMaxX() {
	minX := int64(-1)
	maxX := int64(-1)
	haveSetMinMaxX := false
	for _, seriesInfo := range set.SeriesSet {
		for slotNumber := range seriesInfo.SeriesData {
			if !haveSetMinMaxX {
				minX = slotNumber
				maxX = slotNumber
				haveSetMinMaxX = true
				continue
			}
			if slotNumber < minX {
				minX = slotNumber
			}
			if slotNumber > maxX {
				maxX = slotNumber
			}
		}
	}
	if haveSetMinMaxX {
		minDt := timeutil.UnixMillis(minX).UTC()
		maxDt := timeutil.UnixMillis(maxX).UTC()
		set.MinSlotValue = TimeThin{EpochMs: minX, Time: minDt}
		set.MaxSlotValue = TimeThin{EpochMs: maxX, Time: maxDt}
	}
}

type SlotDataSeries struct {
	SeriesName string
	SeriesData map[int64]int64
}

func (series *SlotDataSeries) Add(name string, x, y int64) {
	if _, ok := series.SeriesData[x]; !ok {
		series.SeriesData[x] = int64(0)
	}
	series.SeriesData[x] += y
}

func (series *SlotDataSeries) AddSlotData(slot SlotData) {
	series.CreateSlotNumberIfNotExists(slot.SlotNumber)
	series.SeriesData[slot.SlotNumber] += slot.SlotValue
}

func (series *SlotDataSeries) CreateSlotNumberIfNotExists(slotNumber int64) {
	if _, ok := series.SeriesData[slotNumber]; !ok {
		series.SeriesData[slotNumber] = 0
	}
}

func (series *SlotDataSeries) DataKeysSorted() []int64 {
	mii := maputil.MapInt64Int64(series.SeriesData)
	return mii.KeysSorted()
}

func (series *SlotDataSeries) DataValuesSortedByKeys() []int64 {
	mii := maputil.MapInt64Int64(series.SeriesData)
	return mii.ValuesSortedByKeys()
}

// SlotDataSeriesSetSimple is useful for C3 Bar Charts
type SlotDataSeriesSetSimple struct {
	SeriesSet map[string]SlotDataSeries
}

func NewSlotDataSeriesSetSimple() SlotDataSeriesSetSimple {
	return SlotDataSeriesSetSimple{SeriesSet: map[string]SlotDataSeries{}}
}

func (set SlotDataSeriesSetSimple) Add(dataSeriesName string, x, y int64) {
	slotDataSeries, ok := set.SeriesSet[dataSeriesName]
	if !ok {
		slotDataSeries = SlotDataSeries{
			SeriesName: dataSeriesName,
			SeriesData: map[int64]int64{}}
	}

	slotDataSeries.AddSlotData(SlotData{
		SeriesName: dataSeriesName,
		SlotValue:  y,
		SlotNumber: x})

	set.SeriesSet[dataSeriesName] = slotDataSeries
}

func (set SlotDataSeriesSetSimple) KeysSorted() []string {
	keys := []string{}
	for k := range set.SeriesSet {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (set SlotDataSeriesSetSimple) MinMaxX() (int64, int64) {
	minX := int64(0)
	maxX := int64(0)
	init := false
	for _, slotDataSeries := range set.SeriesSet {
		for x := range slotDataSeries.SeriesData {
			if !init {
				minX = x
				maxX = x
				init = true
				continue
			}
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
		}
	}
	return minX, maxX
}

type SlotDataSeriesString struct {
	SeriesName string
	SeriesData map[string]int64
}

func NewSlotDataSeriesString(name string) SlotDataSeriesString {
	return SlotDataSeriesString{
		SeriesName: name,
		SeriesData: map[string]int64{}}
}

func (series *SlotDataSeriesString) Add(desc string, y int64) {
	if _, ok := series.SeriesData[desc]; !ok {
		series.SeriesData[desc] = int64(0)
	}
	series.SeriesData[desc] += y
}

type SlotDataSeriesStringSetSimple struct {
	SeriesSet map[string]SlotDataSeriesString
}

func NewSlotDataSeriesStringSetSimple() SlotDataSeriesStringSetSimple {
	return SlotDataSeriesStringSetSimple{SeriesSet: map[string]SlotDataSeriesString{}}
}

func (set *SlotDataSeriesStringSetSimple) Add(seriesName, bucketDesc string, y int64) {
	slotDataSeriesString, ok := set.SeriesSet[seriesName]
	if !ok {
		slotDataSeriesString = NewSlotDataSeriesString(seriesName)
	}
	slotDataSeriesString.Add(bucketDesc, y)

	set.SeriesSet[seriesName] = slotDataSeriesString
}

func AggregateSlotDataSeriesString(pointData SlotDataSeriesSetSimple, bucketSize int64) SlotDataSeriesStringSetSimple {
	agg := NewSlotDataSeriesStringSetSimple()

	for seriesName, slotDataSeries := range pointData.SeriesSet {
		for x, y := range slotDataSeries.SeriesData {
			bucketIndex := BucketIndex(bucketSize, x)
			bucketDesc := fmt.Sprintf("%v", bucketIndex+int64(1))
			/*if 1 == 0 {
				bucketMin, bucketMax := BucketMinMax(bucketSize, bucketIndex)
				bucketDesc = fmt.Sprintf("%v - %v", bucketMin, bucketMax)
			}*/
			agg.Add(seriesName, bucketDesc, y)
		}
	}

	return agg
}

func AggregateSlotDataSeries(pointData SlotDataSeriesSetSimple, bucketSize int64) SlotDataSeriesSetSimple {
	agg := NewSlotDataSeriesSetSimple()

	for seriesName, slotDataSeries := range pointData.SeriesSet {
		for x, y := range slotDataSeries.SeriesData {
			bucketIndex := BucketIndex(bucketSize, x)
			bucketNumber := bucketIndex + int64(1)
			agg.Add(seriesName, bucketNumber, y)
		}
	}

	return agg
}
