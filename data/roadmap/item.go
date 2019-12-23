// roadmap provides data for generating roadmaps
package roadmap

import (
	"fmt"
	"time"

	tu "github.com/grokify/gotilla/time/timeutil"
)

type Item struct {
	MinTime   time.Time
	MaxTime   time.Time
	MinCell   int32 // Inflated by Canvas
	MaxCell   int32 // Inflated by Canvas
	Min       int64 // Inflated by Canvas
	Max       int64 // Inflated by Canvas
	Name      string
	NameShort string
	URL       string
	Color     string
}

func (i *Item) SetMinMaxQuarter(qtrMin, qtrMax int32) error {
	if qtrMax < qtrMin {
		return fmt.Errorf("Max is < min: min [%v] max [%v]", qtrMin, qtrMax)
	}
	err := i.SetMinQuarter(qtrMin)
	if err != nil {
		return err
	}
	return i.SetMaxQuarter(qtrMax)
}

func (i *Item) SetMinQuarter(qtr int32) error {
	qt, err := tu.QuarterInt32StartTime(qtr)
	if err != nil {
		return err
	}
	i.MinTime = qt
	return nil
}

func (i *Item) SetMaxQuarter(qtr int32) error {
	qt, err := tu.QuarterInt32EndTime(qtr)
	if err != nil {
		return err
	}
	i.MaxTime = qt
	return nil
}
