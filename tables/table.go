package tables

const StyleSimple = "border:1px solid #000;border-collapse:collapse"

type TableData struct {
	Id    string
	Style string // border:1px solid #000;border-collapse:collapse
	Rows  [][]string
}
