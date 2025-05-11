// roadmap provides data for generating roadmaps
package roadmap

import (
	"fmt"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

type Item struct {
	MinTime            time.Time
	MaxTime            time.Time
	MinCell            int32 // Inflated by Canvas
	MaxCell            int32 // Inflated by Canvas
	Min                int64 // Inflated by Canvas
	Max                int64 // Inflated by Canvas
	Name               string
	NameShort          string
	URL                string
	ForegroundColorHex string
	BackgroundColorHex string
}

func (i *Item) SetMinMaxQuarter(qtrMin, qtrMax int) error {
	if qtrMax < qtrMin {
		return fmt.Errorf("max is < min: min [%v] max [%v]", qtrMin, qtrMax)
	}
	err := i.SetMinQuarter(qtrMin)
	if err != nil {
		return err
	}
	return i.SetMaxQuarter(qtrMax)
}

func (i *Item) SetMinQuarter(yyyyq int) error {
	qt, err := timeutil.YearQuarterStartTime(yyyyq)
	if err != nil {
		return err
	}
	i.MinTime = qt
	return nil
}

func (i *Item) SetMaxQuarter(yyyyq int) error {
	qt, err := timeutil.YearQuarterEndTime(yyyyq)
	if err != nil {
		return err
	}
	i.MaxTime = qt
	return nil
}
