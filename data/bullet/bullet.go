package bullet

import (
	"fmt"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

type BulletCharts struct {
	TimeStart   time.Time
	TimeCurrent time.Time
	TimeEnd     time.Time
	Charts      map[string]BulletChart
}

func NewBulletCharts() BulletCharts {
	return BulletCharts{Charts: map[string]BulletChart{}}
}

func (charts *BulletCharts) AddTimeDt8(dt8 int32) error {
	dtCur, err := timeutil.TimeForDt8(dt8)
	if err != nil {
		return err
	}
	return charts.AddTimeCurrent(dtCur)
}

func (charts *BulletCharts) AddTimeCurrent(dtCur time.Time) error {
	dtCur = dtCur.UTC()
	charts.TimeCurrent = dtCur
	charts.TimeStart = timeutil.QuarterStart(dtCur)
	charts.TimeEnd = timeutil.QuarterEnd(dtCur)
	return nil
}

func (charts *BulletCharts) InflateChart(key string) error {
	chart, ok := charts.Charts[key]
	if !ok {
		return fmt.Errorf("chart [%v] not found.", key)
	}

	chart.ProjectionData.CreateProjection(
		charts.TimeStart.UTC().Unix(),
		charts.TimeCurrent.UTC().Unix(),
		charts.TimeEnd.UTC().Unix(),
	)

	charts.Charts[key] = chart
	return nil
}

type BulletChart struct {
	Title          string
	Subtitle       string
	ProjectionData ProjectionDataInt64
}

func (bc *BulletChart) Inflate(start, current, end time.Time) {
	bc.ProjectionData.CreateProjection(
		start.UTC().Unix(),
		current.UTC().Unix(),
		end.UTC().Unix(),
	)
}
