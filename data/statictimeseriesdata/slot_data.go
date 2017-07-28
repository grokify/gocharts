package statictimeseriesdata

import (
	"errors"
	"strings"
	"time"
)

// SlotNumber should be epochmilliseconds

type TimeThin struct {
	EpochMs int64
	Time    time.Time
}

func TimeEpochMilliseconds(msec int64) time.Time {
	sec := int64(float64(msec) / 1000.0)
	remainder := msec - (sec * 1000)
	nsec := remainder * 1000000
	return time.Unix(sec, nsec)
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
		for slotNumber, _ := range seriesInfo.SeriesData {
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
		minDt := TimeEpochMilliseconds(minX).UTC()
		maxDt := TimeEpochMilliseconds(maxX).UTC()
		set.MinSlotValue = TimeThin{EpochMs: minX, Time: minDt}
		set.MaxSlotValue = TimeThin{EpochMs: maxX, Time: maxDt}
	}
}

func (set *SlotDataSeriesSet) InflateCanonicalSlots() {
	if strings.ToLower(strings.TrimSpace(set.Interval)) == "quarter" {

	}
	panic("A")
}

type SlotDataSeries struct {
	SeriesName string
	SeriesData map[int64]int64
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
