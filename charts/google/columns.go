package google

const (
	TypeNumber = "number"
	TypeString = "string"
)

type Columns []Column

func (cols Columns) Names() []string {
	var names []string
	for _, c := range cols {
		names = append(names, c.Name)
	}
	return names
}

func (cols Columns) NamesAny() []any {
	var names []any
	for _, c := range cols {
		names = append(names, c.Name)
	}
	return names
}

type Column struct {
	Type string
	Name string
}
