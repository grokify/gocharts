package rickshaw

import (
	"errors"
	"strings"
)

type SlotData struct {
	SeriesName string
	SlotValue  int64
	SlotNumber int64
}

type SlotDataSeriesSet struct {
	SeriesSet map[string]SlotDataSeries
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
