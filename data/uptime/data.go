package uptime

import (
	"errors"
	"os"
	"sort"

	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/shopspring/decimal"
)

type Datum struct {
	Title         string
	Description   string
	UptimeDecimal *decimal.Decimal
	// UptimeFloat   *float64
}

type Data []Datum

func (d Data) Filter(fnIncl func(md Datum) bool) Data {
	out := Data{}
	for _, m := range d {
		if fnIncl(m) {
			out = append(out, m)
		}
	}
	return out
}

func (d Data) Decimals() Decimals {
	var decs Decimals
	for _, m := range d {
		if m.UptimeDecimal != nil {
			decs = append(decs, pointer.Dereference(m.UptimeDecimal))
		}
	}
	return decs
}

func (d Data) HistogramDecimal(categories []decimal.Decimal, suffix string) (*histogram.Histogram, error) {
	vals := d.Decimals()
	return vals.Histogram(categories, suffix)
}

func UptimeCounts(categories, observations []float64) (map[float64]uint, error) {
	out := map[float64]uint{}
	for _, obv := range observations {
		if obv > 100 {
			return out, errors.New("value cannot be greater than 100")
		} else if obv < 0 {
			return out, errors.New("valuyeu cannot be less than 0")
		}
		for _, cat := range categories {
			if obv >= cat {
				out[cat]++
				break
			}
		}
	}
	return out, nil
}

func Float64ToUptimeString(digits uint, f float32) string {
	return strconvutil.Ftoa(f, 32) + "%"
	// return fmt.Sprintf("%."+strconv.Itoa(int(digits))+"f%%", f)
	// return fmt.Sprintf("."+strconv.Itoa(int(digits))+"%f%%", f)
}

func (d Data) Descriptions(fnIncl func(s string) bool, fnFmt func(s string) string, sortAsc bool) []string {
	var out []string
	for _, m := range d {
		if fnIncl != nil && !fnIncl(m.Description) {
			continue
		} else if fnFmt != nil {
			out = append(out, fnFmt(m.Description))
		} else {
			out = append(out, m.Description)
		}
	}
	if sortAsc {
		sort.Strings(out)
	}
	return out
}

func (d Data) Table(name string) *table.Table {
	tbl := table.NewTable(name)
	tbl.Columns = []string{"Title", "Description", "Uptime Percent"}
	tbl.FormatMap = map[int]string{2: table.FormatFloat}
	for _, di := range d {
		uptime := "0"
		if di.UptimeDecimal != nil {
			uptime = di.UptimeDecimal.String()
		}
		tbl.Rows = append(tbl.Rows, []string{di.Title, di.Description, uptime})
	}
	return &tbl
}

func (d Data) WriteXLSX(filename, sheetname string, perm os.FileMode) error {
	tbl := d.Table(sheetname)
	return tbl.WriteXLSX(filename, sheetname)
}
