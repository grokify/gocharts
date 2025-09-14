package wchart

import (
	"io"
	"os"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/grokify/mogo/errors/errorsutil"
)

type ChartType interface {
	Render(rp chartdraw.RendererProvider, w io.Writer) error
}

func WritePNG(w io.Writer, c ChartType) error {
	return c.Render(chartdraw.PNG, w)
}

func WritePNGFile(filename string, c ChartType) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	// err = c.Render(chartdraw.PNG, f)
	err = WritePNG(f, c)
	err2 := f.Close()
	if err != nil && err2 != nil {
		return errorsutil.Wrap(err, err2.Error())
	} else if err != nil {
		return err
	}
	return err2
}
