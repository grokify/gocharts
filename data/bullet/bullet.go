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

func (charts *BulletCharts) AddTimeDT8(dt8 int32) error {
	dt8x := timeutil.DateTime8(dt8)
	dtCur, err := dt8x.Time(time.UTC)
	if err != nil {
		return err
	}
	charts.AddTimeCurrent(dtCur)
	return nil
}

func (charts *BulletCharts) AddTimeCurrent(dtCur time.Time) {
	dtCur = dtCur.UTC()
	dtCurMore := timeutil.NewTimeMore(dtCur, 0)
	charts.TimeCurrent = dtCur
	charts.TimeStart = dtCurMore.QuarterStart()
	charts.TimeEnd = dtCurMore.QuarterEnd()
}

func (charts *BulletCharts) InflateChart(key string) error {
	chart, ok := charts.Charts[key]
	if !ok {
		return fmt.Errorf("chart [%v] not found", key)
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
