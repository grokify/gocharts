package wchart

import (
	"io"
	"os"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/wcharczuk/go-chart"
)

type ChartType interface {
	Render(rp chart.RendererProvider, w io.Writer) error
}

func WritePNG(filename string, thisChart ChartType) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = thisChart.Render(chart.PNG, f)
	err2 := f.Close()
	if err != nil && err2 != nil {
		return errorsutil.Wrap(err, err2.Error())
	} else if err != nil {
		return err
	}
	return err2
}
