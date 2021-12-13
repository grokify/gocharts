package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/grokify/gocharts/charts/d3/d3bullet"
	"github.com/grokify/gocharts/charts/d3/d3bullet/examples/greenred/templates"
	tu "github.com/grokify/mogo/time/timeutil"
)

func main() {
	bulletsInt64 := []d3bullet.BulletInt64{}

	timeNow := time.Now().UTC()

	bulletBuilderQTD := d3bullet.D3BulletChartBuilder{
		Title:    "Progress QTD",
		YTarget:  int64(50),
		YCurrent: int64(100),
		XStart:   tu.QuarterStart(timeNow).Unix(),
		XCurrent: timeNow.Unix(),
		XEnd:     tu.QuarterEnd(timeNow).Unix(),
	}

	bulletsInt64 = append(bulletsInt64, bulletBuilderQTD.D3Bullet())

	bulletBuilderYTD := d3bullet.D3BulletChartBuilder{
		Title:    "Progress YTD",
		YTarget:  int64(50),
		YCurrent: int64(200),
		XStart:   tu.YearStart(timeNow).Unix(),
		XCurrent: timeNow.Unix(),
		XEnd:     tu.YearEnd(timeNow).Unix(),
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
