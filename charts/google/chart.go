package google

import "os"

type Chart interface {
	PageTitle() string
	ChartDivOrDefault() string
	DataTableJSON() []byte
	OptionsJSON() []byte
	WriteFilePage(filename string, perm os.FileMode) error
}
