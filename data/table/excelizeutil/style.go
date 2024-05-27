package excelizeutil

import (
	"github.com/grokify/mogo/image/colors"
	excelize "github.com/xuri/excelize/v2"
)

const (
	TypePattern                = "pattern"
	StyleBackgroundColorFormat = `{"fill":{"type":"pattern","color":["%s"],"pattern":1}}`
)

/*
func DefaultBackgroundColorFunc(colIdx, rowIdx uint) string {

}

func BackgroundColorStyle(hexColor string) (string, error) {
	hexColor = strings.TrimSpace(hexColor)
	if hexColor == "" {
		return nil, nil
	}
	hexRGB, err := colors.CanonicalHex(hexColor, true)
	if err != nil {
		return nil, err
	}
	return excelize.NewStyle(fmt.Sprintf(styleBackgroundColorFormat, hexRGB)), nil

}
*/

// StyleBackgroundColorSimple returns a `*excelize.Style` given a slice of RGB colors.
// A `nil` value is returned if there are no hex values supplied.
func StyleBackgroundColorSimple(hexRGBs []string) (*excelize.Style, error) {
	hexRGBs, err := colors.CanonicalHexes(hexRGBs, true, false, false, false)
	if err != nil {
		return nil, err
	}
	if len(hexRGBs) == 0 {
		return nil, nil
	}
	return &excelize.Style{
		Fill: excelize.Fill{
			Type:    TypePattern,
			Color:   hexRGBs,
			Pattern: 1,
		},
	}, nil
}

// style, err := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#E0EBF5"],"pattern":1}}`)
