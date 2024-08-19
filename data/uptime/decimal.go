package uptime

import (
	"errors"

	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/shopspring/decimal"
)

const PercentSuffix = "%"

type Decimals []decimal.Decimal

func (d Decimals) Strings(suffix string) []string {
	var out []string
	for _, di := range d {
		out = append(out, di.String()+suffix)
	}
	return out
}

func (d Decimals) Floats() []float64 {
	var out []float64
	for _, di := range d {
		if f, err := strconvutil.Atof(di.String(), false); err != nil {
			panic(err)
		} else {
			out = append(out, f)
		}
	}
	return out
}

func (d Decimals) Histogram(categories Decimals, suffix string) (*histogram.Histogram, error) {
	if len(categories) == 0 {
		categories = CategoriesDecimal()
	}
	h := histogram.NewHistogram("")
	h.Order = categories.Strings(suffix)

	dec100, err := decimal.NewFromString("100")
	if err != nil {
		panic(err)
	}
	dec0, err := decimal.NewFromString("0")
	if err != nil {
		panic(err)
	}
	for _, v := range d {
		if v.GreaterThan(dec100) {
			return nil, errors.New("value cannot be greater than 100")
		} else if v.LessThan(dec0) {
			return nil, errors.New("value cannot be less than 0")
		}
		for _, c := range categories {
			if v.GreaterThanOrEqual(c) {
				h.Add(c.String()+"%", 1)
				break
			}
		}
	}

	return h, nil
}

func CategoriesString() []string {
	return []string{"100", "99.999", "99.995", "99.99", "99.95", "99.9", "99.8", "99.5", "99", "98", "97", "95", "90", "85", "80", "75", "50", "0"}
}

func CategoriesDecimal() Decimals {
	var out []decimal.Decimal
	strs := CategoriesString()
	for _, s := range strs {
		if d, err := decimal.NewFromString(s); err != nil {
			panic(err)
		} else {
			out = append(out, d)
		}
	}
	return out
}
