package histogram

import (
	"fmt"
	"io"
	"strconv"

	"github.com/grokify/mogo/type/maputil"
)

func (hsets *HistogramSets) Markdown(w io.Writer, headingPrefix, histColName string, addHeadingNumbers bool, opts *SetTablePivotOpts) error {
	var setNamesOrdered []string
	if len(hsets.Order) > 0 {
		setNamesOrdered = hsets.Order
	} else {
		setNamesOrdered = maputil.Keys(hsets.Items)
	}
	for i, setName := range setNamesOrdered {
		headingNumber := ""
		if addHeadingNumbers {
			headingNumber = strconv.Itoa(i+1) + ". "
		}
		_, err := fmt.Fprintf(w, "%s%s%s\n\n", headingPrefix, headingNumber, setName)
		if err != nil {
			return err
		}
		hset, ok := hsets.Items[setName]
		if !ok {
			continue
		}
		if tbl, err := hset.TablePivot("", histColName, opts); err != nil {
			return err
		} else if _, err = fmt.Fprint(w, tbl.Markdown("\n", true)); err != nil {
			return err
		}
		if i < len(setNamesOrdered)-1 {
			if _, err = fmt.Fprint(w, "\n\n"); err != nil {
				return err
			}
		}
	}
	return nil
}
