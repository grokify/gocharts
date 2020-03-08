package wchart

import (
	"os"

	"github.com/pkg/errors"
	"github.com/wcharczuk/go-chart"
)

func WritePng(filename string, thisChart chart.Chart) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = thisChart.Render(chart.PNG, f)
	err2 := f.Close()
	if err != nil && err2 != nil {
		return errors.Wrap(err, err2.Error())
	} else if err != nil {
		return err
	}
	return err2
}
