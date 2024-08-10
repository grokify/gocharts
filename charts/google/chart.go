package google

type Chart interface {
	PageTitle() string
	ChartDivOrDefault() string
	DataTableJSON() []byte
	OptionsJSON() []byte
}
