package table

type FormatMap map[int]string

func (fm FormatMap) FormatForIdx(idx int) string {
	defFormat := FormatString
	if idx <= -1 {
		idx = -1
	}
	if f, ok := fm[idx]; ok {
		return f
	} else {
		return defFormat
	}
}
