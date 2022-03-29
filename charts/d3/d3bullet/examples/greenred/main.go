package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/grokify/mogo/time/timeutil"

	"github.com/grokify/gocharts/v2/charts/d3/d3bullet"
	"github.com/grokify/gocharts/v2/charts/d3/d3bullet/examples/greenred/templates"
)

func main() {
	bulletsInt64 := []d3bullet.BulletInt64{}

	timeNow := time.Now().UTC()

	bulletBuilderQTD := d3bullet.D3BulletChartBuilder{
		Title:    "Progress QTD",
		YTarget:  int64(50),
		YCurrent: int64(100),
		XStart:   timeutil.QuarterStart(timeNow).Unix(),
		XCurrent: timeNow.Unix(),
		XEnd:     timeutil.QuarterEnd(timeNow).Unix(),
	}

	bulletsInt64 = append(bulletsInt64, bulletBuilderQTD.D3Bullet())

	bulletBuilderYTD := d3bullet.D3BulletChartBuilder{
		Title:    "Progress YTD",
		YTarget:  int64(50),
		YCurrent: int64(200),
		XStart:   timeutil.YearStart(timeNow).Unix(),
		XCurrent: timeNow.Unix(),
		XEnd:     timeutil.YearEnd(timeNow).Unix(),
	}

	bulletsInt64 = append(bulletsInt64, bulletBuilderYTD.D3Bullet())

	chartsData := templates.ChartsData{
		DataInt64: d3bullet.DataInt64{
			Bullets: bulletsInt64,
		},
	}

	ioutil.WriteFile("chart.html", []byte(templates.Charts(chartsData)), 0644)

	fmt.Println("DONE")
}
