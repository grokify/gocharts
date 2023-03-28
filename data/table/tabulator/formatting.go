// tabulator provides helper methods for rendering HTML
// with Tabulator (http://tabulator.info/)
package tabulator

import "encoding/json"

const (
	FormatterDatetime = "datetime"
	FormatterImage    = "image"
	FormatterMoney    = "money"

	ParamDecimal   = "decimal"
	ParamPrecision = "precision"
	ParamThousand  = "thousand"

	ParamInputFormat        = "inputFormat"
	ParamOutputFormat       = "outputFormat"
	ParamInvalidPlaceholder = "invalidPlaceholder"
	ParamTimezone           = "timezone"
)

// FormatterParams returns params for use in https://tabulator.info/docs/5.4/format
type FormatterParams struct {
	Formatter string
	// Money
	Decimal   string
	Thousand  string
	Symbol    string
	Precision int
	// DateTime
	InputFormat        string // (default: yyyy-MM-dd HH:mm:ss), can set to `iso`
	OutputFormat       string // (default: dd/MM/yyyy HH:mm:ss)
	InvalidPlaceholder string // can set to `utc`
	Timezone           string
}

/*

{title:"Example", field:"example", formatter:"datetime", formatterParams:{
    inputFormat:"yyyy-MM-dd HH:ss",
    outputFormat:"dd/MM/yy",
    invalidPlaceholder:"(invalid date)",
    timezone:"America/Los_Angeles",
}}

*/

func (fp *FormatterParams) MarshalJSON() ([]byte, error) {
	// type Alias FormatterParams // http://choly.ca/post/go-json-marshalling/
	msa := fp.JS()
	return json.Marshal(msa)
}

func (fp FormatterParams) JS() map[string]any {
	msa := map[string]any{}
	// Money
	switch fp.Formatter {
	case FormatterMoney:
		if fp.Decimal != "" {
			msa["decimal"] = fp.Decimal
		}
		if fp.Thousand == "" {
			msa[ParamThousand] = false
		} else {
			msa[ParamThousand] = fp.Thousand
		}
		if fp.Symbol != "" {
			msa["symbol"] = fp.Symbol
		}
		if fp.Precision < 0 {
			msa[ParamPrecision] = false
		} else {
			msa[ParamPrecision] = fp.Precision
		}
	case FormatterDatetime:
		// Datetime
		if fp.InputFormat != "" {
			msa[ParamInputFormat] = fp.InputFormat
		}
		if fp.OutputFormat != "" {
			msa[ParamOutputFormat] = fp.OutputFormat
		}
		if fp.InvalidPlaceholder != "" {
			msa[ParamInvalidPlaceholder] = fp.InvalidPlaceholder
		}
		if fp.Timezone != "" {
			msa[ParamTimezone] = fp.Timezone
		}
	}
	return msa
}
