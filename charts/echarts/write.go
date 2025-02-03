package echarts

import (
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/components"
)

func WriteSimple(filename string, charts []components.Charter) error {
	page := components.NewPage()
	page.AddCharts(charts...)
	if f, err := os.Create(filename); err != nil {
		return err
	} else {
		return page.Render(io.MultiWriter(f))
	}
}
